package main

import (
	"testing"
	// "github.com/sebdah/goldie/v2"
)

func TestNewGameStatsFromGame(t *testing.T) {
	// games := loadGames("./testdata/pgn/games.pgn")
	// g := goldie.New(t, goldie.WithNameSuffix(".golden.json"))

	t.Run("carlsen", func(t *testing.T) {

	})
	// t.Run("fools", func(t *testing.T) {
	// 	stats := NewGameStatsFromGame(findGame(games, "Fool's Mate"), "")
	// 	g.AssertJson(t, "fools_mate", stats)
	// })

	// t.Run("scholars", func(t *testing.T) {
	// 	stats := NewGameStatsFromGame(findGame(games, "Scholar's Mate"), "")
	// 	g.AssertJson(t, "scholars_mate", stats)
	// })

	// t.Run("repetition", func(t *testing.T) {
	// 	stats := NewGameStatsFromGame(findGame(games, "Repetition"), "")
	// 	g.AssertJson(t, "repetition", stats)
	// })
}
