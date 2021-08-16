package core

import (
	"encoding/json"

	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/pgn"
)

type FromTo struct {
	From chess.Sq
	To   chess.Sq
}

type FromTos map[FromTo]int

func (ft FromTos) MarshalJSON() ([]byte, error) {
	r := make(map[string]int)

	for k, v := range ft {
		key := k.From.String() + "-" + k.To.String()
		val := v
		r[key] = val
	}

	return json.Marshal(r)
}

type PieceTracker struct {
	PieceMoves map[chess.Sq]FromTos
	squareMap  map[chess.Sq]chess.Sq
}

func NewPieceTracker() *PieceTracker {
	sqmap := make(map[chess.Sq]chess.Sq)
	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sqmap[chess.Square(file, rank)] = chess.NoSquare
		}
	}

	return &PieceTracker{
		PieceMoves: make(map[chess.Sq]FromTos),
		squareMap:  sqmap,
	}
}

func (pt PieceTracker) Track(ptr *pgn.Node) {
	move := ptr.Move

	if move == chess.NullMove {
		return
	}

	if pt.squareMap[move.From] == chess.NoSquare {
		pt.squareMap[move.To] = move.From
	} else {
		pt.squareMap[move.To] = pt.squareMap[move.From]
	}

	origin := pt.squareMap[move.To]
	fromto := FromTo{move.From, move.To}

	if pt.PieceMoves[origin] == nil {
		pt.PieceMoves[origin] = make(FromTos)
	}
	pt.PieceMoves[origin][fromto]++

	if move.Promotion != chess.NoPiece {
		pt.squareMap[move.To] = move.To
	}
}

func (pt PieceTracker) Add(ptr *PieceTracker) {
	for sq, fromtos := range ptr.PieceMoves {
		for fromto, count := range fromtos {
			if pt.PieceMoves[sq] == nil {
				pt.PieceMoves[sq] = make(FromTos)
			}
			pt.PieceMoves[sq][fromto] += count
		}
	}
}

func (pt PieceTracker) MarshalJSON() ([]byte, error) {
	jz := make(map[string]FromTos)
	for k, v := range pt.PieceMoves {
		jz[k.String()] = v
	}
	return json.Marshal(jz)
}
