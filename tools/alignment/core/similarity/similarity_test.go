/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package similarity

import (
	"gonum.org/v1/gonum/mat"
	"reflect"
	"testing"
)

func TestIdentitySimilarity_ComputeBool(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"inequality single letter string",
			args{"A", "B"},
			false,
			false,
		},
		{
			"equality single letter string",
			args{"A", "A"},
			true,
			false,
		},
		{
			"inequality multiple letter string",
			args{"Test string A", "Test string B"},
			false,
			false,
		},
		{
			"equality multiple letter string",
			args{"Test string A", "Test string A"},
			true,
			false,
		},
		{
			"inequality vector",
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 5, 4})},
			false,
			false,
		},
		{
			"equality vector",
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IdentitySimilarity{}
			got, err := s.ComputeBool(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComputeBool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentitySimilarity_ComputeScore(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			"string inequality",
			args{"Test string A", "Test string B"},
			0,
			false,
		},
		{
			"string equality",
			args{"Test string A", "Test string A"},
			1,
			false,
		},
		{
			"vector equality",
			args{mat.NewVecDense(4, []float64{1, 2, 3, 4}), mat.NewVecDense(4, []float64{4, 3, 2, 1})},
			0,
			false,
		},
		{
			"vector inequality",
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IdentitySimilarity{}
			got, err := s.ComputeScore(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComputeScore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdentitySimilarity(t *testing.T) {
	tests := []struct {
		name string
		want *IdentitySimilarity
	}{
		{
			"constructor test",
			&IdentitySimilarity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentitySimilarity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentitySimilarity(0) = %v, want %v", got, tt.want)
			}
		})
	}
}
