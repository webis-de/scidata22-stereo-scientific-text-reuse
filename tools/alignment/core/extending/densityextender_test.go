/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"sort"
	"testing"
)

func TestDBScan_Cluster(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		minPts       int
		eps          float64
		size         int
		reachable    map[Clusterable]map[Clusterable]bool
	}
	type args struct {
		clusterables []Clusterable
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantClusters [][]Clusterable
	}{
		{
			"base",
			fields{
				minPts: 2,
				eps:    4,
			},
			args{clusterables: []Clusterable{
				ValueClusterable{value: 5, id: 0},
				ValueClusterable{value: 7, id: 1},
				ValueClusterable{value: 8, id: 2},
				ValueClusterable{value: 10, id: 3},
				ValueClusterable{value: 100, id: 4},
				ValueClusterable{value: 103, id: 5},
				ValueClusterable{value: 104, id: 6},
				ValueClusterable{value: 0, id: 7},
			}},
			[][]Clusterable{
				{
					ValueClusterable{value: 100, id: 4},
					ValueClusterable{value: 103, id: 5},
					ValueClusterable{value: 104, id: 6},
				}, {
					ValueClusterable{value: 5, id: 0},
					ValueClusterable{value: 7, id: 1},
					ValueClusterable{value: 8, id: 2},
					ValueClusterable{value: 10, id: 3},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DBScan{
				clusterables: tt.fields.clusterables,
				minPts:       tt.fields.minPts,
				eps:          tt.fields.eps,
				size:         tt.fields.size,
				reachable:    tt.fields.reachable,
			}
			gotClusters := e.Cluster(tt.args.clusterables)
			// Sort to make comparable, since order is not guaranteed by algorithm
			for _, slice := range gotClusters {
				sort.Slice(slice, func(i, j int) bool {
					return slice[i].ID() < slice[j].ID()
				})
			}
			sort.Slice(gotClusters, func(i, j int) bool {
				return len(gotClusters[i]) < len(gotClusters[j])
			})
			if !reflect.DeepEqual(gotClusters, tt.wantClusters) {
				t.Errorf("Cluster() = %v, want %v", gotClusters, tt.wantClusters)
			}
		})
	}
}

func TestDBScan_DistanceMatrix(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		minPts       int
		eps          float64
		size         int
	}
	tests := []struct {
		name          string
		fields        fields
		wantDistances map[Clusterable]map[Clusterable]bool
	}{
		{
			"value clusterables",
			fields{
				clusterables: []Clusterable{
					ValueClusterable{value: 5, id: 0},
					ValueClusterable{value: 7, id: 1},
					ValueClusterable{value: 8, id: 2},
					ValueClusterable{value: 10, id: 3},
				},
				minPts: 0,
				eps:    4,
				size:   4,
			},
			map[Clusterable]map[Clusterable]bool{
				ValueClusterable{value: 5, id: 0}: {
					ValueClusterable{value: 5, id: 0}:  true,
					ValueClusterable{value: 7, id: 1}:  true,
					ValueClusterable{value: 8, id: 2}:  true,
					ValueClusterable{value: 10, id: 3}: false,
				},
				ValueClusterable{value: 7, id: 1}: {
					ValueClusterable{value: 5, id: 0}:  true,
					ValueClusterable{value: 7, id: 1}:  true,
					ValueClusterable{value: 8, id: 2}:  true,
					ValueClusterable{value: 10, id: 3}: true,
				},
				ValueClusterable{value: 8, id: 2}: {
					ValueClusterable{value: 5, id: 0}:  true,
					ValueClusterable{value: 7, id: 1}:  true,
					ValueClusterable{value: 8, id: 2}:  true,
					ValueClusterable{value: 10, id: 3}: true,
				},
				ValueClusterable{value: 10, id: 3}: {
					ValueClusterable{value: 5, id: 0}:  false,
					ValueClusterable{value: 7, id: 1}:  true,
					ValueClusterable{value: 8, id: 2}:  true,
					ValueClusterable{value: 10, id: 3}: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DBScan{
				clusterables: tt.fields.clusterables,
				minPts:       tt.fields.minPts,
				eps:          tt.fields.eps,
				size:         tt.fields.size,
			}
			if gotDistances := e.DistanceMatrix(e.clusterables); !reflect.DeepEqual(gotDistances, tt.wantDistances) {
				t.Errorf("DistanceMatrix() = %v, want %v", gotDistances, tt.wantDistances)
			}
		})
	}
}

func TestDBScan_Neighborhood(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		minPts       int
		eps          float64
		size         int
	}
	type args struct {
		q Clusterable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Clusterable
	}{
		{
			"query 1",
			fields{
				clusterables: []Clusterable{
					ValueClusterable{value: 5, id: 0},
					ValueClusterable{value: 7, id: 1},
					ValueClusterable{value: 8, id: 2},
					ValueClusterable{value: 10, id: 3},
					ValueClusterable{value: 100, id: 4},
					ValueClusterable{value: 101, id: 5},
					ValueClusterable{value: 102, id: 6},
				},
				minPts: 0,
				eps:    4,
				size:   7,
			},
			args{
				ValueClusterable{value: 5, id: 0},
			},
			[]Clusterable{
				ValueClusterable{value: 5, id: 0},
				ValueClusterable{value: 7, id: 1},
				ValueClusterable{value: 8, id: 2},
			},
		},
		{
			"query 2",
			fields{
				clusterables: []Clusterable{
					ValueClusterable{value: 5, id: 0},
					ValueClusterable{value: 7, id: 1},
					ValueClusterable{value: 8, id: 2},
					ValueClusterable{value: 10, id: 3},
					ValueClusterable{value: 100, id: 4},
					ValueClusterable{value: 101, id: 5},
					ValueClusterable{value: 102, id: 6},
				},
				minPts: 0,
				eps:    4,
				size:   7,
			},
			args{
				ValueClusterable{value: 100, id: 4},
			},
			[]Clusterable{
				ValueClusterable{value: 100, id: 4},
				ValueClusterable{value: 101, id: 5},
				ValueClusterable{value: 102, id: 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DBScan{
				clusterables: tt.fields.clusterables,
				minPts:       tt.fields.minPts,
				eps:          tt.fields.eps,
				size:         tt.fields.size,
			}
			e.reachable = e.DistanceMatrix(e.clusterables)
			if got := e.Neighborhood(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Neighborhood() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBScan_computeClusters(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		minPts       int
		eps          float64
		size         int
		reachable    map[Clusterable]map[Clusterable]bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[Clusterable]string
	}{
		{
			"base",
			fields{
				clusterables: []Clusterable{
					ValueClusterable{value: 5, id: 0},
					ValueClusterable{value: 7, id: 1},
					ValueClusterable{value: 8, id: 2},
					ValueClusterable{value: 10, id: 3},
					ValueClusterable{value: 100, id: 4},
					ValueClusterable{value: 103, id: 5},
					ValueClusterable{value: 104, id: 6},
					ValueClusterable{value: 0, id: 7},
				},
				minPts: 2,
				eps:    4,
				size:   8,
			},
			map[Clusterable]string{
				ValueClusterable{value: 5, id: 0}:   "1",
				ValueClusterable{value: 7, id: 1}:   "1",
				ValueClusterable{value: 8, id: 2}:   "1",
				ValueClusterable{value: 10, id: 3}:  "1",
				ValueClusterable{value: 100, id: 4}: "2",
				ValueClusterable{value: 103, id: 5}: "2",
				ValueClusterable{value: 104, id: 6}: "2",
				ValueClusterable{value: 0, id: 7}:   "noise",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DBScan{
				clusterables: tt.fields.clusterables,
				minPts:       tt.fields.minPts,
				eps:          tt.fields.eps,
				size:         tt.fields.size,
			}
			e.reachable = e.DistanceMatrix(e.clusterables)
			if got := e.computeClusters(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDensityExtender_Extend(t *testing.T) {
	type fields struct {
		minPoints int
		epsilon   float64
	}
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantExtension []*Seed
		wantErr       bool
	}{
		{
			"base",
			fields{
				minPoints: 2,
				epsilon:   3,
			},
			args{seeds: []*Seed{
				{SourceSpan: &Span{Begin: 0, End: 4, Index: 0}, TargetSpan: &Span{Begin: 0, End: 4, Index: 0}},
				{SourceSpan: &Span{Begin: 5, End: 6, Index: 2}, TargetSpan: &Span{Begin: 17, End: 18, Index: 5}},
				{SourceSpan: &Span{Begin: 6, End: 9, Index: 3}, TargetSpan: &Span{Begin: 19, End: 20, Index: 6}},
				{SourceSpan: &Span{Begin: 12, End: 13, Index: 5}, TargetSpan: &Span{Begin: 28, End: 30, Index: 10}},
				{SourceSpan: &Span{Begin: 14, End: 15, Index: 6}, TargetSpan: &Span{Begin: 31, End: 32, Index: 11}},
				{SourceSpan: &Span{Begin: 21, End: 22, Index: 10}, TargetSpan: &Span{Begin: 54, End: 55, Index: 22}},
			}},
			[]*Seed{
				{SourceSpan: &Span{Begin: 5, End: 9, Index: 2}, TargetSpan: &Span{Begin: 17, End: 20, Index: 5}},
				{SourceSpan: &Span{Begin: 12, End: 15, Index: 5}, TargetSpan: &Span{Begin: 28, End: 32, Index: 10}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &DensityExtender{
				minPoints: tt.fields.minPoints,
				epsilon:   tt.fields.epsilon,
			}
			got, err := e.Extend(tt.args.seeds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extend() = %v, want %v", err, tt.wantErr)
			}
			// Sort since order is not guaranteed
			sort.Slice(got, func(i, j int) bool {
				return got[i].SourceSpan.Begin < got[j].SourceSpan.Begin
			})
			if !reflect.DeepEqual(got, tt.wantExtension) {
				t.Errorf("Extend() = %v, want %v", got, tt.wantExtension)
			}
		})
	}
}

func TestNewDBScan(t *testing.T) {
	type args struct {
		minPts int
		eps    float64
	}
	tests := []struct {
		name string
		args args
		want *DBScan
	}{
		{
			"constructor",
			args{minPts: 5, eps: 10},
			&DBScan{minPts: 5, eps: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDBScan(tt.args.minPts, tt.args.eps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBScan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDensityExtender(t *testing.T) {
	type args struct {
		minPoints int
		epsilon   float64
	}
	tests := []struct {
		name string
		args args
		want *DensityExtender
	}{
		{
			"constructor",
			args{minPoints: 5, epsilon: 10},
			&DensityExtender{minPoints: 5, epsilon: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDensityExtender(tt.args.minPoints, tt.args.epsilon); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDensityExtender() = %v, want %v", got, tt.want)
			}
		})
	}
}
