/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package reducing

import (
	. "picapica.org/picapica4/alignment/core/types"
)

type Reducer interface {
	Process([]*Seed) ([]*Seed, error)
}

// IdentityReducer
// Returns the set of seeds without modifications
type IdentityReducer struct{}

func NewIdentityReducer() *IdentityReducer {
	return &IdentityReducer{}
}

func (p *IdentityReducer) Process(seeds []*Seed) ([]*Seed, error) {
	return seeds, nil
}

// PassageLengthReducer
// Removes all seeds from a set of seeds that have a span of smaller than minLength elements (i.e. characters)
type PassageLengthReducer struct {
	MinLength uint32 `json:"min_length"`
}

func NewPassageLengthReducer(minLength uint32) *PassageLengthReducer {
	return &PassageLengthReducer{MinLength: minLength}
}

func (p *PassageLengthReducer) Process(seeds []*Seed) ([]*Seed, error) {
	result := make([]*Seed, 0, 0)
	for _, s := range seeds {
		if s.SourceSpan.Length() >= p.MinLength && s.SourceSpan.Length() >= p.MinLength {
			result = append(result, s)
		}
	}
	return result, nil
}
