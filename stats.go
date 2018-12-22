package main

import (
	"sync/atomic"

	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/pgn"
)

//FirstBlood counts captures and returns true if one occurred
func FirstBlood(hm *Heatmap, node *pgn.Node) bool {
	to := node.Move.To
	piece := node.Board.Piece[to]
	targetPiece := node.Parent.Board.Piece[to]

	// this avoids castling moves as king's Move.To is always chess.NoPiece
	if piece == chess.NoPiece || targetPiece == chess.NoPiece {
		return false
	}

	hm.Count(piece, to)

	return true
}

//HeatmapStats collects stats for Heatmaps
func HeatmapStats(data *GameStats, node *pgn.Node, lastmove bool) {
	move := node.Move
	piece := node.Board.Piece[node.Move.To]

	king := chess.Piece(node.Board.SideToMove | chess.King)
	rook := chess.Piece(node.Board.SideToMove | chess.Rook)

	if move.From == chess.E1 && move.To == chess.A1 ||
		move.From == chess.E1 && move.To == chess.H1 ||
		move.From == chess.E8 && move.To == chess.A8 ||
		move.From == chess.E8 && move.To == chess.H8 {
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

		data.Heatmaps.SquareUtilization.Count(king, kingMove.To)
		data.Heatmaps.SquareUtilization.Count(rook, rookMove.To)
		data.Heatmaps.MoveSquares.Count(king, kingMove.From)
		data.Heatmaps.MoveSquares.Count(rook, rookMove.From)
	} else {
		data.Heatmaps.SquareUtilization.Count(piece, move.To)
		data.Heatmaps.MoveSquares.Count(piece, move.From)
	}

	if lastmove {
		check, mate := node.Board.IsCheckOrMate()

		if check {
			//if chess.noPiece, it's check by castling - count it as rook check
			if piece == chess.NoPiece {
				data.Heatmaps.CheckSquares.Count(rook, move.To)
			} else {
				data.Heatmaps.CheckSquares.Count(piece, move.To)
			}
		}

		//checkmate
		if check && mate {
			//if chess.noPiece, it's check by castling - count it as rook check
			if piece == chess.NoPiece {
				data.Heatmaps.MateDeliverySquares.Count(rook, move.To)
			} else {
				data.Heatmaps.MateDeliverySquares.Count(piece, move.To)
			}

			enemyKing := chess.Piece(1 - node.Board.SideToMove | chess.King)
			//locate enemy king on the board
			for i := 0; i < 64; i++ {
				if node.Board.Piece[i] == enemyKing {
					data.Heatmaps.MateSquares.Count(enemyKing, chess.Square(i%8, i/8))
				}
			}
		}

		//stalemate
		if !check && mate {
			enemyKing := chess.Piece(1 - node.Board.SideToMove | chess.King)
			//locate enemy king on the board
			for i := 0; i < 64; i++ {
				if node.Board.Piece[i] == enemyKing {
					data.Heatmaps.StalemateSquares.Count(enemyKing, chess.Square(i%8, i/8))
				}
			}
		}
	}

	if node.Board.Piece[node.Move.To] != chess.NoPiece && node.Parent != nil && node.Parent.Board.Piece[node.Move.To] != chess.NoPiece {
		data.Heatmaps.CaptureSquares.Count(piece, move.To)
	}
}

//OpeningStats collects stats for OpeningMoves
func OpeningStats(ptr *OpeningMove, san string) *OpeningMove {
	openingMove := ptr.Find(san)
	if openingMove != nil {
		atomic.AddUint32(&openingMove.Count, 1)
		ptr = openingMove
	} else {
		openingMove = &OpeningMove{
			1, san, make([]*OpeningMove, 0),
		}
		ptr.Children = append(ptr.Children, openingMove)
		ptr = ptr.Children[len(ptr.Children)-1]
	}

	return ptr
}

//MaterialCount returns sum of material and difference (white - black)
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
