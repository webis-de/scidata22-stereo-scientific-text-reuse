/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"gonum.org/v1/gonum/mat"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/normalizing"
	. "picapica.org/picapica4/alignment/core/types"
	"picapica.org/picapica4/alignment/core/utility"
	"reflect"
	"sort"
	"testing"
)

func TestHashSeeder_Seed(t *testing.T) {
	type fields struct {
		BaseSeeder BaseSeeder
		Normalizer normalizing.Normalizer
		Combinator combinating.ValueCombinator
		n          int
		overlap    int
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
			"string values simple",
			fields{
				Normalizer: normalizing.NewIdentityNormalizer(),
				Combinator: combinating.NewStringValueCombinator(""),
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
			"string values complex",
			fields{
				Normalizer: normalizing.NewIdentityNormalizer(),
				Combinator: combinating.NewStringValueCombinator(" "),
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
			[]*Seed{},
			false,
		},
		{
			"vector values",
			fields{
				Normalizer: normalizing.NewIdentityNormalizer(),
				Combinator: combinating.NewVectorValueCombinator(),
				n:          2,
				overlap:    1,
			},
			args{
				source: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: []float32{3}},
					{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: []float32{4}},
					{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: []float32{5}},
				},
				target: []*Element{
					{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: []float32{1}},
					{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: []float32{2}},
					{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: []float32{3}},
					{Span: &Span{Begin: 11, End: 16, Index: 3}, Value: []float32{5}},
					{Span: &Span{Begin: 18, End: 21, Index: 4}, Value: []float32{4}},
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
			s := &HashSeeder{
				BaseSeeder: tt.fields.BaseSeeder,
				Normalizer: tt.fields.Normalizer,
				Combinator: tt.fields.Combinator,
				N:          tt.fields.n,
				Overlap:    tt.fields.overlap,
			}
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
			if !reflect.DeepEqual(gotSeeds, tt.want) {
				t.Errorf("Seed() = %v, want %v", gotSeeds, tt.want)
			}
		})
	}
}

func TestHashSeeder_computeSeedCandidates(t *testing.T) {
	type fields struct {
		combinator combinating.ValueCombinator
		normalizer normalizing.Normalizer
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
				combinator: combinating.NewStringValueCombinator(" "),
				normalizer: normalizing.NewIdentityNormalizer(),
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
			"string values complex",
			fields{
				combinator: combinating.NewStringValueCombinator(" "),
				normalizer: normalizing.NewIdentityNormalizer(),
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
				combinator: combinating.NewStringValueCombinator(" "),
				normalizer: normalizing.NewIdentityNormalizer(),
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
				combinator: combinating.NewVectorValueCombinator(),
				normalizer: normalizing.NewIdentityNormalizer(),
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
			"no overlap vector",
			fields{
				combinator: combinating.NewVectorValueCombinator(),
				normalizer: normalizing.NewIdentityNormalizer(),
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
				combinator: combinating.NewStringValueCombinator(" "),
				normalizer: normalizing.NewIdentityNormalizer(),
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
				combinator: combinating.NewVectorValueCombinator(),
				normalizer: normalizing.NewIdentityNormalizer(),
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
				combinator: combinating.NewStringValueCombinator(" "),
				normalizer: normalizing.NewIdentityNormalizer(),
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
			s := NewHashSeeder(tt.fields.n, tt.fields.overlap, tt.fields.normalizer, tt.fields.combinator)
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

func TestHashSeeder_computeSeeds(t *testing.T) {
	type fields struct {
		BaseSeeder BaseSeeder
		Combinator combinating.ValueCombinator
		n          int
		overlap    int
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
			"base test",
			fields{
				Combinator: combinating.NewStringValueCombinator(" "),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HashSeeder{
				BaseSeeder: tt.fields.BaseSeeder,
				Combinator: tt.fields.Combinator,
				N:          tt.fields.n,
				Overlap:    tt.fields.overlap,
			}
			gotSeeds, err := s.computeSeeds(tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeSeeds() error = %v, wantErr %v", err, tt.wantErr)
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
			if !reflect.DeepEqual(gotSeeds, tt.want) {
				t.Errorf("computeSeed() = %v, want %v", gotSeeds, tt.want)
			}
		})
	}
}

// Benchmarks for the computeSeeds function
func benchmarkHashSeedercomputeSeeds(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewHashSeeder(5, 4, normalizing.NewIdentityNormalizer(), combinating.NewVectorValueCombinator())
	// run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 100)
		target := utility.RandomVecSeq(length, 100)
		_, _ = s.computeSeeds(source, target)
	}
}
func BenchmarkHashSeedercomputeSeeds10(b *testing.B)   { benchmarkHashSeedercomputeSeeds(10, b) }
func BenchmarkHashSeedercomputeSeeds100(b *testing.B)  { benchmarkHashSeedercomputeSeeds(100, b) }
func BenchmarkHashSeedercomputeSeeds1000(b *testing.B) { benchmarkHashSeedercomputeSeeds(1000, b) }

// Benchmarks for the Seed function
func benchmarkHashSeederSeed(length int, b *testing.B) {
	b.ReportAllocs()
	s := NewHashSeeder(5, 4, normalizing.NewIdentityNormalizer(), combinating.NewVectorValueCombinator())
	// run the seed function b.N times, generating new input data every time
	for n := 0; n < b.N; n++ {
		source := utility.RandomVecSeq(length, 100)
		target := utility.RandomVecSeq(length, 100)
		_, _ = s.Seed(source, target)
	}
}
func BenchmarkHashSeederSeed10(b *testing.B)   { benchmarkHashSeederSeed(10, b) }
func BenchmarkHashSeederSeed100(b *testing.B)  { benchmarkHashSeederSeed(100, b) }
func BenchmarkHashSeederSeed1000(b *testing.B) { benchmarkHashSeederSeed(1000, b) }
