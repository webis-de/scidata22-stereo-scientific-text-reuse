/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package similarity

import (
	"picapica.org/picapica4/alignment/core/errors"
	"picapica.org/picapica4/alignment/core/utility"
	"strings"
)

type JaccardStringSimilarity struct {
	threshold float64 `json:"threshold"`
	separator string  `json:"threshold"`
}

func NewJaccardStringSimilarity(threshold float64, separator string) *JaccardStringSimilarity {
	return &JaccardStringSimilarity{threshold: threshold, separator: separator}
}

func (s JaccardStringSimilarity) ComputeScore(source interface{}, target interface{}) (float64, error) {
	switch t1 := source.(type) {
	case string:
		switch t2 := target.(type) {
		case string:
			sourceSeq := strings.Split(source.(string), s.separator)
			targetSeq := strings.Split(target.(string), s.separator)

			intersectionLength := float64(len(utility.Intersection(sourceSeq, targetSeq)))
			unionLength := float64(len(utility.Union(sourceSeq, targetSeq)))

			return intersectionLength / unionLength, nil

		default:
			return 0, errors.NewTypeError(t2)
		}
	default:
		return 0, errors.NewTypeError(t1)
	}
}

func (s JaccardStringSimilarity) ComputeBool(source interface{}, target interface{}) (bool, error) {
	value, err := s.ComputeScore(source, target)
	return value <= s.threshold, err
}

type LevenshteinStringSimilarity struct {
	threshold float64 `json:"threshold"`
}

func NewLevenshteinStringSimilarity(threshold float64) *LevenshteinStringSimilarity {
	return &LevenshteinStringSimilarity{threshold: threshold}
}

func (s LevenshteinStringSimilarity) ComputeScore(source interface{}, target interface{}) (float64, error) {
	switch t1 := source.(type) {
	case string:
		switch t2 := target.(type) {
		case string:
			sourceLen := len(source.(string))
			targetLen := len(target.(string))

			d := make([][]int, sourceLen+1)
			for i := range d {
				d[i] = make([]int, targetLen+1)
			}
			for i := range d {
				d[i][0] = i
			}
			for j := range d[0] {
				d[0][j] = j
			}
			for j := 1; j <= targetLen; j++ {
				for i := 1; i <= sourceLen; i++ {
					if source.(string)[i-1] == target.(string)[j-1] {
						d[i][j] = d[i-1][j-1]
					} else {
						min := d[i-1][j]
						if d[i][j-1] < min {
							min = d[i][j-1]
						}
						if d[i-1][j-1] < min {
							min = d[i-1][j-1]
						}
						d[i][j] = min + 1
					}
				}

			}
			return float64(d[sourceLen][targetLen]), nil
		default:
			return 0, errors.NewTypeError(t2)
		}
	default:
		return 0, errors.NewTypeError(t1)
	}
}

func (s LevenshteinStringSimilarity) ComputeBool(source interface{}, target interface{}) (bool, error) {
	value, err := s.ComputeScore(source, target)
	return value <= s.threshold, err
}
