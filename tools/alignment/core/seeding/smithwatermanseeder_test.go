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

func TestSmithWatermanSeeder_Seed(t *testing.T) {
	type fields struct {
		BaseSeeder   BaseSeeder
		mismatchCost float64
		gapCost      float64
		matchAward   float64
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
	}{
		{
			"base test",
			fields{
				BaseSeeder:   NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchCost: -3,
				gapCost:      -2,
				matchAward:   3,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "T"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "G"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "G"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "T"},
					{Span: &Span{Begin: 8, End: 9, Index: 8}, Value: "A"},
				},
			},
			[]*Seed{{
				SourceSpan: &Span{
					Begin: 1,
					End:   6,
					Index: 1,
				},
				TargetSpan: &Span{
					Begin: 1,
					End:   7,
					Index: 1,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SmithWatermanSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchCost,
				GapCost:      tt.fields.gapCost,
				MatchAward:   tt.fields.matchAward,
			}
			if gotSeeds, _ := s.Seed(tt.args.source, tt.args.target); !reflect.DeepEqual(gotSeeds, tt.wantSeeds) {
				t.Errorf("Seed() = %v, want %v", gotSeeds, tt.wantSeeds)
			}
		})
	}
}

func TestSmithWatermanSeeder_score(t *testing.T) {
	type fields struct {
		BaseSeeder   BaseSeeder
		mismatchCost float64
		gapCost      float64
		matchAward   float64
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
				BaseSeeder:   NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchCost: -1,
				gapCost:      -1,
				matchAward:   1,
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
				BaseSeeder:   NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchCost: -1,
				gapCost:      -1,
				matchAward:   1,
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
			s := &SmithWatermanSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchCost,
				GapCost:      tt.fields.gapCost,
				MatchAward:   tt.fields.matchAward,
			}
			if got := s.score(tt.args.source, tt.args.target); got != tt.want {
				t.Errorf("score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSmithWatermanSeeder_scoringTable(t *testing.T) {
	type fields struct {
		BaseSeeder   BaseSeeder
		mismatchCost float64
		gapCost      float64
		matchAward   float64
	}
	type args struct {
		source []*Element
		target []*Element
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantScores   *mat.Dense
		wantMaxCoord [2]int
	}{
		{
			"base test",
			fields{
				BaseSeeder:   NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchCost: -3,
				gapCost:      -2,
				matchAward:   3,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "T"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "G"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "G"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "T"},
					{Span: &Span{Begin: 8, End: 9, Index: 8}, Value: "A"},
				},
			},
			mat.NewDense(9, 10, []float64{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 3, 3, 1, 0, 0, 3, 1,
				0, 3, 3, 1, 1, 6, 4, 2, 1, 0,
				0, 1, 1, 6, 4, 4, 3, 1, 5, 3,
				0, 0, 0, 4, 9, 7, 5, 3, 4, 2,
				0, 0, 0, 2, 7, 6, 10, 8, 6, 7,
				0, 0, 0, 0, 5, 4, 8, 13, 11, 9,
				0, 3, 3, 1, 3, 8, 6, 11, 10, 8,
				0, 3, 6, 4, 2, 6, 5, 9, 8, 7,
			}),
			[2]int{6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SmithWatermanSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchCost,
				GapCost:      tt.fields.gapCost,
				MatchAward:   tt.fields.matchAward,
			}
			gotScores, gotMaxCoord := s.scoringTable(tt.args.source, tt.args.target)
			if !reflect.DeepEqual(gotScores, tt.wantScores) {
				t.Errorf("scoringTable() got Scores = %v, want %v", gotScores, tt.wantScores)
			}
			if !reflect.DeepEqual(gotMaxCoord, tt.wantMaxCoord) {
				t.Errorf("scoringTable() got MaxCoord = %v, want %v", gotMaxCoord, tt.wantMaxCoord)
			}
		})
	}
}

func TestSmithWatermanSeeder_align(t *testing.T) {
	type fields struct {
		BaseSeeder   BaseSeeder
		mismatchCost float64
		gapCost      float64
		matchAward   float64
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
			"base test",
			fields{
				BaseSeeder:   NewBaseSeeder(similarity.NewIdentitySimilarity()),
				mismatchCost: -3,
				gapCost:      -2,
				matchAward:   3,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "T"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "G"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "G"},
				},
				target: []*Element{
					{Span: &Span{Begin: 0, End: 1, Index: 0}, Value: "G"},
					{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
					{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
					{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
					{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
					{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
					{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
					{Span: &Span{Begin: 7, End: 8, Index: 7}, Value: "T"},
					{Span: &Span{Begin: 8, End: 9, Index: 8}, Value: "A"},
				},
			},
			[]*Element{
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "C"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "A"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: nil},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
			},
			[]*Element{
				{Span: &Span{Begin: 6, End: 7, Index: 6}, Value: "C"},
				{Span: &Span{Begin: 5, End: 6, Index: 5}, Value: "A"},
				{Span: &Span{Begin: 4, End: 5, Index: 4}, Value: "G"},
				{Span: &Span{Begin: 3, End: 4, Index: 3}, Value: "T"},
				{Span: &Span{Begin: 2, End: 3, Index: 2}, Value: "T"},
				{Span: &Span{Begin: 1, End: 2, Index: 1}, Value: "G"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SmithWatermanSeeder{
				BaseSeeder:   tt.fields.BaseSeeder,
				MismatchCost: tt.fields.mismatchCost,
				GapCost:      tt.fields.gapCost,
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

// Benchmark for the Seed function
func benchmarkSmithWatermanSeederSeed(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewSmithWatermanSeeder(similarity.NewCosineVectorSimilarity(0.9), -3, -2, 3)
	// run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 300)
		target := utility.RandomVecSeq(length, 300)
		_, _ = s.Seed(source, target)
	}
}

// Benchmarks for different sequence lengths
func BenchmarkSmithWatermanSeeder_Seed10(b *testing.B)   { benchmarkSmithWatermanSeederSeed(10, b) }
func BenchmarkSmithWatermanSeeder_Seed100(b *testing.B)  { benchmarkSmithWatermanSeederSeed(100, b) }
func BenchmarkSmithWatermanSeeder_Seed1000(b *testing.B) { benchmarkSmithWatermanSeederSeed(1000, b) }
