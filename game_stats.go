package main

import (
	"github.com/dylhunn/dragontoothmg"
)

// GetStats collects statistics from a game
func GetStats(c <-chan *Game, gs chan<- *GameStats) {
	for Game := range c {
		// openingPtr := data.Openings
		// castle := ""

		// atomic.AddUint32(&data.TotalGames, 1)
		// atomic.AddUint32(&data.Openings.Count, 1)

		stats := NewGameStats()

		ply := -1
		// var lastPosition *pgn.Node
		var firstCapture = false

		for gamePtr := Game.PgnGame.Root; gamePtr != nil; gamePtr = gamePtr.Next {
			ply++
			// lastPosition = gamePtr

			// move := gamePtr.Move
			// rawMove := Game.Moves[ply]
			// piece := gamePtr.Board.Piece[move.To]

			if ply > 0 && !firstCapture {
				firstCapture = FirstBlood(&stats.Heatmaps.FirstBlood, gamePtr)
			}

			// HeatmapStats(data, gamePtr)

			// if rawMove == "O-O" || rawMove == "O-O-O" {
			// 	CastlingStats(data, rawMove, gamePtr.Board.SideToMove)

			// 	if castle == "" {
			// 		castle = rawMove
			// 	} else {
			// 		if rawMove == castle {
			// 			atomic.AddUint32(&data.Castling.Side.Same, 1)
			// 		} else {
			// 			atomic.AddUint32(&data.Castling.Side.Opposite, 1)
			// 		}
			// 	}
			// }

			// if ply < 10 {
			// 	openingPtr = OpeningStats(openingPtr, rawMove)
			// }

			//BranchingFactor
			board := dragontoothmg.ParseFen(gamePtr.Board.Fen())
			branchingFactor := float64(len(board.GenerateLegalMoves()))
			stats.BranchingFactor[ply] += branchingFactor

			//MaterialCount
			count, diff := MaterialCount(gamePtr.Board)
			stats.MaterialCount[ply] = float64(count)
			stats.MaterialDiff[ply] = float64(diff)
			if ply == len(Game.Moves)-1 {
				stats.GameEndMaterialCount[ply] = float64(count)
				stats.GameEndMaterialDiff[ply] = float64(diff)
			}

			//PromotionSquares
			// if move.Promotion != chess.NoPiece {
			// 	data.Heatmaps.PromotionSquares.Count(move.Promotion, move.To)
			// }

			// //EnPassantSquares
			// if gamePtr.Board.EpSquare != chess.NoSquare {
			// 	data.Heatmaps.EnPassantSquares.Count(gamePtr.Parent.Board.Piece[move.From], gamePtr.Board.EpSquare)
			// }
		}

		//results
		// switch Game.Moves[len(Game.Moves)-1] {
		// case "1-0":
		// 	atomic.AddUint32(&data.Results.White, 1)
		// case "0-1":
		// 	atomic.AddUint32(&data.Results.Black, 1)
		// case "1/2-1/2":
		// 	atomic.AddUint32(&data.Results.Draw, 1)
		// default:
		// 	atomic.AddUint32(&data.Results.NA, 1)
		// }

		//ratings
		// if elo, ok := Game.PgnGame.Tags["WhiteElo"]; ok {
		// 	eloNum64, _ := strconv.Atoi(elo)
		// 	eloNum := uint32(eloNum64)

		// 	if eloNum < atomic.LoadUint32(&data.Ratings.Min) {
		// 		atomic.StoreUint32(&data.Ratings.Min, eloNum)
		// 	}

		// 	if eloNum > atomic.LoadUint32(&data.Ratings.Max) {
		// 		atomic.StoreUint32(&data.Ratings.Max, eloNum)
		// 	}
		// }

		// if elo, ok := Game.PgnGame.Tags["BlackElo"]; ok {
		// 	eloNum64, _ := strconv.Atoi(elo)
		// 	eloNum := uint32(eloNum64)

		// 	if eloNum < atomic.LoadUint32(&data.Ratings.Min) {
		// 		atomic.StoreUint32(&data.Ratings.Min, eloNum)
		// 	}

		// 	if eloNum > atomic.LoadUint32(&data.Ratings.Max) {
		// 		atomic.StoreUint32(&data.Ratings.Max, eloNum)
		// 	}
		// }

		//dates
		// if date, ok := Game.PgnGame.Tags["Date"]; ok {
		// 	year64, _ := strconv.Atoi(date[:4])
		// 	year := uint32(year64)

		// 	if year < atomic.LoadUint32(&data.Dates.Min) {
		// 		atomic.StoreUint32(&data.Dates.Min, year)
		// 	}

		// 	if year > atomic.LoadUint32(&data.Dates.Max) {
		// 		atomic.StoreUint32(&data.Dates.Max, year)
		// 	}
		// }

		//GameLengths
		stats.GameLengths[ply] = 1

		//GamesEndingWith
		// check, mate := lastPosition.Board.IsCheckOrMate()
		// //check
		// if check && !mate {
		// 	atomic.AddUint32(&data.GamesEndingWith.Check, 1)
		// }
		// //mate
		// if check && mate {
		// 	atomic.AddUint32(&data.GamesEndingWith.Mate, 1)
		// }
		// //stalemate
		// if !check && mate {
		// 	atomic.AddUint32(&data.GamesEndingWith.Stalemate, 1)
		// }

		gs <- stats
	}
}
