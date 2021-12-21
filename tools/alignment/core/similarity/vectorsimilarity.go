/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package similarity

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/errors"
)

type IdentityVectorSimilarity struct{}

type CosineVectorSimilarity struct {
	threshold float64 `json:"threshold"`
}

func NewCosineVectorSimilarity(threshold float64) *CosineVectorSimilarity {
	return &CosineVectorSimilarity{threshold: threshold}
}

func (s CosineVectorSimilarity) ComputeScore(source interface{}, target interface{}) (float64, error) {
	switch t1 := source.(type) {
	case mat.Vector:
		switch t2 := target.(type) {
		case mat.Vector:
			return mat.Dot(source.(mat.Vector), target.(mat.Vector)) / (mat.Norm(source.(mat.Vector), 2) * mat.Norm(target.(mat.Vector), 2)), nil
		default:
			return 0, errors.NewTypeError(t2)
		}

	default:
		return 0, errors.NewTypeError(t1)
	}

}

func (s CosineVectorSimilarity) ComputeBool(source interface{}, target interface{}) (bool, error) {
	value, err := s.ComputeScore(source, target)
	return value >= s.threshold, err
}
