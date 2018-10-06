package main

import (
	"sync/atomic"

	"github.com/malbrecht/chess"
)

//WB is White/Black counts
type WB struct {
	W uint32
	B uint32
}

//Heatsquare is a square in the Heatmap
type Heatsquare struct {
	All WB
	K   WB
	Q   WB
	R   WB
	B   WB
	N   WB
	P   WB
}

//Heatmap is a collection of HeatSquares in the shape of a chess board
type Heatmap [64]Heatsquare

//Count increments the counts of HeatSquares depending on piece and square
func (hm *Heatmap) Count(piece chess.Piece, sq chess.Sq) {
	index := (7-sq.Rank())*8 + sq.File()
	hmPtr := &hm[index]

	switch piece {
	case chess.WP:
		atomic.AddUint32(&hmPtr.P.W, 1)
	case chess.BP:
		atomic.AddUint32(&hmPtr.P.B, 1)
	case chess.WN:
		atomic.AddUint32(&hmPtr.N.W, 1)
	case chess.BN:
		atomic.AddUint32(&hmPtr.N.B, 1)
	case chess.WB:
		atomic.AddUint32(&hmPtr.B.W, 1)
	case chess.BB:
		atomic.AddUint32(&hmPtr.B.B, 1)
	case chess.WR:
		atomic.AddUint32(&hmPtr.R.W, 1)
	case chess.BR:
		atomic.AddUint32(&hmPtr.R.B, 1)
	case chess.WQ:
		atomic.AddUint32(&hmPtr.Q.W, 1)
	case chess.BQ:
		atomic.AddUint32(&hmPtr.Q.B, 1)
	case chess.WK:
		atomic.AddUint32(&hmPtr.K.W, 1)
	case chess.BK:
		atomic.AddUint32(&hmPtr.K.B, 1)
	}

	switch piece.Color() {
	case chess.White:
		atomic.AddUint32(&hmPtr.All.W, 1)
	case chess.Black:
		atomic.AddUint32(&hmPtr.All.B, 1)
	}
}
