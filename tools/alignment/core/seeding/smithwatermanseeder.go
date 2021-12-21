/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
	"picapica.org/picapica4/alignment/core/utility"
)

// Identifies common passages in source and target using the Smith-Waterman algorithm to compute the ideal local alignment.
// Smith, Temple F. & Waterman, Michael S. (1981). "Identification of Common Molecular Subsequences"  10.1016/0022-2836(81)90087-5.
type SmithWatermanSeeder struct {
	BaseSeeder
	MismatchCost float64 `json:"mismatch_cost"`
	GapCost      float64 `json:"gap_cost"`
	MatchAward   float64 `json:"match_award"`
}

func NewSmithWatermanSeeder(similarity similarity.Similarity, mismatchCost float64, gapCost float64, matchAward float64) *SmithWatermanSeeder {
	return &SmithWatermanSeeder{
		BaseSeeder:   NewBaseSeeder(similarity),
		MismatchCost: mismatchCost,
		GapCost:      gapCost,
		MatchAward:   matchAward,
	}
}

func (s *SmithWatermanSeeder) Seed(source []*Element, target []*Element) (seeds []*Seed, err error) {
	// Align both sequences
	sourceSeq, targetSeq := s.align(source, target)

	// Iterate over both sequences and find span of complete sequence
	var (
		beginSource uint32 = math.MaxUint32
		endSource   uint32 = 0
		indexSource uint32 = 0
	)
	for i := range sourceSeq {
		if beginSource > sourceSeq[i].Span.Begin {
			beginSource = sourceSeq[i].Span.Begin
			indexSource = sourceSeq[i].Span.Index
		}
		if endSource < sourceSeq[i].Span.End {
			endSource = sourceSeq[i].Span.End
		}
	}
	var (
		beginTarget uint32 = math.MaxUint32
		endTarget   uint32 = 0
		indexTarget uint32 = 0
	)
	for j := range targetSeq {
		if beginTarget > targetSeq[j].Span.Begin {
			beginTarget = targetSeq[j].Span.Begin
			indexTarget = targetSeq[j].Span.Index
		}
		if endTarget < targetSeq[j].Span.End {
			endTarget = targetSeq[j].Span.End
		}
	}

	return []*Seed{{
		SourceSpan: &Span{Begin: beginSource, End: endSource, Index: indexSource},
		TargetSpan: &Span{Begin: beginTarget, End: endTarget, Index: indexTarget},
	}}, nil
}

func (s *SmithWatermanSeeder) score(source *Element, target *Element) float64 {
	if similar, err := s.Similarity.ComputeBool(source.Value, target.Value); similar && err == nil {
		return s.MatchAward
	} else {
		return s.MismatchCost
	}
}

func (s *SmithWatermanSeeder) align(source []*Element, target []*Element) ([]*Element, []*Element) {
	// Populate the score matrix
	scores, maxCoord := s.scoringTable(source, target)

	// Create empty sequences to hold the results
	sourceSeq := make([]*Element, 0, len(source))
	targetSeq := make([]*Element, 0, len(target))

	// Backtrack optimal sequence using the score matrix, starting from matrix maximum
	i := maxCoord[0]
	j := maxCoord[1]
	for (i > 0) && (j > 0) {
		// Lookahead if local alignment is finished
		if scores.At(i-1, j-1) == 0 {
			sourceSeq = append(sourceSeq, source[i-1])
			targetSeq = append(targetSeq, target[j-1])
			break
		} else {
			//  Find neighboring cell with highest value for next move
			argMax, _ := utility.ArgMaxSlice([]float64{
				scores.At(i-1, j-1), // Cell diagonally above
				scores.At(i-1, j),   // Cell above
				scores.At(i, j-1),   // Cell to the left
			})
			switch argMax {
			// Move diagonally
			case 0:
				sourceSeq = append(sourceSeq, source[i-1])
				targetSeq = append(targetSeq, target[j-1])
				i--
				j--
			// Move up
			case 1:
				sourceSeq = append(sourceSeq, source[i-1])
				targetSeq = append(targetSeq, &Element{
					Span:  target[j-1].Span,
					Value: nil,
				})
				i--
			// Move left
			case 2:
				sourceSeq = append(sourceSeq, &Element{
					Span:  source[i-1].Span,
					Value: nil,
				})
				targetSeq = append(targetSeq, target[j-1])
				j--
			}
		}
	}
	return sourceSeq, targetSeq
}

func (s *SmithWatermanSeeder) scoringTable(source []*Element, target []*Element) (*mat.Dense, [2]int) {
	n := len(source)
	m := len(target)
	var (
		globalMaxVal   float64 = 0
		globalMaxCoord         = [2]int{0, 0}
	)
	// Initialize new zero matrix for scores
	scores := mat.NewDense(n+1, m+1, nil)

	// Fill the rest
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			scoreVec := []float64{0, 0, 0, 0}
			// Cell diagonally above + similarity scores
			scoreVec[0] = scores.At(i-1, j-1) + s.score(source[i-1], target[j-1])
			// Row Maximum + gapLength * gapCost
			for k := 1; k <= i; k++ {
				if comp := scores.At(i-k, j) + s.GapCost*float64(k); comp >= scoreVec[1] {
					scoreVec[1] = comp
				}
			}
			// Column Maximum + gapLength * gapCost
			for l := 1; l <= j; l++ {
				if comp := scores.At(i, j-l) + s.GapCost*float64(l); comp >= scoreVec[2] {
					scoreVec[2] = comp
				}
			}
			if max, err := utility.ArgMaxSlice(scoreVec); err == nil {
				scores.Set(i, j, scoreVec[max])
				if scoreVec[max] > globalMaxVal {
					globalMaxVal = scoreVec[max]
					globalMaxCoord = [2]int{i, j}
				}
			}
		}
	}
	return scores, globalMaxCoord
}
