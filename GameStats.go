package main

import (
	"encoding/json"
	"sort"
)

type PosMap map[string]int

//PlyMap is a map[int]float64
type PlyMap map[int]float64

//MarshalJSON marshals to json
func (m PlyMap) MarshalJSON() ([]byte, error) {
	if len(m) < 1 {
		return json.Marshal([]int{})
	}

	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	max := keys[len(keys)-1]

	sorted := make([]float64, 0, max)
	for i := 0; i <= max; i++ {
		sorted = append(sorted, m[i])
	}

	return json.Marshal(sorted)
}

//Heatmaps are all the Heatmaps we return
type Heatmaps struct {
	SquareUtilization   Heatmap
	MoveSquares         Heatmap
	CaptureSquares      Heatmap
	CheckSquares        Heatmap
	FirstBlood          Heatmap
	PromotionSquares    Heatmap
	EnPassantSquares    Heatmap
	MateSquares         Heatmap
	MateDeliverySquares Heatmap
	StalemateSquares    Heatmap
}

//GameStats is the statistics for games
type GameStats struct {
	Color                string
	Total                uint64
	GameLengths          PlyMap
	BranchingFactor      PlyMap
	MaterialCount        PlyMap
	MaterialDiff         PlyMap
	GameEndMaterialCount PlyMap
	GameEndMaterialDiff  PlyMap
	Years                map[string]int
	Ratings              map[string]int
	Heatmaps             Heatmaps
	Openings             *OpeningMove
	Trax                 PieceTracker
	Positions            PosMap
	TotalPositions       int
	UniquePositions      PosMap
	TotalUniquePositions int
}

//NewGameStats creates new GameStats
func NewGameStats() *GameStats {
	return &GameStats{
		// "" or "w" or "b"
		Color:                "",
		Total:                0,
		GameLengths:          make(map[int]float64),
		BranchingFactor:      make(map[int]float64),
		MaterialCount:        make(map[int]float64),
		MaterialDiff:         make(map[int]float64),
		GameEndMaterialCount: make(map[int]float64),
		GameEndMaterialDiff:  make(map[int]float64),
		Years:                make(map[string]int),
		Ratings:              make(map[string]int),
		Heatmaps: Heatmaps{
			SquareUtilization:   *NewHeatmap(),
			MoveSquares:         *NewHeatmap(),
			CaptureSquares:      *NewHeatmap(),
			CheckSquares:        *NewHeatmap(),
			FirstBlood:          *NewHeatmap(),
			PromotionSquares:    *NewHeatmap(),
			EnPassantSquares:    *NewHeatmap(),
			MateSquares:         *NewHeatmap(),
			MateDeliverySquares: *NewHeatmap(),
			StalemateSquares:    *NewHeatmap(),
		},
		Trax:                 *NewPieceTracker(),
		Positions:            make(PosMap),
		TotalPositions:       0,
		UniquePositions:      make(PosMap),
		TotalUniquePositions: 0,
	}
}

//Add adds GameStats together
func (gs *GameStats) Add(ad *GameStats) {
	for k, v := range ad.GameLengths {
		gs.GameLengths[k] += v
	}

	for k, v := range ad.BranchingFactor {
		gs.BranchingFactor[k] += v
	}

	for k, v := range ad.MaterialCount {
		gs.MaterialCount[k] += v
	}

	for k, v := range ad.MaterialDiff {
		gs.MaterialDiff[k] += v
	}

	for k, v := range ad.GameEndMaterialCount {
		gs.GameEndMaterialCount[k] += v
	}

	for k, v := range ad.GameEndMaterialDiff {
		gs.GameEndMaterialDiff[k] += v
	}

	for k, v := range ad.Years {
		gs.Years[k] += v
	}

	for k, v := range ad.Ratings {
		gs.Ratings[k] += v
	}

	for k, v := range ad.Positions {
		gs.Positions[k] += v
	}

	gs.TotalPositions += ad.TotalPositions
	gs.TotalUniquePositions += ad.TotalUniquePositions

	gs.Heatmaps.SquareUtilization.Add(&ad.Heatmaps.SquareUtilization)
	gs.Heatmaps.MoveSquares.Add(&ad.Heatmaps.MoveSquares)
	gs.Heatmaps.CaptureSquares.Add(&ad.Heatmaps.CaptureSquares)
	gs.Heatmaps.CheckSquares.Add(&ad.Heatmaps.CheckSquares)
	gs.Heatmaps.FirstBlood.Add(&ad.Heatmaps.FirstBlood)
	gs.Heatmaps.PromotionSquares.Add(&ad.Heatmaps.PromotionSquares)
	gs.Heatmaps.EnPassantSquares.Add(&ad.Heatmaps.EnPassantSquares)
	gs.Heatmaps.MateSquares.Add(&ad.Heatmaps.MateSquares)
	gs.Heatmaps.MateDeliverySquares.Add(&ad.Heatmaps.MateDeliverySquares)
	gs.Heatmaps.StalemateSquares.Add(&ad.Heatmaps.StalemateSquares)

	gs.Trax.Add(&ad.Trax)

	gs.Total++
}

//Average averages statistics by total games
func (gs *GameStats) Average() {
	for k, v := range gs.BranchingFactor {
		gs.BranchingFactor[k] = v / float64(gs.Total)
	}

	for k, v := range gs.MaterialCount {
		gs.MaterialCount[k] = v / float64(gs.Total)
	}

	for k, v := range gs.MaterialDiff {
		gs.MaterialDiff[k] = v / float64(gs.Total)
	}

	for k, v := range gs.GameEndMaterialCount {
		gs.GameEndMaterialCount[k] = v / gs.GameLengths[k]
	}

	for k, v := range gs.GameEndMaterialDiff {
		gs.GameEndMaterialDiff[k] = v / gs.GameLengths[k]
	}
}
