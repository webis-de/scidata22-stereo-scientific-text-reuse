/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package similarity

import (
	"reflect"
	"testing"
)

func TestNewLevenshteinStringSimilarity(t *testing.T) {
	type args struct {
		threshold float64
	}
	tests := []struct {
		name string
		args args
		want *LevenshteinStringSimilarity
	}{
		{
			"constructor test",
			args{0.9},
			&LevenshteinStringSimilarity{0.9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLevenshteinStringSimilarity(tt.args.threshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLevenshteinStringSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevenshteinStringSimilarity_ComputeScore(t *testing.T) {
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
			"inequal strings",
			fields{0},
			args{"kitten", "sitting"},
			3,
			false,
		},
		{
			"equal strings",
			fields{0},
			args{"kitten", "kitten"},
			0,
			false,
		},
		{
			"one-sided type error",
			fields{5},
			args{"kitten", 0},
			0,
			true,
		},
		{
			"two-sided type error",
			fields{5},
			args{0, 0},
			0,
			true,
		},
		{
			"empty string",
			fields{5},
			args{"", ""},
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LevenshteinStringSimilarity{
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

func TestLevenshteinStringSimilarity_ComputeBool(t *testing.T) {
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
			"identity threshold",
			fields{0},
			args{"kitten", "sitting"},
			false,
			false,
		},
		{
			"small threshold",
			fields{2},
			args{"kitten", "sitting"},
			false,
			false,
		},
		{
			"large threshold",
			fields{5},
			args{"kitten", "sitting"},
			true,
			false,
		},
		{
			"one-sided type error",
			fields{5},
			args{"kitten", 0},
			true,
			true,
		},
		{
			"two-sided type error",
			fields{5},
			args{0, 0},
			true,
			true,
		},
		{
			"empty strings",
			fields{5},
			args{"", ""},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := LevenshteinStringSimilarity{
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

func TestNewJaccardStringSimilarity(t *testing.T) {
	type args struct {
		threshold float64
		separator string
	}
	tests := []struct {
		name string
		args args
		want *JaccardStringSimilarity
	}{
		{
			"constructor test",
			args{
				threshold: 0.9,
				separator: "",
			},
			&JaccardStringSimilarity{
				threshold: 0.9,
				separator: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJaccardStringSimilarity(tt.args.threshold, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJaccardStringSimilarity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJaccardStringSimilarity_ComputeScore(t *testing.T) {
	type fields struct {
		threshold float64
		separator string
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
			"character-level jaccard similarity",
			fields{
				threshold: 0.9,
				separator: "",
			},
			args{
				source: "abcde",
				target: "cdef",
			},
			0.5,
			false,
		},
		{
			"word-level jaccard similarity",
			fields{
				threshold: 0.9,
				separator: " ",
			},
			args{
				source: "this is a test text",
				target: "this might be a another text",
			},
			0.375,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := JaccardStringSimilarity{
				threshold: tt.fields.threshold,
				separator: tt.fields.separator,
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
