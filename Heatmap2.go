package main

import "github.com/malbrecht/chess"

type HeatSquare map[string]uint64

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

type Heatmap [64]HeatSquare

func NewHeatmap() *Heatmap {
	var hm Heatmap
	for i := 0; i < 64; i++ {
		hm[i] = make(HeatSquare)
	}
	return &hm
}

func (hm *Heatmap) Add(add *Heatmap) {
	for i, square := range add {
		for piece := range square {
			hm[i][piece] += add[i][piece]
		}
	}
}
