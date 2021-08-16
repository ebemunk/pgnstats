package core

import "github.com/malbrecht/chess"

//HeatSquare is a square on a chess board
type HeatSquare map[string]uint64

//Add increments the piece count of the square
func (hs HeatSquare) Add(piece chess.Piece) {
	var key string

	switch piece {
	case chess.WK:
		key = "K"
	case chess.BK:
		key = "k"
	case chess.WQ:
		key = "Q"
	case chess.BQ:
		key = "q"
	case chess.WR:
		key = "R"
	case chess.BR:
		key = "r"
	case chess.WB:
		key = "B"
	case chess.BB:
		key = "b"
	case chess.WN:
		key = "N"
	case chess.BN:
		key = "n"
	case chess.WP:
		key = "P"
	case chess.BP:
		key = "p"
	}

	hs[key]++
}

//Heatmap is a chessboard made up of 64 HeatSquares
type Heatmap [64]HeatSquare

//NewHeatmap returns an initialized Heatmap
func NewHeatmap() *Heatmap {
	var hm Heatmap
	for i := 0; i < 64; i++ {
		hm[i] = make(HeatSquare)
	}
	return &hm
}

//Add adds two Heatmaps together
func (hm *Heatmap) Add(add *Heatmap) {
	for i, square := range add {
		for piece := range square {
			hm[i][piece] += add[i][piece]
		}
	}
}

//Count increments the value for a given square and piece
func (hm *Heatmap) Count(piece chess.Piece, square chess.Sq) {
	i := (7-square.Rank())*8 + square.File()
	hm[i].Add(piece)
}
