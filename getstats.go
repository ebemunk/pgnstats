package main

import (
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/dylhunn/dragontoothmg"
	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/pgn"
)

//GetStats collects statistics from a single game
func GetStats(Game *pgn.Game, openingsPtr *OpeningMove, filterPlayer string) *GameStats {
	oPtr := openingsPtr

	stats := NewGameStats()
	stats.Total = 1

	ply := -1
	var firstCapture = false

	if Game.Tags["White"] == filterPlayer {
		stats.Color = "w"
	} else if Game.Tags["Black"] == filterPlayer {
		stats.Color = "b"
	} else if filterPlayer != "" {
		return nil
	}

	for gamePtr := Game.Root; gamePtr != nil; gamePtr = gamePtr.Next {
		ply++

		if ply == 0 {
			// start position does not have a valid Move
			continue
		}

		var isFilteredPlayersMove bool
		if filterPlayer == "" {
			isFilteredPlayersMove = true
		} else if Game.Tags["White"] == filterPlayer && gamePtr.Board.SideToMove == chess.Black {
			// SideToMove is the side *after* the move has been played
			isFilteredPlayersMove = true
		} else if Game.Tags["Black"] == filterPlayer && gamePtr.Board.SideToMove == chess.White {
			// SideToMove is the side *after* the move has been played
			isFilteredPlayersMove = true
		} else {
			isFilteredPlayersMove = false
		}

		move := gamePtr.Move
		isLastMove := gamePtr.Next == nil
		fen := gamePtr.Board.Fen()

		//Openings
		if ply > 0 && ply < 10 {
			atomic.AddUint32(&oPtr.Count, 1)
			oPtr = OpeningStats(oPtr, gamePtr.Move.San(gamePtr.Parent.Board))
		}

		//BranchingFactor
		board := dragontoothmg.ParseFen(fen)
		branchingFactor := float64(len(board.GenerateLegalMoves()))
		stats.BranchingFactor[ply] += branchingFactor

		//MaterialCount
		count, diff := MaterialCount(gamePtr.Board)
		stats.MaterialCount[ply] = float64(count)
		stats.MaterialDiff[ply] = float64(diff)
		if isLastMove {
			stats.GameEndMaterialCount[ply] = float64(count)
			stats.GameEndMaterialDiff[ply] = float64(diff)
		}

		if !isFilteredPlayersMove {
			continue
		}

		//Heatmaps
		HeatmapStats(stats, gamePtr, isLastMove)

		if ply > 0 && !firstCapture {
			firstCapture = FirstBlood(&stats.Heatmaps.FirstBlood, gamePtr)
		}

		//PromotionSquares
		if move.Promotion != chess.NoPiece {
			stats.Heatmaps.PromotionSquares.Count(move.Promotion, move.To)
		}

		//EnPassantSquares
		if gamePtr.Board.EpSquare != chess.NoSquare {
			stats.Heatmaps.EnPassantSquares.Count(gamePtr.Parent.Board.Piece[move.From], gamePtr.Board.EpSquare)
		}

		//TrackMoves
		stats.Trax.Track(gamePtr)

		stats.Positions[fen]++
		boardEqualityRegexp, _ := regexp.Compile(`.+ [bw] (-|[KQkq]+) (-|[a-h]\d)`)
		uniquePos := strings.Join(boardEqualityRegexp.FindAllString(fen, -1), "")
		stats.UniquePositions[uniquePos]++
	}

	for _, v := range stats.Positions {
		stats.TotalPositions += v
	}
	stats.TotalUniquePositions = len(stats.UniquePositions)

	//Ratings
	if elo, ok := Game.Tags["WhiteElo"]; ok {
		stats.Ratings[elo] = 1
	}

	if elo, ok := Game.Tags["BlackElo"]; ok {
		stats.Ratings[elo] = 1
	}

	//Years
	if date, ok := Game.Tags["Date"]; ok {
		year64, _ := strconv.Atoi(date[:4])
		year := strconv.Itoa(year64)

		stats.Years[year] = 1
	} else {
		if date, ok := Game.Tags["UTCDate"]; ok {
			year64, _ := strconv.Atoi(date[:4])
			year := strconv.Itoa(year64)

			stats.Years[year] = 1
		}
	}

	//GameLengths
	stats.GameLengths[ply] = 1

	return stats
}
