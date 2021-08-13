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
var filterPlayer = flag.String("fp", "Carlsen,M", "filter by player")

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
	Openings.San = "start"

	//collect stats
	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			GetStats(parsedC, gsC, Openings, *filterPlayer)
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
		// totals or White when -fp
		wgs := NewGameStats()
		wgs.Color = "W"
		// Black when -fp
		bgs := NewGameStats()
		bgs.Color = "B"

		for gamst := range gsC {
			if gamst.Color == "b" {
				bgs.Add(gamst)
			} else {
				wgs.Add(gamst)
			}
		}

		wgs.Average()
		bgs.Average()

		if *filterPlayer != "" {
			log.Printf("analyzed %d games (%d white, %d black)\n", wgs.Total+bgs.Total, wgs.Total, bgs.Total)
		} else {
			log.Printf("analyzed %d games\n", wgs.Total)
		}

		pruneThreshold := int(float32(wgs.Total) * 0.0001)
		if *verbose {
			log.Printf("prune param %d\n", pruneThreshold)
		}

		Openings.Prune(pruneThreshold)
		wgs.Openings = Openings

		prunedPos := make(PosMap)

		for k, v := range wgs.Positions {
			if v > pruneThreshold {
				prunedPos[k] = v
			}
		}
		wgs.Positions = prunedPos

		for k, v := range bgs.Positions {
			if v > pruneThreshold {
				prunedPos[k] = v
			}
		}
		bgs.Positions = prunedPos

		if *filterPlayer == "" {
			writeJSON(wgs, "all")
		} else {
			writeJSON(wgs, "w")
		}

		if *filterPlayer != "" {
			writeJSON(bgs, "b")
		}

		wg3.Done()
	}()

	wg3.Wait()

	log.Printf("done!")
}

func writeJSON(gs *GameStats, suffix string) {
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

	filePath := *outputPath + "-" + suffix + ".json"
	err = ioutil.WriteFile(filePath, js, 0644)
	if err != nil {
		log.Fatalf("error writing file: %s\n", err)
	}
	log.Printf("wrote to %v", filePath)
}
