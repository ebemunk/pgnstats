package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/malbrecht/chess/pgn"
)

//Read reads a PGN file in chunks and constructs []byte with contents of a single game
func Read(f *os.File) <-chan []byte {
	c := make(chan []byte)

	scanner := bufio.NewScanner(f)
	var bytes []byte

	go func() {
		defer close(c)
		if *verbose {
			defer log.Println("close read chan")
		}

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
func Parse(r <-chan []byte, s chan<- *pgn.Game) {
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

		if _, ok := db.Games[0].Tags["Result"]; !ok {
			if *verbose {
				log.Println("no result for game! skipping", db.Games[0].Tags)
			}
			continue
		}

		db.ParseMoves(db.Games[0])

		s <- db.Games[0]
	}
}

func WriteJSON(ps interface{}, suffix string) {
	var js []byte
	var err error

	if *indent {
		js, err = json.MarshalIndent(ps, "", "  ")
	} else {
		js, err = json.Marshal(ps)
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
