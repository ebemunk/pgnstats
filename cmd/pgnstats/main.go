package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/malbrecht/chess/pgn"
	"github.com/pkg/profile"

	"github.com/ebemunk/pgnstats/core"
)

//flags
var pgnPath = flag.String("f", "", "path of the PGN file")
var concurrencyLevel = flag.Int("c", 10, "concurrency level for parsing")
var outputPath = flag.String("o", "", "output path for JSON file")
var perf = flag.Bool("p", false, "write profile to ./prof/")
var verbose = flag.Bool("v", false, "verbose mode")
var indent = flag.Bool("i", false, "indent json output")
var filterPlayer = flag.String("fp", "", "filter by player")

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

	readC := Read(f)
	parsedC := make(chan *pgn.Game)
	gsC := make(chan *core.GameStats)

	//read the file & parse
	var wg sync.WaitGroup
	wg.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			Parse(readC, parsedC)
			wg.Done()
		}()
	}

	//close parsedC when parsing is complete
	go func() {
		wg.Wait()
		close(parsedC)
		if *verbose {
			log.Println("close parsedC")
		}
	}()

	playerStats := core.NewPlayerStats()

	var Openings []*core.OpeningMove = make([]*core.OpeningMove, 0)
	Openings = append(Openings, &core.OpeningMove{San: "start"})
	playerStats.All.Openings = Openings[0]
	if *filterPlayer != "" {
		Openings = append(Openings, &core.OpeningMove{San: "start"})
		playerStats.White.Openings = Openings[1]
		Openings = append(Openings, &core.OpeningMove{San: "start"})
		playerStats.Black.Openings = Openings[2]
	}

	//collect stats
	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			for Game := range parsedC {
				stats := core.NewGameStatsFromGame(Game, *filterPlayer, Openings)

				if stats != nil {
					gsC <- stats
				}
			}
			wg2.Done()
		}()
	}

	//close gsC when complete
	go func() {
		wg2.Wait()
		close(gsC)
		if *verbose {
			log.Println("close gsC")
		}
	}()

	//combine stats
	var wg3 sync.WaitGroup
	wg3.Add(1)
	go func() {
		for gamst := range gsC {
			playerStats.All.Add(gamst)
			if gamst.Color == "w" {
				playerStats.White.Add(gamst)
			} else if gamst.Color == "b" {
				playerStats.Black.Add(gamst)
			}
		}

		if *filterPlayer != "" {
			log.Printf("analyzed %d games (%d white, %d black)\n", playerStats.All.Total, playerStats.White.Total, playerStats.Black.Total)
		} else {
			log.Printf("analyzed %d games\n", playerStats.All.Total)
		}

		pruneThreshold := int(float32(playerStats.All.Total) * 0.005)
		if *verbose {
			log.Printf("prune param %d\n", pruneThreshold)
		}

		playerStats.All.Average()
		playerStats.White.Average()
		playerStats.Black.Average()

		playerStats.All.Openings.Prune(pruneThreshold)
		if *filterPlayer != "" {
			playerStats.White.Openings.Prune(pruneThreshold)
			playerStats.Black.Openings.Prune(pruneThreshold)
		}

		playerStats.All.Positions.Prune(pruneThreshold)
		playerStats.White.Positions.Prune(pruneThreshold)
		playerStats.Black.Positions.Prune(pruneThreshold)

		playerStats.All.UniquePositions.Prune(pruneThreshold)
		playerStats.White.UniquePositions.Prune(pruneThreshold)
		playerStats.Black.UniquePositions.Prune(pruneThreshold)

		if *filterPlayer == "" {
			WriteJSON(playerStats.All, "allgames")
		} else {
			WriteJSON(playerStats, "filtered")
		}

		wg3.Done()
	}()

	wg3.Wait()

	log.Printf("done!")
}
