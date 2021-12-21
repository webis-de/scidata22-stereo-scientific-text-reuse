/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
	"picapica.org/picapica4/alignment/core/utility"
	"reflect"
	"testing"
)

func TestNeedlemanWunschSeeder_Seed(t *testing.T) {
	type fields struct {
		BaseSeeder      BaseSeeder
		mismatchPenalty float64
		gapPenalty      float64
		matchAward      float64
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Seed
	}{
		{
			name: "base test",
			fields: fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args: args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "C"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "U"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				},
			},
			want: []*Seed{
				{&Span{Begin: 0, End: 1, Index: 0}, &Span{Begin: 0, End: 1, Index: 0}},
				{&Span{Begin: 2, End: 3, Index: 2}, &Span{Begin: 1, End: 2, Index: 1}},
				{&Span{Begin: 3, End: 4, Index: 3}, &Span{Begin: 3, End: 4, Index: 3}},
				{&Span{Begin: 4, End: 5, Index: 4}, &Span{Begin: 4, End: 5, Index: 4}},
				{&Span{Begin: 5, End: 6, Index: 5}, &Span{Begin: 5, End: 6, Index: 5}},
				{&Span{Begin: 6, End: 7, Index: 6}, &Span{Begin: 6, End: 7, Index: 6}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NeedlemanWunschSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchPenalty,
				GapCost:      tt.fields.gapPenalty,
				MatchAward:   tt.fields.matchAward,
			}
			if got, _ := s.Seed(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Seed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeedlemanWunschSeeder_align(t *testing.T) {
	type fields struct {
		BaseSeeder      BaseSeeder
		mismatchPenalty float64
		gapPenalty      float64
		matchAward      float64
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSource []*Element
		wantTarget []*Element
	}{
		{
			"GCATGCU/GATTACA",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "C"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "U"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "U"},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: nil},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "C"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
			},
			[]*Element{
				{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: nil},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
			},
		},
		{
			"CGAGAGA/GAGAGA",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "C"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "G"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "A"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "C"},
			},
			[]*Element{
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "A"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "G"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: nil},
			},
		},
		{
			"CGAGAGA/GAGAGA",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
			},
			[]*Element{
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: nil},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "G"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "G"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
				{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NeedlemanWunschSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchPenalty,
				GapCost:      tt.fields.gapPenalty,
				MatchAward:   tt.fields.matchAward,
			}
			gotSource, gotTarget := s.align(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(gotSource, tt.wantSource) {
				t.Errorf("align() gotSource = %v, want %v", gotSource, tt.wantSource)
			}
			if !reflect.DeepEqual(gotTarget, tt.wantTarget) {
				t.Errorf("align() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
			}
		})
	}
}

func TestNeedlemanWunschSeeder_scoringTable(t *testing.T) {
	type fields struct {
		BaseSeeder      BaseSeeder
		mismatchPenalty float64
		gapPenalty      float64
		matchAward      float64
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantScore     *mat.Dense
		wantTraceback *[][][2]uint16
	}{
		{
			"base test",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "C"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "A"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "U"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "A"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "A"},
				},
			},
			mat.NewDense(8, 8, []float64{
				0, -1, -2, -3, -4, -5, -6, -7,
				-1, 1, 0, -1, -2, -3, -4, -5,
				-2, 0, 0, -1, -2, -3, -2, -3,
				-3, -1, 1, 0, -1, -1, -2, -1,
				-4, -2, 0, 2, 1, 0, -1, -2,
				-5, -3, -1, 1, 1, 0, -1, -2,
				-6, -4, -2, 0, 0, 0, 1, 0,
				-7, -5, -3, -1, -1, -1, 0, 0,
			}),
			&[][][2]uint16{
				{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}},
				{{0, 0}, {0, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}},
				{{0, 0}, {1, 1}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {1, 5}, {2, 6}},
				{{0, 0}, {2, 1}, {2, 1}, {3, 2}, {3, 3}, {2, 4}, {3, 5}, {2, 6}},
				{{0, 0}, {3, 1}, {3, 2}, {3, 2}, {3, 3}, {4, 4}, {4, 5}, {4, 6}},
				{{0, 0}, {4, 0}, {4, 2}, {4, 3}, {4, 3}, {4, 4}, {4, 5}, {4, 6}},
				{{0, 0}, {5, 1}, {5, 2}, {5, 3}, {5, 3}, {5, 4}, {5, 5}, {6, 6}},
				{{0, 0}, {6, 1}, {6, 2}, {6, 3}, {6, 3}, {6, 4}, {6, 6}, {6, 6}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NeedlemanWunschSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchPenalty,
				GapCost:      tt.fields.gapPenalty,
				MatchAward:   tt.fields.matchAward,
			}
			gotScore, gotTraceback := s.scoringTable(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(gotScore, tt.wantScore) {
				t.Errorf("populateScoringTable() Scores = %v, want %v", gotScore, tt.wantScore)
			}
			if !reflect.DeepEqual(gotTraceback, tt.wantTraceback) {
				t.Errorf("populateScoringTable() Traceback = %v, want %v", gotTraceback, tt.wantTraceback)
			}
		})
	}
}

func TestNeedlemanWunschSeeder_score(t *testing.T) {
	type fields struct {
		BaseSeeder      BaseSeeder
		mismatchPenalty float64
		gapPenalty      float64
		matchAward      float64
	}
	type args struct {
		source *Element
		target *Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"unequal values",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: &Element{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
				target: &Element{Span: &Span{Begin: 1, End: 2, Index: 0}, Value: "B"},
			},
			-1,
		},
		{
			"equal values",
			fields{
				BaseSeeder:      NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchPenalty: -1,
				gapPenalty:      -1,
				matchAward:      1,
			},
			args{
				source: &Element{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
				target: &Element{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "A"},
			},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NeedlemanWunschSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchPenalty,
				GapCost:      tt.fields.gapPenalty,
				MatchAward:   tt.fields.matchAward,
			}
			if got := s.score(tt.args.source, tt.args.target); got != tt.want {
				t.Errorf("score() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmarks

// Benchmark for the Seed function
func benchmarkNeedlemanWunschSeederSeed(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewNeedlemanWunschSeeder(similarity.NewCosineVectorSimilarity(0.9), -1, -1, 1) // run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 300)
		target := utility.RandomVecSeq(length, 300)
		_, _ = s.Seed(source, target)
	}
}

// Benchmarks for different sequence lengths
func BenchmarkNeedlemanWunschSeeder_Seed10(b *testing.B) { benchmarkNeedlemanWunschSeederSeed(10, b) }
func BenchmarkNeedlemanWunschSeeder_Seed100(b *testing.B) {
	benchmarkNeedlemanWunschSeederSeed(100, b)
}
func BenchmarkNeedlemanWunschSeeder_Seed1000(b *testing.B) {
	benchmarkNeedlemanWunschSeederSeed(1000, b)
}
