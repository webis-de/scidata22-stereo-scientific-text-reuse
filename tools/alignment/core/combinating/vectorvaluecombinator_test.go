/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package combinating

import (
	"gonum.org/v1/gonum/mat"
	"reflect"
	"testing"
)

func TestNewVectorValueCombinator(t *testing.T) {
	tests := []struct {
		name string
		want VectorValueCombinator
	}{
		{
			"constructor test",
			VectorValueCombinator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVectorValueCombinator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVectorValueCombinator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorValueCombinator_Combine(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"normal input",
			args{input: [][]float32{{0, 1}, {2, 3}}},
			mat.NewVecDense(4, []float64{0, 1, 2, 3}),
			false,
		},
		{
			"scalar input",
			args{input: [][]float32{{0}}},
			mat.NewVecDense(1, []float64{0}),
			false,
		},
		{
			"malformed input",
			args{input: []string{"This", "is", "a", "test", "text"}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := VectorValueCombinator{}
			// Create interface slice to get type conversion working correctly
			// https://golang.org/doc/faq#convert_slice_of_interface
			var s []interface{}
			switch tt.args.input.(type) {
			case []string:
				s = make([]interface{}, len(tt.args.input.([]string)))
				for i, v := range tt.args.input.([]string) {
					s[i] = v
				}
			case [][]float32:
				s = make([]interface{}, len(tt.args.input.([][]float32)))
				for i, v := range tt.args.input.([][]float32) {
					s[i] = v
				}
			}
			got, err := c.Combine(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Combine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combine() got = %v, want %v", got, tt.want)
			}
		})
	}
}
