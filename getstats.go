package main

import (
	"strconv"
	"sync/atomic"

	"github.com/dylhunn/dragontoothmg"
	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/pgn"
)

//GetStats collects statistics from a game
func GetStats(c <-chan *pgn.Game, gs chan<- *GameStats, openingsPtr *OpeningMove) {
	for Game := range c {
		oPtr := openingsPtr

		stats := NewGameStats()

		ply := -1
		var firstCapture = false

		for gamePtr := Game.Root; gamePtr != nil; gamePtr = gamePtr.Next {
			ply++

			move := gamePtr.Move
			isLastMove := gamePtr.Next == nil

			//Openings
			if ply > 0 && ply < 10 {
				atomic.AddUint32(&oPtr.Count, 1)
				oPtr = OpeningStats(oPtr, gamePtr.Move.San(gamePtr.Parent.Board))
			}

			//BranchingFactor
			board := dragontoothmg.ParseFen(gamePtr.Board.Fen())
			branchingFactor := float64(len(board.GenerateLegalMoves()))
			stats.BranchingFactor[ply] += branchingFactor

			//Heatmaps
			HeatmapStats(stats, gamePtr, isLastMove)

			if ply > 0 && !firstCapture {
				firstCapture = FirstBlood(&stats.Heatmaps.FirstBlood, gamePtr)
			}

			//MaterialCount
			count, diff := MaterialCount(gamePtr.Board)
			stats.MaterialCount[ply] = float64(count)
			stats.MaterialDiff[ply] = float64(diff)
			if isLastMove {
				stats.GameEndMaterialCount[ply] = float64(count)
				stats.GameEndMaterialDiff[ply] = float64(diff)
			}

			//PromotionSquares
			if move.Promotion != chess.NoPiece {
				stats.Heatmaps.PromotionSquares.Count(move.Promotion, move.To)
			}

			//EnPassantSquares
			if gamePtr.Board.EpSquare != chess.NoSquare {
				stats.Heatmaps.EnPassantSquares.Count(gamePtr.Parent.Board.Piece[move.From], gamePtr.Board.EpSquare)
			}
		}

		//Ratings
		if elo, ok := Game.Tags["WhiteElo"]; ok {
			stats.Ratings[elo] = 1
		}

		if elo, ok := Game.Tags["BlackElo"]; ok {
			stats.Ratings[elo] = 1
		}

		//Years
		if date, ok := Game.Tags["UTCDate"]; ok {
			year64, _ := strconv.Atoi(date[:4])
			year := strconv.Itoa(year64)

			stats.Years[year] = 1
		}

		//GameLengths
		stats.GameLengths[ply] = 1

		gs <- stats
	}
}
