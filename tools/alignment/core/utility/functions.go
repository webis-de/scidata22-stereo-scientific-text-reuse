/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package utility

import (
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"picapica.org/picapica4/alignment/core/errors"
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
)

// Computes the union of two slices
func Union(a interface{}, b interface{}) []interface{} {
	hash := make(map[interface{}]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		hash[el] = true
	}

	union := make([]interface{}, len(hash))
	i := 0
	for k := range hash {
		union[i] = k
		i++
	}
	return union
}

// Computes the intersection of two slices
func Intersection(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	hash := make(map[interface{}]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		if _, found := hash[el]; found {
			set = append(set, el)
		}
	}
	return set
}

// Checks if a slice contains a value
// TODO: extend for other types
func Contains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Returns the maximum value in a slice
// TODO: extend for other types
func MaxSlice(s interface{}) (interface{}, error) {
	switch t := s.(type) {
	case []float64:
		max := s.([]float64)[0]
		for _, x := range s.([]float64) {
			if x > max {
				max = x
			}
		}
		return max, nil
	default:
		return nil, errors.NewTypeError(t)
	}
}

// Returns the index of the maximum value in a slice
// TODO: extend for other types
func ArgMaxSlice(s interface{}) (int, error) {
	switch t := s.(type) {
	case []float64:
		max := s.([]float64)[0]
		var index int
		for i, x := range s.([]float64) {
			if x > max {
				max = x
				index = i
			}
		}
		return index, nil
	default:
		return 0, errors.NewTypeError(t)
	}
}

// Generates an element sequence of length n with random vectors of length m as values
func RandomVecSeq(n int, m int) []*Element {
	seq := make([]*Element, n, n)
	for i := 0; i < n; i++ {
		data := make([]float64, m, m)
		for j := 0; j < m; j++ {
			data[j] = rand.Float64()
		}
		seq[i] = &Element{
			Span:  &Span{Begin: uint32(i), End: uint32(i + 1), Index: uint32(i)},
			Value: mat.NewVecDense(m, data),
		}
	}
	return seq
}
