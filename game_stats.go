package main

import (
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/dylhunn/dragontoothmg"
)

// GameStats collects statistics from a game
func GameStats(c <-chan *Game, data *Result) {
	for Game := range c {
		gamePtr := Game.PgnGame.Root
		openingPtr := data.Openings
		castle := ""
		var firstCapture = false

		atomic.AddUint32(&data.TotalGames, 1)
		atomic.AddUint32(&data.Openings.Count, 1)

		for ply := 0; ply < len(Game.Moves)-1; ply++ {
			gamePtr = gamePtr.Next
			move := gamePtr.Move
			rawMove := Game.Moves[ply]
			piece := gamePtr.Board.Piece[move.To]

			if !firstCapture {
				firstCapture = FirstBlood(&data.Heatmaps.FirstBlood, gamePtr)
			}

			HeatmapStats(data, move, piece, rawMove)

			if rawMove == "O-O" || rawMove == "O-O-O" {
				CastlingStats(data, rawMove, gamePtr.Board.SideToMove)

				if castle == "" {
					castle = rawMove
				} else {
					if rawMove == castle {
						atomic.AddUint32(&data.Castling.Side.Same, 1)
					} else {
						atomic.AddUint32(&data.Castling.Side.Opposite, 1)
					}
				}
			}

			if ply < 10 {
				openingPtr = OpeningStats(openingPtr, rawMove)
			}

			//BranchingFactor
			board := dragontoothmg.ParseFen(gamePtr.Board.Fen())
			branchingFactor := float64(len(board.GenerateLegalMoves()))

			val, loaded := data.BranchingFactor.LoadOrStore(ply, branchingFactor)
			if loaded {
				data.BranchingFactor.Store(ply, ((val.(float64)*float64(data.TotalGames))+branchingFactor)/(float64(data.TotalGames)+1))
			}

			//MaterialCount
			count, diff := MaterialCount(gamePtr.Board)
			val, loaded = data.MaterialCount.LoadOrStore(ply, float64(count))
			if loaded {
				data.MaterialCount.Store(ply, ((val.(float64)*float64(data.TotalGames))+float64(count))/(float64(data.TotalGames)+1))
			}

			//MaterialDiff
			val, loaded = data.MaterialDiff.LoadOrStore(ply, float64(diff))
			if loaded {
				data.MaterialDiff.Store(ply, ((val.(float64)*float64(data.TotalGames))+float64(diff))/(float64(data.TotalGames)+1))
			}
		}

		//results
		switch Game.Moves[len(Game.Moves)-1] {
		case "1-0":
			atomic.AddUint32(&data.Results.White, 1)
		case "0-1":
			atomic.AddUint32(&data.Results.Black, 1)
		case "1/2-1/2":
			atomic.AddUint32(&data.Results.Draw, 1)
		default:
			atomic.AddUint32(&data.Results.NA, 1)
		}

		//ratings
		if elo, ok := Game.PgnGame.Tags["WhiteElo"]; ok {
			eloNum64, _ := strconv.Atoi(elo)
			eloNum := uint32(eloNum64)

			if eloNum < atomic.LoadUint32(&data.Ratings.Min) {
				atomic.StoreUint32(&data.Ratings.Min, eloNum)
			}

			if eloNum > atomic.LoadUint32(&data.Ratings.Max) {
				atomic.StoreUint32(&data.Ratings.Max, eloNum)
			}
		}

		if elo, ok := Game.PgnGame.Tags["BlackElo"]; ok {
			eloNum64, _ := strconv.Atoi(elo)
			eloNum := uint32(eloNum64)

			if eloNum < atomic.LoadUint32(&data.Ratings.Min) {
				atomic.StoreUint32(&data.Ratings.Min, eloNum)
			}

			if eloNum > atomic.LoadUint32(&data.Ratings.Max) {
				atomic.StoreUint32(&data.Ratings.Max, eloNum)
			}
		}

		//dates
		if date, ok := Game.PgnGame.Tags["Date"]; ok {
			year64, _ := strconv.Atoi(date[:4])
			year := uint32(year64)

			if year < atomic.LoadUint32(&data.Dates.Min) {
				atomic.StoreUint32(&data.Dates.Min, year)
			}

			if year > atomic.LoadUint32(&data.Dates.Max) {
				atomic.StoreUint32(&data.Dates.Max, year)
			}
		}

		//game lengths
		totalPlies := len(Game.Moves) - 1
		val, loaded := data.GameLengths.LoadOrStore(totalPlies, float64(1))
		if loaded {
			data.GameLengths.Store(totalPlies, val.(float64)+1)
		}

		if len(Game.Moves) < 2 {
			continue
		}

		//games ending with check/mate
		lastMove := Game.Moves[len(Game.Moves)-2]

		if strings.ContainsRune(lastMove, '+') {
			atomic.AddUint32(&data.GamesEndingWith.Check, 1)
		}

		if strings.ContainsRune(lastMove, '#') {
			atomic.AddUint32(&data.GamesEndingWith.Mate, 1)
		}
	}
}
