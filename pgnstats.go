package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/pkg/profile"
)

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
		Ratings: MinMax{
			Min: 3000,
			Max: 0,
		},
		Dates: MinMax{
			Min: 3000,
			Max: 0,
		},
		BranchingFactor: FloatMap{},
	}

	readC := Read(f)
	parsedC := make(chan *Game)

	var wg sync.WaitGroup
	wg.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		// log.Printf("parse worker %d\n", i)
		go func() {
			Parse(readC, parsedC)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(parsedC)
		// log.Println("closed ")
	}()

	// var i = 0
	// for game := range parsedC {
	// 	i += game.Kek
	// }
	// log.Printf("analyzed %d games\n", i)

	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		// log.Printf("stats worker %d\n", i)
		go func() {
			GetStats(parsedC, stats)
			wg2.Done()
		}()
	}

	wg2.Wait()

	log.Printf("analyzed %d games\n", stats.TotalGames)

	// pruneThreshold := int(float32(stats.TotalGames) * 0.001)

	// if *verbose {
	// 	log.Printf("prune param %d\n", pruneThreshold)
	// }

	// stats.Openings.Prune(pruneThreshold)

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
