package core

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/malbrecht/chess/pgn"
	"github.com/sebdah/goldie/v2"
)

func loadGames(file string) []*pgn.Game {
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

func TestNewGameStatsFromGame(t *testing.T) {
	games := loadGames("./testdata/pgn/games.pgn")
	g := goldie.New(t, goldie.WithNameSuffix(".golden.json"))

	var Openings []*OpeningMove = make([]*OpeningMove, 1)
	Openings[0] = &OpeningMove{San: "start"}

	t.Run("fools", func(t *testing.T) {
		stats := NewGameStatsFromGame(findGame(games, "Fool's Mate"), "", Openings)
		g.AssertJson(t, "fools_mate", stats)
	})

	t.Run("scholars", func(t *testing.T) {
		stats := NewGameStatsFromGame(findGame(games, "Scholar's Mate"), "", Openings)
		g.AssertJson(t, "scholars_mate", stats)
	})

	t.Run("repetition", func(t *testing.T) {
		stats := NewGameStatsFromGame(findGame(games, "Repetition"), "", Openings)
		g.AssertJson(t, "repetition", stats)
	})

	t.Run("repetition-opening", func(t *testing.T) {
		var Openings []*OpeningMove = make([]*OpeningMove, 3)
		Openings[0] = &OpeningMove{San: "start"}
		Openings[1] = &OpeningMove{San: "start"}
		Openings[2] = &OpeningMove{San: "start"}
		NewGameStatsFromGame(findGame(games, "Repetition"), "White", Openings)
		g.AssertJson(t, "repetition-opening", Openings)
	})
}
