package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/malbrecht/chess/pgn"
	"github.com/pkg/profile"
)

//Stats is what we'll write in as JSON
type Stats struct {
	TotalGames uint32       `json:"totalGames"`
	Openings   *OpeningMove `json:"openings"`
	Heatmaps   struct {
		SquareUtilization Heatmap `json:"squareUtilization"`
		MoveSquares       Heatmap `json:"moveSquares"`
		CaptureSquares    Heatmap `json:"captureSquares"`
		CheckSquares      Heatmap `json:"checkSquares"`
	} `json:"heatmaps"`
	Results struct {
		White uint32
		Black uint32
		Draw  uint32
		NA    uint32
	}
	GamesEndingWith struct {
		Check uint32
		Mate  uint32
	}
	GameLengths   ConcurrentMap
	MaterialCount ConcurrentMap
	MaterialDiff  ConcurrentMap
	Castling      struct {
		White struct {
			Kingside  uint32
			Queenside uint32
		}
		Black struct {
			Kingside  uint32
			Queenside uint32
		}
		Side struct {
			Same     uint32
			Opposite uint32
		}
	}
}

//Game sure
type Game struct {
	PgnGame *pgn.Game
	Moves   []string
}

//Read reads the PGN file in chunks and constructs []byte with contents of a single game
func Read(f *os.File) <-chan []byte {
	c := make(chan []byte)

	scanner := bufio.NewScanner(f)
	var bytes []byte

	go func() {
		defer close(c)
		for scanner.Scan() {
			line := scanner.Bytes()
			lineStr := string(line)
			bytes = append(bytes, line...)
			bytes = append(bytes, '\n')

			match := strings.HasSuffix(lineStr, "1-0") ||
				strings.HasSuffix(lineStr, "0-1") ||
				strings.HasSuffix(lineStr, "1/2-1/2") ||
				strings.HasSuffix(lineStr, "*")

			if !match {
				continue
			}

			c <- []byte(bytes)
			bytes = make([]byte, 0)
		}
	}()

	return c
}

//Parse parses games coming from r and send them off to s
func Parse(r <-chan []byte, s chan<- *Game) {
	tagsRegex := regexp.MustCompile("\\[.+\\]")
	movesRegex := regexp.MustCompile("\\d+\\.")

	for game := range r {
		db := pgn.DB{}
		err := db.Parse(string(game))
		if err != nil {
			log.Printf("error parsing game: %s\n", err)
			continue
		}

		if _, ok := db.Games[0].Tags["SetUp"]; ok {
			if *verbose {
				log.Println("SetUp tag found, won't bother parsing it")
			}
			continue
		}

		db.ParseMoves(db.Games[0])

		if _, ok := db.Games[0].Tags["Result"]; !ok {
			if *verbose {
				log.Println("no result for game! skipping", db.Games[0].Tags)
			}
			continue
		}

		movesBytes := tagsRegex.ReplaceAll(game, nil)
		movesBytes = movesRegex.ReplaceAll(movesBytes, nil)
		movesString := strings.Trim(string(movesBytes), "\n ")
		moves := strings.Fields(movesString)

		s <- &Game{
			db.Games[0],
			moves,
		}
	}
}

//GetStats collects statistics from games
func GetStats(c <-chan *Game, data *Stats) {
	for Game := range c {
		gamePtr := Game.PgnGame.Root
		openingPtr := data.Openings
		castle := ""

		atomic.AddUint32(&data.TotalGames, 1)
		atomic.AddUint32(&data.Openings.Count, 1)

		for i := 0; i < len(Game.Moves)-1; i++ {
			gamePtr = gamePtr.Next
			move := gamePtr.Move
			rawMove := Game.Moves[i]
			piece := gamePtr.Board.Piece[move.To]

			HeatmapStats(data, move, piece, rawMove)
			MaterialCountStats(data, gamePtr.Board, i)

			if rawMove == "O-O" || rawMove == "O-O-O" {
				CastlingStats(data, rawMove, gamePtr.Board.SideToMove)

				if castle == "" {
					castle = rawMove
				} else {
					if rawMove == castle {
						atomic.AddUint32(&data.Castling.Side.Same, 1)
					} else {
						atomic.AddUint32(&data.Castling.Side.Opposite, 1)
					}
				}
			}

			if i < 10 {
				openingPtr = OpeningStats(openingPtr, rawMove)
			}
		}

		//results
		switch Game.Moves[len(Game.Moves)-1] {
		case "1-0":
			atomic.AddUint32(&data.Results.White, 1)
		case "0-1":
			atomic.AddUint32(&data.Results.Black, 1)
		case "1/2-1/2":
			atomic.AddUint32(&data.Results.Draw, 1)
		default:
			atomic.AddUint32(&data.Results.NA, 1)
		}

		//game lengths histogram
		gamePlyStr := strconv.Itoa(len(Game.Moves) - 1)

		if val, ok := data.GameLengths.Get(gamePlyStr); ok {
			data.GameLengths.Set(gamePlyStr, val+1)
		} else {
			data.GameLengths.Set(gamePlyStr, 1)
		}

		if len(Game.Moves) < 2 {
			continue
		}

		//games ending with check/mate
		lastMove := Game.Moves[len(Game.Moves)-2]

		if strings.ContainsRune(lastMove, '+') {
			atomic.AddUint32(&data.GamesEndingWith.Check, 1)
		}

		if strings.ContainsRune(lastMove, '#') {
			atomic.AddUint32(&data.GamesEndingWith.Mate, 1)
		}
	}
}

//flags
var pgnPath = flag.String("f", "./pgn/a.pgn", "path of the PGN file")
var concurrencyLevel = flag.Int("c", 10, "concurrency level for parsing")
var outputPath = flag.String("o", "./data.json", "output path for JSON file")
var perf = flag.Bool("p", false, "write profile to ./prof/")
var verbose = flag.Bool("v", false, "verbose mode")
var indent = flag.Bool("i", false, "indent json output")

func main() {
	flag.Parse()

	if *perf {
		defer profile.Start(profile.ProfilePath("./prof")).Stop()
	}

	//open file
	f, err := os.Open(*pgnPath)
	if err != nil {
		log.Fatalf("failed to open file: %s\n", err)
	}

	log.Printf("starting")

	//init data
	stats := &Stats{
		Openings: &OpeningMove{
			San:      "start",
			Children: make([]*OpeningMove, 0),
		},
		GameLengths:   NewCmap(),
		MaterialCount: NewCmap(),
		MaterialDiff:  NewCmap(),
	}

	readC := Read(f)
	parsedC := make(chan *Game)

	var wg sync.WaitGroup
	wg.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			Parse(readC, parsedC)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(parsedC)
		log.Println("parsing done")
	}()

	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			GetStats(parsedC, stats)
			wg2.Done()
		}()
	}

	wg2.Wait()

	log.Printf("analyzed %d games\n", stats.TotalGames)

	pruneThreshold := int(float32(stats.TotalGames) * 0.001)

	if *verbose {
		log.Printf("prune param %d\n", pruneThreshold)
	}

	stats.Openings.Prune(pruneThreshold)

	var js []byte
	if *indent {
		js, err = json.MarshalIndent(stats, "", "  ")
	} else {
		js, err = json.Marshal(stats)
	}
	if err != nil {
		log.Fatalf("error converting to json: %s\n", err)
	}

	err = ioutil.WriteFile(*outputPath, js, 0644)
	if err != nil {
		log.Fatalf("error writing file: %s\n", err)
	}

	log.Printf("done!")
}
