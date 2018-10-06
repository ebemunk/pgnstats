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

	//init data
	stats := NewResult()
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
	}()

	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			GameStats(parsedC, stats)
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
