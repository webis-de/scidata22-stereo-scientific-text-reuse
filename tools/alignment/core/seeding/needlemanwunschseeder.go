/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
	"picapica.org/picapica4/alignment/core/utility"
	"sort"
)

// Identifies common passages in source and target using the Needleman-Wunsch algorithm to compute the ideal global alignment.
//  Needleman, Saul B. & Wunsch, Christian D. (1970). "A general method applicable to the search for similarities in the amino acid sequence of two proteins". doi:10.1016/0022-2836(70)90057-4.
type NeedlemanWunschSeeder struct {
	BaseSeeder
	MismatchCost float64 `json:"mismatch_cost"`
	GapCost      float64 `json:"gap_cost"`
	MatchAward   float64 `json:"match_award"`
}

func NewNeedlemanWunschSeeder(similarity similarity.Similarity, mismatchCost float64, gapCost float64, matchAward float64) *NeedlemanWunschSeeder {
	return &NeedlemanWunschSeeder{
		BaseSeeder:   NewBaseSeeder(similarity),
		MismatchCost: mismatchCost,
		GapCost:      gapCost,
		MatchAward:   matchAward,
	}
}

func (s *NeedlemanWunschSeeder) Seed(source []*Element, target []*Element) (seeds []*Seed, err error) {
	// Align both sequences
	sourceSeq, targetSeq := s.align(source, target)

	// Iterate over both sequences and only return those elements which have a match in the other sequence
	for i := range sourceSeq {
		if sourceSeq[i].Value != nil && targetSeq[i].Value != nil {
			seeds = append(seeds, &Seed{
				SourceSpan: sourceSeq[i].Span,
				TargetSpan: targetSeq[i].Span,
			})
		}
	}

	// Restore correct order of the seeds
	sort.Slice(seeds, func(i int, j int) bool {
		return seeds[i].SourceSpan.Index < seeds[j].SourceSpan.Index
	})
	return
}

func (s *NeedlemanWunschSeeder) score(source *Element, target *Element) float64 {
	if similar, err := s.Similarity.ComputeBool(source.Value, target.Value); similar && err == nil {
		return s.MatchAward
	} else {
		return s.MismatchCost
	}
}

func (s *NeedlemanWunschSeeder) align(source []*Element, target []*Element) ([]*Element, []*Element) {
	i := len(source)
	j := len(target)

	// Populate the scoring table
	_, traceback := s.scoringTable(source, target)

	// Create empty sequences to hold the results
	sourceSeq := make([]*Element, 0, i)
	targetSeq := make([]*Element, 0, j)

	// Backtrack optimal sequence using the traceback matrix, starting from bottom right cell
	for (i > 0) && (j > 0) {
		current := (*traceback)[i][j]
		switch current {
		case [2]uint16{uint16(i) - 1, uint16(j) - 1}:
			sourceSeq = append(sourceSeq, source[i-1])
			targetSeq = append(targetSeq, target[j-1])
			i--
			j--
		case [2]uint16{uint16(i) - 1, uint16(j)}:
			sourceSeq = append(sourceSeq, source[i-1])
			targetSeq = append(targetSeq, &Element{
				Span:  target[j-1].Span,
				Value: nil,
			})
			i--
		case [2]uint16{uint16(i), uint16(j) - 1}:
			sourceSeq = append(sourceSeq, &Element{
				Span:  source[i-1].Span,
				Value: nil,
			})
			targetSeq = append(targetSeq, target[j-1])
			j--
		}
	}
	// Finish tracing up to the top left cell
	for i > 0 {
		sourceSeq = append(sourceSeq, source[i-1])
		targetSeq = append(targetSeq, &Element{
			Span:  target[j].Span,
			Value: nil,
		})
		i--
	}
	for j > 0 {
		sourceSeq = append(sourceSeq, &Element{
			Span:  target[i].Span,
			Value: nil,
		})
		targetSeq = append(targetSeq, target[j-1])
		j--
	}
	traceback = nil
	return sourceSeq, targetSeq
}

func (s *NeedlemanWunschSeeder) scoringTable(source []*Element, target []*Element) (*mat.Dense, *[][][2]uint16) {
	n := len(source)
	m := len(target)
	// Initialize new zero matrix for scores
	score := mat.NewDense(n+1, m+1, nil)
	// Initialize new matrix for backtracking
	tracebackMatrix := make([][][2]uint16, n+1, n+1)
	for i := range tracebackMatrix {
		tracebackMatrix[i] = make([][2]uint16, m+1, m+1)
	}
	// Initialize first row and column
	score.Set(0, 0, 0)
	for i := 1; i <= n; i++ {
		score.Set(i, 0, s.GapCost*float64(i))
	}
	for j := 1; j <= m; j++ {
		score.Set(0, j, s.GapCost*float64(j))
	}
	// Fill the rest
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			e := []float64{
				score.At(i-1, j-1) + s.score(source[i-1], target[j-1]), // Cell diagonally above + similarity score
				score.At(i, j-1) + s.GapCost,                           // Cell above + gap cost
				score.At(i-1, j) + s.GapCost,                           // Cell to the left + gap cost
			}
			if max, err := utility.ArgMaxSlice(e); err == nil {
				switch max {
				case 0:
					tracebackMatrix[i][j] = [2]uint16{uint16(i) - 1, uint16(j) - 1}
				case 1:
					tracebackMatrix[i][j] = [2]uint16{uint16(i), uint16(j) - 1}
				case 2:
					tracebackMatrix[i][j] = [2]uint16{uint16(i) - 1, uint16(j)}
				}
				score.Set(i, j, e[max])
			}
		}
	}
	return score, &tracebackMatrix
}
