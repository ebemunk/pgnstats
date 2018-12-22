package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/malbrecht/chess/pgn"
	"github.com/pkg/profile"
)

//flags
var pgnPath = flag.String("f", "./pgn/a.pgn", "path of the PGN file")
var concurrencyLevel = flag.Int("c", 10, "concurrency level for parsing")
var outputPath = flag.String("o", "./data/data.json", "output path for JSON file")
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

	readC := Read(f)
	parsedC := make(chan *pgn.Game)
	gsC := make(chan *GameStats)

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

	Openings := &OpeningMove{}

	//collect stats
	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			GetStats(parsedC, gsC, Openings)
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
		fgs := NewGameStats()

		for gamst := range gsC {
			fgs.Add(gamst)
		}

		fgs.Average()

		log.Printf("analyzed %d games\n", fgs.Total)

		pruneThreshold := int(float32(fgs.Total) * 0.0001)
		if *verbose {
			log.Printf("prune param %d\n", pruneThreshold)
		}

		Openings.Prune(pruneThreshold)
		fgs.Openings = Openings

		prunedPos := make(PosMap)

		for k, v := range fgs.Positions {
			if v > pruneThreshold {
				prunedPos[k] = v
			}
		}
		fgs.Positions = prunedPos

		writeJSON(fgs)

		wg3.Done()
	}()

	wg3.Wait()
}

func writeJSON(gs *GameStats) {
	var js []byte
	var err error

	if *indent {
		js, err = json.MarshalIndent(gs, "", "  ")
	} else {
		js, err = json.Marshal(gs)
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
