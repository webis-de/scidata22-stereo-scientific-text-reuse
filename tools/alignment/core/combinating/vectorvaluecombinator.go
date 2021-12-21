/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package combinating

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/errors"
)

type VectorValueCombinator struct{}

func NewVectorValueCombinator() VectorValueCombinator {
	return VectorValueCombinator{}
}

func (c VectorValueCombinator) Combine(input []interface{}) (interface{}, error) {
	seq := make([][]float32, len(input), len(input))
	totalLength := 0
	for i, elem := range input {
		switch t := elem.(type) {
		case []float32:
			seq[i] = elem.([]float32)
			totalLength += len(elem.([]float32))
		default:
			return nil, errors.NewTypeError(t)
		}
	}
	data := make([]float64, totalLength, totalLength)
	index := 0
	for _, vec := range seq {
		for i := 0; i < len(vec); i++ {
			data[index] = float64(vec[i])
			index++
		}
	}
	return mat.NewVecDense(len(data), data), nil
}
