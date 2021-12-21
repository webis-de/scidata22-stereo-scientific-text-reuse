/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
	"picapica.org/picapica4/alignment/core/utility"
	"reflect"
	"sort"
	"testing"
)

func TestNGramSeeder_Seed(t *testing.T) {
	type fields struct {
		similarity similarity.Similarity
		combinator combinating.ValueCombinator
		n          uint32
		overlap    uint32
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantSeeds []*Seed
		wantErr   bool
	}{
		{
			"string values no error",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          2,
				overlap:    1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 16, Index: 3}, Value: "string"},
					{Span: &Span{Begin: 18, End: 21, Index: 4}, Value: "test"},
				},
			},
			[]*Seed{
				{SourceSpan: &Span{Begin: 1, End: 7, Index: 0}, TargetSpan: &Span{Begin: 1, End: 7, Index: 0}},
				{SourceSpan: &Span{Begin: 6, End: 9, Index: 1}, TargetSpan: &Span{Begin: 6, End: 9, Index: 1}},
			},
			false,
		},
		{
			"error with incompatible combinator",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          2,
				overlap:    1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 16, Index: 3}, Value: "string"},
					{Span: &Span{Begin: 18, End: 21, Index: 4}, Value: "test"},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNGramSeeder(tt.fields.n, tt.fields.overlap, tt.fields.similarity, tt.fields.combinator)
			gotSeeds, err := s.Seed(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Sort seeds by index, since seed computation is not deterministic
			sort.Slice(gotSeeds, func(i, j int) bool {
				if gotSeeds[i].SourceSpan.Index < gotSeeds[j].SourceSpan.Index {
					return true
				} else {
					return false
				}
			})
			if !reflect.DeepEqual(gotSeeds, tt.wantSeeds) {
				t.Errorf("Seed() = %v, want %v", gotSeeds, tt.wantSeeds)
			}
		})
	}
}

func TestNGramSeeder_computeSeedCandidates(t *testing.T) {
	type fields struct {
		similarity similarity.Similarity
		combinator combinating.ValueCombinator
		n          uint32
		overlap    uint32
	}
	type args struct {
		seq []*Element
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantSeedCandidates []*Element
		wantErr            bool
	}{
		{
			"string values simple",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          2,
				overlap:    1,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: "This is"},
				{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: "is a"},
				{Span: &Span{Begin: 9, End: 14, Index: 2}, Value: "a test"},
				{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: "test string"},
			},
			false,
		},
		{
			"string values error",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          2,
				overlap:    1,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
				},
			},
			nil,
			true,
		},
		{
			"string values complex",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          5,
				overlap:    4,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 3, Index: 0}, Value: "Die"},
					{Span: &Span{Begin: 5, End: 10, Index: 1}, Value: "Elster"},
					{Span: &Span{Begin: 12, End: 16, Index: 2}, Value: "(Pica"},
					{Span: &Span{Begin: 18, End: 22, Index: 3}, Value: "pica)"},
					{Span: &Span{Begin: 24, End: 26, Index: 4}, Value: "ist"},
					{Span: &Span{Begin: 28, End: 31, Index: 5}, Value: "eine"},
					{Span: &Span{Begin: 33, End: 40, Index: 6}, Value: "Vogelart"},
					{Span: &Span{Begin: 42, End: 44, Index: 7}, Value: "aus"},
					{Span: &Span{Begin: 46, End: 48, Index: 8}, Value: "der"},
					{Span: &Span{Begin: 50, End: 56, Index: 9}, Value: "Familie"},
					{Span: &Span{Begin: 58, End: 60, Index: 10}, Value: "der"},
					{Span: &Span{Begin: 61, End: 70, Index: 11}, Value: "Rabenvögel"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 26, Index: 0}, Value: "Die Elster (Pica pica) ist"},
				{Span: &Span{Begin: 5, End: 31, Index: 1}, Value: "Elster (Pica pica) ist eine"},
				{Span: &Span{Begin: 12, End: 40, Index: 2}, Value: "(Pica pica) ist eine Vogelart"},
				{Span: &Span{Begin: 18, End: 44, Index: 3}, Value: "pica) ist eine Vogelart aus"},
				{Span: &Span{Begin: 24, End: 48, Index: 4}, Value: "ist eine Vogelart aus der"},
				{Span: &Span{Begin: 28, End: 56, Index: 5}, Value: "eine Vogelart aus der Familie"},
				{Span: &Span{Begin: 33, End: 60, Index: 6}, Value: "Vogelart aus der Familie der"},
				{Span: &Span{Begin: 42, End: 70, Index: 7}, Value: "aus der Familie der Rabenvögel"},
			},
			false,
		},
		{
			"string values complex 2",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          5,
				overlap:    4,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 3, Index: 0}, Value: "Die"},
					{Span: &Span{Begin: 5, End: 16, Index: 1}, Value: "Hudsonelster"},
					{Span: &Span{Begin: 18, End: 22, Index: 2}, Value: "(Pica"},
					{Span: &Span{Begin: 24, End: 32, Index: 3}, Value: "hudsonia)"},
					{Span: &Span{Begin: 34, End: 36, Index: 4}, Value: "ist"},
					{Span: &Span{Begin: 38, End: 41, Index: 5}, Value: "eine"},
					{Span: &Span{Begin: 45, End: 54, Index: 6}, Value: "Singvogelart"},
					{Span: &Span{Begin: 56, End: 58, Index: 7}, Value: "aus"},
					{Span: &Span{Begin: 60, End: 62, Index: 8}, Value: "der"},
					{Span: &Span{Begin: 64, End: 70, Index: 9}, Value: "Familie"},
					{Span: &Span{Begin: 72, End: 74, Index: 10}, Value: "der"},
					{Span: &Span{Begin: 76, End: 85, Index: 11}, Value: "Rabenvögel"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 36, Index: 0}, Value: "Die Hudsonelster (Pica hudsonia) ist"},
				{Span: &Span{Begin: 5, End: 41, Index: 1}, Value: "Hudsonelster (Pica hudsonia) ist eine"},
				{Span: &Span{Begin: 18, End: 54, Index: 2}, Value: "(Pica hudsonia) ist eine Singvogelart"},
				{Span: &Span{Begin: 24, End: 58, Index: 3}, Value: "hudsonia) ist eine Singvogelart aus"},
				{Span: &Span{Begin: 34, End: 62, Index: 4}, Value: "ist eine Singvogelart aus der"},
				{Span: &Span{Begin: 38, End: 70, Index: 5}, Value: "eine Singvogelart aus der Familie"},
				{Span: &Span{Begin: 45, End: 74, Index: 6}, Value: "Singvogelart aus der Familie der"},
				{Span: &Span{Begin: 56, End: 85, Index: 7}, Value: "aus der Familie der Rabenvögel"},
			},
			false,
		},
		{
			"vector values simple",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          2,
				overlap:    1,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: []float32{3}},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: []float32{4}},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: []float32{5}},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: mat.NewVecDense(2, []float64{1, 2})},
				{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: mat.NewVecDense(2, []float64{2, 3})},
				{Span: &Span{Begin: 9, End: 14, Index: 2}, Value: mat.NewVecDense(2, []float64{3, 4})},
				{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: mat.NewVecDense(2, []float64{4, 5})},
			},
			false,
		},
		{
			"vector values error",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(""),
				n:          2,
				overlap:    1,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: []float32{3}},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: []float32{4}},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: []float32{5}},
				},
			},
			nil,
			true,
		},
		{
			"no overlap vector",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          1,
				overlap:    0,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: []float32{3}},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: []float32{4}},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: []float32{5}},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: mat.NewVecDense(1, []float64{1})},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: mat.NewVecDense(1, []float64{2})},
				{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: mat.NewVecDense(1, []float64{3})},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: mat.NewVecDense(1, []float64{4})},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: mat.NewVecDense(1, []float64{5})},
			},
			false,
		},
		{
			"no overlap string",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          1,
				overlap:    0,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
				{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
			},
			false,
		},
		{
			"type error",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          1,
				overlap:    0,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float64{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float64{2}},
				},
			},
			nil,
			true,
		},
		{
			"incompatible combinator error",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          1,
				overlap:    0,
			},
			args{
				seq: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNGramSeeder(tt.fields.n, tt.fields.overlap, tt.fields.similarity, tt.fields.combinator)
			gotSeedCandidates, err := s.computeSeedCandidates(tt.args.seq)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeSeedCandidates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSeedCandidates, tt.wantSeedCandidates) {
				t.Errorf("computeSeedCandidates() gotSeedCandidates = %v, want %v", gotSeedCandidates, tt.wantSeedCandidates)
			}
		})
	}
}

func TestNGramSeeder_computeSeeds(t *testing.T) {
	type fields struct {
		similarity similarity.Similarity
		combinator combinating.ValueCombinator
		n          uint32
		overlap    uint32
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Seed
		wantErr bool
	}{
		{
			"string value simple",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          2,
				overlap:    1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: "This is"},
					{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: "is a"},
					{Span: &Span{Begin: 9, End: 15, Index: 2}, Value: "a test"},
					{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: "test string"},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: "This is"},
					{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: "is a"},
					{Span: &Span{Begin: 9, End: 16, Index: 2}, Value: "a string"},
					{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: "string test"},
				},
			},
			[]*Seed{
				{SourceSpan: &Span{Begin: 1, End: 7, Index: 0}, TargetSpan: &Span{Begin: 1, End: 7, Index: 0}},
				{SourceSpan: &Span{Begin: 6, End: 9, Index: 1}, TargetSpan: &Span{Begin: 6, End: 9, Index: 1}},
			},
			false,
		},
		{
			"string value complex",
			fields{
				similarity: similarity.NewLevenshteinStringSimilarity(5),
				combinator: combinating.NewStringValueCombinator(" "),
				n:          5,
				overlap:    4,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 26, Index: 0}, Value: "Die Elster (Pica pica) ist"},
					{Span: &Span{Begin: 5, End: 31, Index: 1}, Value: "Elster (Pica pica) ist eine"},
					{Span: &Span{Begin: 12, End: 40, Index: 2}, Value: "(Pica pica) ist eine Vogelart"},
					{Span: &Span{Begin: 18, End: 44, Index: 3}, Value: "pica) ist eine Vogelart aus"},
					{Span: &Span{Begin: 24, End: 48, Index: 4}, Value: "ist eine Vogelart aus der"},
					{Span: &Span{Begin: 28, End: 56, Index: 5}, Value: "eine Vogelart aus der Familie"},
					{Span: &Span{Begin: 33, End: 60, Index: 6}, Value: "Vogelart aus der Familie der"},
					{Span: &Span{Begin: 42, End: 70, Index: 7}, Value: "aus der Familie der Rabenvögel"},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 36, Index: 0}, Value: "Die Hudsonelster (Pica hudsonia) ist"},
					{Span: &Span{Begin: 5, End: 41, Index: 1}, Value: "Hudsonelster (Pica hudsonia) ist eine"},
					{Span: &Span{Begin: 18, End: 54, Index: 2}, Value: "(Pica hudsonia) ist eine Singvogelart"},
					{Span: &Span{Begin: 24, End: 58, Index: 3}, Value: "hudsonia) ist eine Singvogelart aus"},
					{Span: &Span{Begin: 34, End: 62, Index: 4}, Value: "ist eine Singvogelart aus der"},
					{Span: &Span{Begin: 38, End: 70, Index: 5}, Value: "eine Singvogelart aus der Familie"},
					{Span: &Span{Begin: 45, End: 74, Index: 6}, Value: "Singvogelart aus der Familie der"},
					{Span: &Span{Begin: 56, End: 85, Index: 7}, Value: "aus der Familie der Rabenvögel"},
				},
			},
			[]*Seed{
				{SourceSpan: &Span{Begin: 24, End: 48, Index: 4}, TargetSpan: &Span{Begin: 34, End: 62, Index: 4}},
				{SourceSpan: &Span{Begin: 28, End: 56, Index: 5}, TargetSpan: &Span{Begin: 38, End: 70, Index: 5}},
				{SourceSpan: &Span{Begin: 33, End: 60, Index: 6}, TargetSpan: &Span{Begin: 45, End: 74, Index: 6}},
				{SourceSpan: &Span{Begin: 42, End: 70, Index: 7}, TargetSpan: &Span{Begin: 56, End: 85, Index: 7}},
			},
			false,
		},
		{
			"vector value simple",
			fields{
				similarity: similarity.NewIdentitySimilarity(),
				combinator: combinating.NewVectorValueCombinator(),
				n:          2,
				overlap:    1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: []float32{0, 1}},
					{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: []float32{1, 2}},
					{Span: &Span{Begin: 9, End: 15, Index: 2}, Value: []float32{2, 3}},
					{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: []float32{3, 4}},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 7, Index: 0}, Value: []float32{0, 1}},
					{Span: &Span{Begin: 6, End: 9, Index: 1}, Value: []float32{1, 2}},
					{Span: &Span{Begin: 9, End: 16, Index: 2}, Value: []float32{2, 4}},
					{Span: &Span{Begin: 11, End: 21, Index: 3}, Value: []float32{4, 3}},
				},
			},
			[]*Seed{
				{SourceSpan: &Span{Begin: 1, End: 7, Index: 0}, TargetSpan: &Span{Begin: 1, End: 7, Index: 0}},
				{SourceSpan: &Span{Begin: 6, End: 9, Index: 1}, TargetSpan: &Span{Begin: 6, End: 9, Index: 1}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewNGramSeeder(tt.fields.n, tt.fields.overlap, tt.fields.similarity, tt.fields.combinator)
			got, err := s.computeSeeds(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Since order is not guaranteed for the parallelization, restore order here for comparison purposes
			sort.Slice(got, func(i, j int) bool {
				return got[i].SourceSpan.Index < got[j].SourceSpan.Index
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeSeedsWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmarks for the computeSeeds function
func benchmarkNGramSeedercomputeSeeds(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewNGramSeeder(5, 4, similarity.NewCosineVectorSimilarity(0.9), combinating.NewVectorValueCombinator())
	// run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 100)
		target := utility.RandomVecSeq(length, 100)
		_, _ = s.computeSeeds(source, target)
	}
}
func BenchmarkNGramSeeder_computeSeeds10(b *testing.B)   { benchmarkNGramSeedercomputeSeeds(10, b) }
func BenchmarkNGramSeeder_computeSeeds100(b *testing.B)  { benchmarkNGramSeedercomputeSeeds(100, b) }
func BenchmarkNGramSeeder_computeSeeds1000(b *testing.B) { benchmarkNGramSeedercomputeSeeds(1000, b) }

// Benchmarks for the Seed function
func benchmarkNGramSeederSeed(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewNGramSeeder(5, 4, similarity.NewCosineVectorSimilarity(0.9), combinating.NewVectorValueCombinator())
	// run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 100)
		target := utility.RandomVecSeq(length, 100)
		_, _ = s.Seed(source, target)
	}
}
func BenchmarkNGramSeeder_Seed10(b *testing.B)   { benchmarkNGramSeederSeed(10, b) }
func BenchmarkNGramSeeder_Seed100(b *testing.B)  { benchmarkNGramSeederSeed(100, b) }
func BenchmarkNGramSeeder_Seed1000(b *testing.B) { benchmarkNGramSeederSeed(1000, b) }
