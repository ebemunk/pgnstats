package main

//MinMax is a min and max uint32
type MinMax struct {
	Min uint32
	Max uint32
}

type Result struct {
	TotalGames uint32
	Ratings    MinMax
	Dates      MinMax
	Results    struct {
		White uint32
		Black uint32
		Draw  uint32
		NA    uint32
	}
	Castling struct {
		White struct {
			Kingside  uint32
			Queenside uint32
		}
		Black struct {
			Kingside  uint32
			Queenside uint32
		}
		Side struct {
			Same     uint32
			Opposite uint32
		}
	}
	GamesEndingWith struct {
		Check uint32
		Mate  uint32
	}
	GameLengths     PlyMap
	MaterialCount   PlyMap
	MaterialDiff    PlyMap
	BranchingFactor PlyMap
	Heatmaps        struct {
		SquareUtilization Heatmap
		MoveSquares       Heatmap
		CaptureSquares    Heatmap
		CheckSquares      Heatmap
	}
	Openings *OpeningMove
}

func NewResult() *Result {
	return &Result{
		Openings: &OpeningMove{
			San:      "start",
			Children: make([]*OpeningMove, 0),
		},
		GameLengths:   PlyMap{},
		MaterialCount: PlyMap{},
		MaterialDiff:  PlyMap{},
		Ratings: MinMax{
			Min: 3000,
			Max: 0,
		},
		Dates: MinMax{
			Min: 3000,
			Max: 0,
		},
		BranchingFactor: PlyMap{},
	}

}
