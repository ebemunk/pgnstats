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

	readC := Read(f)
	parsedC := make(chan *Game)
	gsC := make(chan *GameStats)

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
		log.Println("close parsedC")
	}()

	var wg2 sync.WaitGroup
	wg2.Add(*concurrencyLevel)
	for i := 0; i < *concurrencyLevel; i++ {
		go func() {
			GetStats(parsedC, gsC)
			wg2.Done()
		}()
	}

	go func() {
		wg2.Wait()
		close(gsC)
		log.Println("close gsC")
	}()

	var viji sync.WaitGroup
	viji.Add(1)
	go func() {
		fgs := NewGameStats()

		for gamst := range gsC {
			fgs.Add(gamst)
		}

		fgs.Average()

		WriteJson(fgs)

		viji.Done()

	}()

	log.Println("before viji wait")
	viji.Wait()
	log.Println("after viji wait")

	// log.Printf("analyzed %d games\n", stats.TotalGames)

	// pruneThreshold := int(float32(stats.TotalGames) * 0.001)

	// if *verbose {
	// 	log.Printf("prune param %d\n", pruneThreshold)
	// }

	// stats.Openings.Prune(pruneThreshold)
}

func WriteJson(gs *GameStats) {
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
