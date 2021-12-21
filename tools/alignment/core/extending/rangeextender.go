/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	"math"
	. "picapica.org/picapica4/alignment/core/types"
	"sort"
)

// RangeExtender
// Groups spans that are close to each other in both source and target sequence
// Two Spans are considered close if their index difference is smaller than theta
// Efstathios Stamatatos: "Plagiarism detection using stopword n-grams". https://doi.org/10.1002/asi.21630
type RangeExtender struct {
	theta uint32
}

func NewRangeExtender(theta uint32) *RangeExtender {
	return &RangeExtender{theta: theta}
}

// Iterate over all seeds and identifies candidate passages of close spans
// A passage is a group of spans that have an incremental distance in the source document smaller than theta
func (e *RangeExtender) candidatePassages(seeds []*Seed) [][]*Seed {
	switch len(seeds) {
	case 0:
		return [][]*Seed{}
	case 1:
		return [][]*Seed{seeds}
	default:
		passages := make([][]*Seed, 0)
		// Establish correct order for source span comparison
		sort.Slice(seeds, func(i, j int) bool {
			return seeds[i].SourceSpan.Index < seeds[j].SourceSpan.Index
		})
		group := make([]*Seed, 0)
		// Start with first element
		group = append(group, seeds[0])
		// Iterate over the rest, starting new groups if necessary
		for i := 1; i < len(seeds); i++ {
			// Compare each element to its predecessor
			// If distance > theta, start a new group, otherwise append to existing one
			if uint32(math.Abs(float64(seeds[i].SourceSpan.Index)-float64(seeds[i-1].SourceSpan.Index))) > e.theta {
				passages = append(passages, group)
				group = make([]*Seed, 0)
				group = append(group, seeds[i])
			} else {
				group = append(group, seeds[i])
			}
		}
		// Flush remaining group if needed
		if len(group) > 0 {
			passages = append(passages, group)
		}
		return passages
	}
}

// Iterate over all candidate passages and splits up them if necessary
// A group is split up if the distance between two passages is higher than theta in the target document
// Assumes that seeds in each group are ordered ascending by target index
func (e *RangeExtender) splitCandidates(candidateGroups [][]*Seed) [][]*Seed {
	if len(candidateGroups) == 0 {
		return candidateGroups
	} else {
		passages := make([][]*Seed, 0)
		// Iterate over all groups
		for _, candidateGroup := range candidateGroups {
			// Restore order within the group, ascending by target index
			sort.Slice(candidateGroup, func(i, j int) bool {
				return candidateGroup[i].TargetSpan.Index < candidateGroup[j].TargetSpan.Index
			})
			group := make([]*Seed, 0)
			// Start with first element in candidateGroup
			group = append(group, candidateGroup[0])
			// Iterate over the rest in candidateGroup, starting new split groups if necessary
			for i := 1; i < len(candidateGroup); i++ {
				// Compare each element to its predecessor
				// If distance > theta, start a new split group, otherwise append to existing one
				if uint32(math.Abs(float64(candidateGroup[i].TargetSpan.Index)-float64(candidateGroup[i-1].TargetSpan.Index))) > e.theta {
					passages = append(passages, group)
					group = make([]*Seed, 0)
					group = append(group, candidateGroup[i])
				} else {
					group = append(group, candidateGroup[i])
				}
			}
			// Flush remaining split group if needed
			if len(group) > 0 {
				passages = append(passages, group)
			}
		}
		return passages
	}
}

// Concatenates groups of seeds into single seeds
func (e *RangeExtender) concatenatePassage(sourcePassages [][]*Seed) (passages []*Seed) {
	for _, passage := range sourcePassages {
		// Gets the start and end of the source span of the passage
		sort.Slice(passage, func(i, j int) bool {
			return passage[i].SourceSpan.Index <= passage[j].SourceSpan.Index
		})
		sourceStart := passage[0]
		sourceEnd := passage[len(passage)-1]

		// Gets the start and end of the target span of the passage
		sort.Slice(passage, func(i, j int) bool {
			return passage[i].TargetSpan.Index <= passage[j].TargetSpan.Index
		})
		targetStart := passage[0]
		targetEnd := passage[len(passage)-1]

		// Creates a new seed with the identified ranges and adds it to the result
		passages = append(passages, &Seed{
			SourceSpan: &Span{
				Begin: sourceStart.SourceSpan.Begin,
				End:   sourceEnd.SourceSpan.End,
				Index: sourceStart.SourceSpan.Index,
			},
			TargetSpan: &Span{
				Begin: targetStart.TargetSpan.Begin,
				End:   targetEnd.TargetSpan.End,
				Index: targetStart.TargetSpan.Index,
			},
		})
	}
	return
}

func (e *RangeExtender) Extend(seeds []*Seed) ([]*Seed, error) {
	candidates := e.candidatePassages(seeds)
	candidates = e.splitCandidates(candidates)
	passages := e.concatenatePassage(candidates)
	return passages, nil
}
