/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
)

// A seeder that checks each element in the source for identity to each element in the target
type IdentitySeeder struct {
	Similarity similarity.Similarity `json:"similarity"`
}

// Constructor for IdentitySeeder. Takes uint32 arguments for n and overlap to ensure proper domain of values, but
// converts them to int internally to make them usable in loops.
func NewIdentitySeeder() *IdentitySeeder {
	return &IdentitySeeder{similarity.NewIdentitySimilarity()}
}

func (s *IdentitySeeder) Seed(source []*Element, target []*Element) ([]*Seed, error) {
	// Compute seed candidates for source and target sequence
	// Find matching seeds in both sequences according to the specified similarity function
	seeds, _ := s.computeSeeds(source, target)
	// Return the slice
	return seeds, nil
}

func (s *IdentitySeeder) computeSeeds(sourceSeeds []*Element, targetSeeds []*Element) ([]*Seed, error) {
	var numJobs = len(sourceSeeds) * len(targetSeeds)

	// TODO: implement this parameter in a general conf
	const numWorkers = 8

	jobs := make(chan [2]*Element, numJobs)
	results := make(chan *Seed, numJobs)
	errorChan := make(chan error) // Channel to collect potential errors

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
