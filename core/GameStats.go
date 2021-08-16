package core

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/dylhunn/dragontoothmg"
	"github.com/malbrecht/chess"
	"github.com/malbrecht/chess/pgn"
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
	// Years                map[string]int
	Ratings              map[string]int
	Heatmaps             Heatmaps
	Openings             *OpeningMove
	PiecePaths           PieceTracker
	Positions            PosMap
	TotalPositions       int
	UniquePositions      PosMap
	TotalUniquePositions int
}

type PlayerStats struct {
	White GameStats
	Black GameStats
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
		// Years:                make(map[string]int),
		Ratings: make(map[string]int),
		Heatmaps: Heatmaps{
			SquareUtilization:   *NewHeatmap(),
			MoveSquares:         *NewHeatmap(),
			CaptureSquares:      *NewHeatmap(),
			CheckSquares:        *NewHeatmap(),
			FirstBlood:          *NewHeatmap(),
			PromotionSquares:    *NewHeatmap(),
			MateSquares:         *NewHeatmap(),
			MateDeliverySquares: *NewHeatmap(),
			StalemateSquares:    *NewHeatmap(),
		},
		PiecePaths:           *NewPieceTracker(),
		Positions:            make(PosMap),
		TotalPositions:       0,
		UniquePositions:      make(PosMap),
		TotalUniquePositions: 0,
	}
}

//GetStatsFromGame returns GameStats from a single game
func NewGameStatsFromGame(game *pgn.Game, filterPlayer string) *GameStats {
	gs := NewGameStats()
	// 1 because this is for a single game
	gs.Total = 1

	// counter for plies of the game
	ply := -1
	// boolean to track when first capture happens
	var firstCapture = false

	if game.Tags["White"] == filterPlayer {
		gs.Color = "w"
	} else if game.Tags["Black"] == filterPlayer {
		gs.Color = "b"
	} else if filterPlayer != "" {
		return nil
	}

	for gamePtr := game.Root; gamePtr != nil; gamePtr = gamePtr.Next {
		ply++

		// start position does not have a valid Move
		if ply == 0 {
			continue
		}

		var isFilteredPlayersMove bool
		if filterPlayer == "" {
			isFilteredPlayersMove = true
		} else if game.Tags["White"] == filterPlayer && gamePtr.Board.SideToMove == chess.Black {
			// SideToMove is the side *after* the move has been played
			isFilteredPlayersMove = true
		} else if game.Tags["Black"] == filterPlayer && gamePtr.Board.SideToMove == chess.White {
			// SideToMove is the side *after* the move has been played
			isFilteredPlayersMove = true
		} else {
			isFilteredPlayersMove = false
		}

		// move made to reach this position
		move := gamePtr.Move
		isLastMove := gamePtr.Next == nil
		fen := gamePtr.Board.Fen()

		// branching factor is the number of legal moves from a position
		board := dragontoothmg.ParseFen(fen)
		branchingFactor := float64(len(board.GenerateLegalMoves()))
		gs.BranchingFactor[ply] += branchingFactor

		// count material on the board
		count, diff := MaterialCount(gamePtr.Board)
		gs.MaterialCount[ply] = float64(count)
		gs.MaterialDiff[ply] = float64(diff)
		// if last move of a game, also save it in GameEndMaterial
		if isLastMove {
			gs.GameEndMaterialCount[ply] = float64(count)
			gs.GameEndMaterialDiff[ply] = float64(diff)
		}

		// all metrics up until this point is shared
		// the rest of the metrics need to be player-specific, if specified
		if !isFilteredPlayersMove {
			continue
		}

		//Heatmaps
		HeatmapStats(gs, gamePtr, isLastMove)

		if ply > 0 && !firstCapture {
			firstCapture = FirstBlood(&gs.Heatmaps.FirstBlood, gamePtr)
		}

		//PromotionSquares
		if move.Promotion != chess.NoPiece {
			gs.Heatmaps.PromotionSquares.Count(move.Promotion, move.To)
		}

		gs.PiecePaths.Track(gamePtr)

		// keep track of all positions by FEN
		gs.Positions[fen]++
		// unique positions, where game board state is reached regardless of
		// halfmove clock or fullmove counter, which are stipped from the FEN
		// with the regexp below
		boardEqualityRegexp, _ := regexp.Compile(`.+ [bw] (-|[KQkq]+) (-|[a-h]\d)`)
		uniquePos := strings.Join(boardEqualityRegexp.FindAllString(fen, -1), "")
		gs.UniquePositions[uniquePos]++

	}

	// count number of positions encountered
	for _, v := range gs.Positions {
		gs.TotalPositions += v
	}
	// count number of unique positions
	gs.TotalUniquePositions = len(gs.UniquePositions)

	// save the elo ratings for both players
	if elo, ok := game.Tags["WhiteElo"]; ok {
		gs.Ratings[elo] = 1
	}
	if elo, ok := game.Tags["BlackElo"]; ok {
		gs.Ratings[elo] = 1
	}

	gs.GameLengths[ply] = 1

	return gs
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

	// for k, v := range ad.Years {
	// 	gs.Years[k] += v
	// }

	for k, v := range ad.Ratings {
		gs.Ratings[k] += v
	}

	for k, v := range ad.Positions {
		gs.Positions[k] += v
	}

	for k, v := range ad.UniquePositions {
		gs.UniquePositions[k] += v
	}

	gs.TotalPositions += ad.TotalPositions
	gs.TotalUniquePositions += ad.TotalUniquePositions

	gs.Heatmaps.SquareUtilization.Add(&ad.Heatmaps.SquareUtilization)
	gs.Heatmaps.MoveSquares.Add(&ad.Heatmaps.MoveSquares)
	gs.Heatmaps.CaptureSquares.Add(&ad.Heatmaps.CaptureSquares)
	gs.Heatmaps.CheckSquares.Add(&ad.Heatmaps.CheckSquares)
	gs.Heatmaps.FirstBlood.Add(&ad.Heatmaps.FirstBlood)
	gs.Heatmaps.PromotionSquares.Add(&ad.Heatmaps.PromotionSquares)
	// gs.Heatmaps.EnPassantSquares.Add(&ad.Heatmaps.EnPassantSquares)
	gs.Heatmaps.MateSquares.Add(&ad.Heatmaps.MateSquares)
	gs.Heatmaps.MateDeliverySquares.Add(&ad.Heatmaps.MateDeliverySquares)
	gs.Heatmaps.StalemateSquares.Add(&ad.Heatmaps.StalemateSquares)

	gs.PiecePaths.Add(&ad.PiecePaths)

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
