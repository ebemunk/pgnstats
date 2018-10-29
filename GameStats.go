package main

import (
	"encoding/json"
	"sort"
)

//PlyMap is a map[int]float64
type PlyMap map[int]float64

//GameStats is the statistics for games
type GameStats struct {
	Total                uint64
	GameLengths          PlyMap
	BranchingFactor      PlyMap
	MaterialCount        PlyMap
	MaterialDiff         PlyMap
	GameEndMaterialCount PlyMap
	GameEndMaterialDiff  PlyMap
	Heatmaps             struct {
		// 		SquareUtilization Heatmap
		// 		MoveSquares       Heatmap
		// 		CaptureSquares    Heatmap
		// 		CheckSquares      Heatmap
		FirstBlood Heatmap
		// 		PromotionSquares  Heatmap
		// 		EnPassantSquares  Heatmap
	}
}

//NewGameStats creates new GameStats
func NewGameStats() *GameStats {
	return &GameStats{
		Total:                0,
		GameLengths:          make(map[int]float64),
		BranchingFactor:      make(map[int]float64),
		MaterialCount:        make(map[int]float64),
		MaterialDiff:         make(map[int]float64),
		GameEndMaterialCount: make(map[int]float64),
		GameEndMaterialDiff:  make(map[int]float64),
		Heatmaps: struct {
			FirstBlood Heatmap
		}{
			FirstBlood: *NewHeatmap(),
		},
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

	gs.Heatmaps.FirstBlood.Add(&ad.Heatmaps.FirstBlood)

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
		gs.GameEndMaterialCount[k] = v / float64(gs.Total)
	}

	for k, v := range gs.GameEndMaterialDiff {
		gs.GameEndMaterialDiff[k] = v / float64(gs.Total)
	}
}

//MarshalJSON marshals to json
func (m PlyMap) MarshalJSON() ([]byte, error) {
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
