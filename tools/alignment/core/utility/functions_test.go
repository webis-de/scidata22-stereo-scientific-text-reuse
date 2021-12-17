/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package utility

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		s []uint64
		e uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test for true",
			args{
				s: []uint64{0, 1, 2, 3, 4, 5},
				e: 0,
			},
			true,
		},
		{
			"test for false",
			args{
				s: []uint64{0, 1, 2, 3, 4, 5},
				e: 6,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxSlice(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"[]float64 test positive",
			args{[]float64{1, 5, 4, 8, 5, 3, 7, 2, 6, 3}},
			float64(8),
			false,
		},
		{
			"[]float64 test negative ",
			args{[]float64{-1, -5, -4, -8, -5, -3, -7, -2, -6, -3}},
			float64(-1),
			false,
		},
		{
			"[]float64 mixed",
			args{[]float64{-1, 5, -4, 8, -5, 3, -7, 2, -6, 3}},
			float64(8),
			false,
		},
		{
			"test error",
			args{[]string{"this", "is", "a", "test", "string"}},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MaxSlice(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("MaxSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MaxSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgMaxSlice(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"[]float64 test positive",
			args{[]float64{1, 5, 4, 8, 5, 3, 7, 2, 6, 3}},
			3,
			false,
		},
		{
			"[]float64 test negative ",
			args{[]float64{-1, -5, -4, -8, -5, -3, -7, -2, -6, -3}},
			0,
			false,
		},
		{
			"[]float64 mixed",
			args{[]float64{-1, 5, -4, 8, -5, 3, -7, 2, -6, 3}},
			3,
			false,
		},
		{
			"test error",
			args{[]string{"this", "is", "a", "test", "string"}},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArgMaxSlice(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArgMaxSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArgMaxSlice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
