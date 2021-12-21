/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"fmt"
	"hash/maphash"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/errors"
	"picapica.org/picapica4/alignment/core/filtering"
	"picapica.org/picapica4/alignment/core/normalizing"
	. "picapica.org/picapica4/alignment/core/types"
	"sync"
)

type HashSeeder struct {
	BaseSeeder
	Normalizer normalizing.Normalizer      `json:"normalizer"`
	Combinator combinating.ValueCombinator `json:"combinator"`
	N          int                         `json:"n"`
	Overlap    int                         `json:"overlap"`
}

// Constructor for NGramSeeder. Takes uint32 arguments for n and overlap to ensure proper domain of values, but
// converts them to int internally to make them usable in loops.
func NewHashSeeder(n uint32, overlap uint32, normalizer normalizing.Normalizer, combinator combinating.ValueCombinator) *HashSeeder {
	baseSeeder := BaseSeeder{
		PreFilterPipeline:  filtering.NewFilterPipeline(),
		SeedFilterPipeline: filtering.NewFilterPipeline(),
	}
	return &HashSeeder{BaseSeeder: baseSeeder, Normalizer: normalizer, Combinator: combinator, N: int(n), Overlap: int(overlap)}
}

func (s *HashSeeder) Seed(source []*Element, target []*Element) ([]*Seed, error) {
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
func (s *HashSeeder) computeSeedCandidates(seq []*Element) (seedCandidates []*Element, err error) {
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
		seedValues, err := s.Normalizer.Normalize(seedValues)
		if err != nil {
			return seedCandidates, err
		}

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

// Computes matching elements between two sequences by comparing the hashes of elements
func (s *HashSeeder) computeSeeds(sourceSeeds []*Element, targetSeeds []*Element) ([]*Seed, error) {
	// Setup the hashmap & all other needed stuff for writing to it
	seedMap := make(map[uint64][2][]*Span)
	hash := maphash.Hash{}
	mapMutex := sync.Mutex{}
	hashMutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(sourceSeeds) + len(targetSeeds))
	// TODO: implement this parameter in a general conf
	const numWorkers = 16

	// Write the source seeds to the first value of the hash map
	var (
		numSourceJobs = len(sourceSeeds)
		sourceJobs    = make(chan *Element, numSourceJobs)
	)
	for w := 1; w <= numWorkers; w++ {
		go func(id int, jobs <-chan *Element) {
			for x := range jobs {
				hashMutex.Lock()
				_, _ = hash.Write([]byte(fmt.Sprintf("%v", x.Value)))
				key := hash.Sum64()
				hash.Reset()
				hashMutex.Unlock()
				mapMutex.Lock()
				if existing, ok := seedMap[key]; ok {
					seedMap[key] = [2][]*Span{append(existing[0], x.Span), existing[1]}
				} else {
					seedMap[key] = [2][]*Span{{x.Span}, existing[1]}
				}
				mapMutex.Unlock()
				wg.Done()
			}
		}(w, sourceJobs)
	}
	for _, x := range sourceSeeds {
		sourceJobs <- x
	}
	close(sourceJobs)

	// Write the target seeds to the second value of the hash map
	var (
		numTargetJobs = len(targetSeeds)
		targetJobs    = make(chan *Element, numTargetJobs)
	)
	for w := 1; w <= numWorkers; w++ {
		go func(id int, jobs <-chan *Element) {
			for y := range jobs {
				hashMutex.Lock()
				_, _ = hash.Write([]byte(fmt.Sprintf("%v", y.Value)))
				key := hash.Sum64()
				hash.Reset()
				hashMutex.Unlock()
				mapMutex.Lock()
				if existing, ok := seedMap[key]; ok {
					seedMap[key] = [2][]*Span{existing[0], append(existing[1], y.Span)}
				} else {
					seedMap[key] = [2][]*Span{existing[0], {y.Span}}
				}
				mapMutex.Unlock()
				wg.Done()
			}
		}(w, targetJobs)
	}
	for _, y := range targetSeeds {
		targetJobs <- y
	}
	close(targetJobs)

	// Turn the hash map into seeds; wherever there are elements in both
	// source and target part of the value, construct seeds from every
	// combination of them
	seeds := make([]*Seed, 0, 0)
	wg.Wait()
	for _, v := range seedMap {
		for _, source := range v[0] {
			for _, target := range v[1] {
				seeds = append(seeds, &Seed{
					SourceSpan: source,
					TargetSpan: target,
				})
			}
		}
	}
	return seeds, nil
}
