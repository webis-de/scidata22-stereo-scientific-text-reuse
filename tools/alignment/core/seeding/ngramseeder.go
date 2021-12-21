/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/errors"
	"picapica.org/picapica4/alignment/core/filtering"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
)

type NGramSeeder struct {
	BaseSeeder
	Combinator combinating.ValueCombinator `json:"combinator"`
	N          int                         `json:"n"`
	Overlap    int                         `json:"overlap"`
}

// Constructor for NGramSeeder. Takes uint32 arguments for n and overlap to ensure proper domain of values, but
// converts them to int internally to make them usable in loops.
func NewNGramSeeder(n uint32, overlap uint32, similarity similarity.Similarity, combinator combinating.ValueCombinator) *NGramSeeder {
	baseSeeder := BaseSeeder{
		PreFilterPipeline:  filtering.NewFilterPipeline(),
		SeedFilterPipeline: filtering.NewFilterPipeline(),
		Similarity:         similarity,
	}
	return &NGramSeeder{BaseSeeder: baseSeeder, Combinator: combinator, N: int(n), Overlap: int(overlap)}
}

func (s *NGramSeeder) Seed(source []*Element, target []*Element) ([]*Seed, error) {
	// Compute seed candidates for source and target sequence
	sourceSeeds, err := s.computeSeedCandidates(source)
	if err != nil {
		return nil, err
	}
	targetSeeds, err := s.computeSeedCandidates(target)
	if err != nil {
		return nil, err
	}
	// Find matching seeds in both sequences according to the specified similarity function
	seeds, err := s.computeSeeds(sourceSeeds, targetSeeds)
	if err != nil {
		return nil, err
	}
	// Return the slice
	return seeds, nil
}

// Iterates over all elements in a sequence, processes them, filters them and adds them to a set of seedCandidates if
// the defined SeedFilterPipeline is applicable and the SeedProcessingPipeline is applied
func (s *NGramSeeder) computeSeedCandidates(seq []*Element) (seedCandidates []*Element, err error) {

	s.PreFilterPipeline.Reduce(seq...)

	var index uint32 = 0
	for i := 0; i+s.N <= len(seq); i += s.N - s.Overlap {
		span := Span{
			Begin: seq[i].Span.Begin,
			End:   seq[i+s.N-1].Span.End,
			Index: index,
		}
		index++

		seedValues := make([]interface{}, 0, 0)
		for j := i; j < i+s.N; j++ {
			switch t := seq[j].Value.(type) {
			case string:
				seedValues = append(seedValues, seq[j].Value.(string))
			case []float32:
				seedValues = append(seedValues, seq[j].Value.([]float32))
			default:
				return seedCandidates, errors.NewTypeError(t)
			}
		}

		//s.SeedFilterPipeline.Reduce(seedValues...)

		if seedValue, err := s.Combinator.Combine(seedValues); err == nil {
			seedCandidates = append(seedCandidates, &Element{
				Span:  &span,
				Value: seedValue,
			})
		} else {
			return seedCandidates, err
		}
	}

	return seedCandidates, nil
}

// Computes matching elements between two sequences by applying the similarity function
// to each combination of elements. Pairs for which the similarity returns true are added
// to the set of seeds. Parallel computation using the worker concurrency pattern.
func (s *NGramSeeder) computeSeeds(sourceSeeds []*Element, targetSeeds []*Element) ([]*Seed, error) {
	var numJobs = len(sourceSeeds) * len(targetSeeds)

	// TODO: implement this parameter in a general conf
	const numWorkers = 8

	jobs := make(chan [2]*Element, numJobs) // Channel for comparison pairs
	results := make(chan *Seed, numJobs)    // Channel for comparison results
	errorChan := make(chan error)           // Channel to collect potential errors

	for w := 1; w <= numWorkers; w++ {
		go func(id int, jobs <-chan [2]*Element, results chan<- *Seed, errorChan chan<- error) {
			for elements := range jobs {
				sim, err := s.Similarity.ComputeBool(elements[0].Value, elements[1].Value)
				if err != nil {
					errorChan <- err
				} else {
					if sim {
						results <- &Seed{
							SourceSpan: elements[0].Span,
							TargetSpan: elements[1].Span,
						}
					} else {
						results <- nil
					}
				}
			}
		}(w, jobs, results, errorChan)
	}
	for _, x := range sourceSeeds {
		for _, y := range targetSeeds {
			jobs <- [2]*Element{x, y}
		}
	}
	close(jobs)

	seeds := make([]*Seed, 0, numJobs)
	for a := 1; a <= numJobs; a++ {
		select {
		case msg := <-results:
			if msg != nil {
				seeds = append(seeds, msg)
			}
		case err := <-errorChan:
			return nil, err
		}
	}
	return seeds, nil
}
