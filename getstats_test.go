package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/malbrecht/chess/pgn"
)

func GetGame(file string) *pgn.Game {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("cannot read test file %s\n", file)
	}

	db := pgn.DB{}
	db.Parse(string(dat))
	err = db.ParseMoves(db.Games[0])
	if err != nil {
		log.Fatalf("cannot parse game in %s\n", file)
	}

	return db.Games[0]
}

func writeTestJSON(stats *GameStats, path string) {
	json, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatalf("error marshalling json: %s\n", err)
	}

	err = ioutil.WriteFile(path, json, 0644)
	if err != nil {
		log.Fatalf("error writing file: %s\n", err)
	}
}

func TestGetStats(t *testing.T) {
	game := GetGame("./testdata/fools_mate.pgn")
	op := OpeningMove{}
	stats := GetStats(game, &op, "")
	log.Printf("%+v\n", stats)
	writeTestJSON(stats, "./testdata/fools_mate.json")
}
