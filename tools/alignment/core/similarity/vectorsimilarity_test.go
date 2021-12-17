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

func TestCosineVectorSimilarity_ComputeBool(t *testing.T) {
	type fields struct {
		threshold float64
	}
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			"Test for inequality",
			fields{1},
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 5, 4})},
			false,
			false,
		},
		{
			"Test for equality",
			fields{1},
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			true,
			false,
		},
		{
			"Test for left side error",
			fields{1},
			args{1, mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			false,
			true,
		},
		{
			"Test for right side error",
			fields{1},
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), 5},
			false,
			true,
		},
		{
			"Test for both sided error",
			fields{1},
			args{"", 5},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CosineVectorSimilarity{
				threshold: tt.fields.threshold,
			}
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

func TestCosineVectorSimilarity_ComputeScore(t *testing.T) {
	type fields struct {
		threshold float64
	}
	type args struct {
		source interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			"Test for inequality",
			fields{1},
			args{mat.NewVecDense(4, []float64{1, 2, 3, 4}), mat.NewVecDense(4, []float64{4, 3, 2, 1})},
			float64(2) / float64(3),
			false,
		},
		{
			"Test for equality",
			fields{1},
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			1,
			false,
		},
		{
			"Test for left side error",
			fields{1},
			args{1, mat.NewVecDense(5, []float64{1, 2, 3, 4, 5})},
			0,
			true,
		},
		{
			"Test for right side error",
			fields{1},
			args{mat.NewVecDense(5, []float64{1, 2, 3, 4, 5}), 5},
			0,
			true,
		},
		{
			"Test for both sided error",
			fields{1},
			args{"", 5},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CosineVectorSimilarity{
				threshold: tt.fields.threshold,
			}
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

func TestNewCosineVectorSimilarity(t *testing.T) {
	type args struct {
		threshold float64
	}
	tests := []struct {
		name string
		args args
		want *CosineVectorSimilarity
	}{
		{
			"constructor test",
			args{
				threshold: 0.9,
			},
			&CosineVectorSimilarity{
				threshold: 0.9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCosineVectorSimilarity(tt.args.threshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCosineVectorSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}
