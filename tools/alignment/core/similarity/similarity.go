/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package similarity

import (
	"reflect"
)

// The similarity interface aggregates functions to compute how similar two given things are
// An implementation has to provide two functions: one returns a boolean output (are two things
// similar or not?) and the other a float score (how similar are the two, exactly?). In most
// cases, the boolean function is a wrapper around the score function using a defined threshold.
type Similarity interface {
	ComputeBool(source interface{}, target interface{}) (bool, error)
	ComputeScore(source interface{}, target interface{}) (float64, error)
}

type IdentitySimilarity struct{}

func NewIdentitySimilarity() *IdentitySimilarity {
	return &IdentitySimilarity{}
}

func (s IdentitySimilarity) ComputeScore(source interface{}, target interface{}) (float64, error) {
	if reflect.DeepEqual(source, target) {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (s IdentitySimilarity) ComputeBool(source interface{}, target interface{}) (bool, error) {
	value, err := s.ComputeScore(source, target)
	return value == 1, err
}
