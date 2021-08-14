package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/malbrecht/chess/pgn"
)

func LoadGames(file string) []*pgn.Game {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("cannot read test file %s\n", file)
	}

	db := pgn.DB{}
	db.Parse(string(dat))

	for i := range db.Games {
		err = db.ParseMoves(db.Games[i])
		if err != nil {
			log.Fatalf("cannot parse game in %s\n", err)
		}
	}

	return db.Games
}

func findGame(games []*pgn.Game, name string) *pgn.Game {
	for _, game := range games {
		if game.Tags["Event"] == name {
			return game
		}
	}

	return nil
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
	games := LoadGames("./testdata/pgn/fools_mate.pgn")
	op := OpeningMove{}

	t.Run("fools", func(t *testing.T) {
		stats := GetStats(findGame(games, "Fool's Mate"), &op, "")
		writeTestJSON(stats, "./testdata/results/fools_mate.json")
	})

	t.Run("scholars", func(t *testing.T) {
		stats := GetStats(findGame(games, "Scholar's Mate"), &op, "")
		writeTestJSON(stats, "./testdata/results/scholars_mate.json")
	})

	t.Run("repetition", func(t *testing.T) {
		stats := GetStats(findGame(games, "Repetition"), &op, "")
		writeTestJSON(stats, "./testdata/results/repetition.json")
	})
}
