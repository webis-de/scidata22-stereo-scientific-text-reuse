/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package combinating

import (
	"picapica.org/picapica4/alignment/core/errors"
	"strings"
)

type StringValueCombinator struct {
	sep string `json:"separator"`
}

func NewStringValueCombinator(sep string) *StringValueCombinator {
	return &StringValueCombinator{sep: sep}
}

func (c StringValueCombinator) Combine(input []interface{}) (interface{}, error) {
	seq := make([]string, len(input), len(input))
	for i, elem := range input {
		switch t := elem.(type) {
		case string:
			seq[i] = elem.(string)
		default:
			return seq, errors.NewTypeError(t)
		}
	}
	return strings.Join(seq, c.sep), nil
}
