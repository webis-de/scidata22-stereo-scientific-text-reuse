/*******************************************************************************
 * Last modified: 02.09.21, 10:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package normalizing

import (
	"picapica.org/picapica4/alignment/core/errors"
	"sort"
)

type Normalizer interface {
	Normalize([]interface{}) ([]interface{}, error)
}

// Identity
// Returns the seed without modifcation
type IdentityNormalizer struct{}

func NewIdentityNormalizer() *IdentityNormalizer {
	return &IdentityNormalizer{}
}

func (n *IdentityNormalizer) Normalize(values []interface{}) ([]interface{}, error) {
	return values, nil
}

// SortingNormalizer
// Sorts all values alphabetically
type SortingNormalizer struct{}

func NewSortingNormalizer() *SortingNormalizer {
	return &SortingNormalizer{}
}

func (n *SortingNormalizer) Normalize(values []interface{}) ([]interface{}, error) {
	switch t := values[0].(type) {
	case string:
		sort.Slice(values, func(i, j int) bool {
			return values[i].(string) < values[j].(string)
		})
		return values, nil
	default:
		return values, errors.NewTypeError(t)
	}
}
