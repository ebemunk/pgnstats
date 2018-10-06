package main

import (
	"strings"
	"sync/atomic"

	"github.com/malbrecht/chess"
)

//HeatmapStats collects stats for Heatmaps
func HeatmapStats(data *Result, move chess.Move, piece chess.Piece, rawMove string) {
	if move.From == chess.E1 && move.To == chess.A1 ||
		move.From == chess.E1 && move.To == chess.H1 ||
		move.From == chess.E8 && move.To == chess.A8 ||
		move.From == chess.E8 && move.To == chess.H8 {
		var color uint8
		var kingMove chess.Move
		var rookMove chess.Move

		if move.To-move.From == 3 {
			kingMove = chess.Move{
				From:      move.From,
				To:        move.From + 2,
				Promotion: chess.NoPiece,
			}
			rookMove = chess.Move{
				From:      move.To,
				To:        move.To - 2,
				Promotion: chess.NoPiece,
			}
		} else {
			kingMove = chess.Move{
				From:      move.From,
				To:        move.From - 2,
				Promotion: chess.NoPiece,
			}
			rookMove = chess.Move{
				From:      move.To,
				To:        move.To + 3,
				Promotion: chess.NoPiece,
			}
		}

		if move.From == chess.E1 {
			color = chess.White
		} else {
			color = chess.Black
		}

		king := chess.Piece(color | chess.King)
		rook := chess.Piece(color | chess.Rook)

		data.Heatmaps.SquareUtilization.Count(king, kingMove.To)
		data.Heatmaps.SquareUtilization.Count(rook, rookMove.To)
		data.Heatmaps.MoveSquares.Count(king, kingMove.From)
		data.Heatmaps.MoveSquares.Count(rook, rookMove.From)
	} else {
		data.Heatmaps.SquareUtilization.Count(piece, move.To)
		data.Heatmaps.MoveSquares.Count(piece, move.From)
	}

	if strings.ContainsRune(rawMove, '+') {
		data.Heatmaps.CheckSquares.Count(piece, move.To)
	}

	if strings.ContainsRune(rawMove, 'x') {
		data.Heatmaps.CaptureSquares.Count(piece, move.To)
	}
}

//OpeningStats collects stats for OpeningMoves
func OpeningStats(ptr *OpeningMove, rawMove string) *OpeningMove {
	openingMove := ptr.Find(rawMove)
	if openingMove != nil {
		atomic.AddUint32(&openingMove.Count, 1)
		ptr = openingMove
	} else {
		openingMove = &OpeningMove{
			1, rawMove, make([]*OpeningMove, 0),
		}
		ptr.Children = append(ptr.Children, openingMove)
		ptr = ptr.Children[len(ptr.Children)-1]
	}

	return ptr
}

//CastlingStats counts the number of kingside and queenside castles by both colors
func CastlingStats(data *Result, rawMove string, sideToMove int) {
	if rawMove == "O-O" {
		if sideToMove == chess.White {
			atomic.AddUint32(&data.Castling.White.Kingside, 1)
		} else {
			atomic.AddUint32(&data.Castling.Black.Kingside, 1)
		}
	} else {
		if sideToMove == chess.White {
			atomic.AddUint32(&data.Castling.White.Queenside, 1)
		} else {
			atomic.AddUint32(&data.Castling.Black.Queenside, 1)
		}
	}
}

func MaterialCount(board *chess.Board) (int, int) {
	countW := 0
	countB := 0
	var countPtr *int

	for _, p := range board.Piece {
		switch chess.Piece(p).Color() {
		case chess.White:
			countPtr = &countW
		case chess.Black:
			countPtr = &countB
		}

		switch chess.Piece(p).Type() {
		case chess.Pawn:
			*countPtr++
		case chess.Knight:
			*countPtr += 3
		case chess.Bishop:
			*countPtr += 3
		case chess.Rook:
			*countPtr += 5
		case chess.Queen:
			*countPtr += 9
		}
	}

	count := countW + countB
	diff := countW - countB

	return count, diff
}
