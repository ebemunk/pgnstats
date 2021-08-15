package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/malbrecht/chess/pgn"
	"github.com/sebdah/goldie/v2"
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

func TestGetStats(t *testing.T) {
	games := LoadGames("./testdata/pgn/games.pgn")
	op := OpeningMove{}
	g := goldie.New(t, goldie.WithNameSuffix(".golden.json"))

	t.Run("fools", func(t *testing.T) {
		stats := GetStats(findGame(games, "Fool's Mate"), &op, "")
		g.AssertJson(t, "fools_mate", stats)
	})

	t.Run("scholars", func(t *testing.T) {
		stats := GetStats(findGame(games, "Scholar's Mate"), &op, "")
		g.AssertJson(t, "scholars_mate", stats)
	})

	t.Run("repetition", func(t *testing.T) {
		stats := GetStats(findGame(games, "Repetition"), &op, "")
		g.AssertJson(t, "repetition", stats)
	})
}
