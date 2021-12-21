/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	"math"
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func TestIndexCluster_FromSeeds(t *testing.T) {
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name string
		args args
		want []Clusterable
	}{
		{
			"test with one seed",
			args{seeds: []*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 3, Index: 0},
					TargetSpan: &Span{Begin: 5, End: 9, Index: 1},
				},
			}},
			[]Clusterable{
				IndexClusterable{x: 0, y: 1, id: 0},
			},
		},
		{
			"test with multiple seeds",
			args{seeds: []*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 3, Index: 0},
					TargetSpan: &Span{Begin: 5, End: 9, Index: 1},
				},
				{
					SourceSpan: &Span{Begin: 5, End: 13, Index: 4},
					TargetSpan: &Span{Begin: 29, End: 22, Index: 8},
				},
				{
					SourceSpan: &Span{Begin: 144, End: 149, Index: 16},
					TargetSpan: &Span{Begin: 5, End: 9, Index: 2},
				},
			}},
			[]Clusterable{
				IndexClusterable{x: 0, y: 1, id: 0},
				IndexClusterable{x: 4, y: 8, id: 0},
				IndexClusterable{x: 16, y: 2, id: 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewIndexCluster()
			c.FromSeeds(tt.args.seeds)
			// Check manually for occurrence of begin/ends of seeds, since hashes are only known at runtime
			for _, clusterable := range c.Clusterables() {
				found := false
				for _, check := range tt.want {
					if clusterable.(IndexClusterable).x == check.(IndexClusterable).x && clusterable.(IndexClusterable).y == check.(IndexClusterable).y {
						found = true
					}
				}
				if !found {
					t.Errorf("FromSeeds(), %v not found in %v", clusterable, tt.want)
				}
			}
		})
	}
}

func TestIndexCluster_Lookup(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		mapping      map[uint64]*Seed
	}
	type args struct {
		clusterables Clusterable
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
				clusterables: []Clusterable{IndexClusterable{0, 0, 0}, IndexClusterable{1, 1, 1}},
				mapping: map[uint64]*Seed{
					0: {
						SourceSpan: &Span{Begin: 0, End: 0, Index: 0},
						TargetSpan: &Span{Begin: 1, End: 1, Index: 0},
					},
					1: {
						SourceSpan: &Span{Begin: 0, End: 0, Index: 1},
						TargetSpan: &Span{Begin: 1, End: 1, Index: 1},
					},
				},
			},
			args{clusterables: &IndexClusterable{1, 1, 1}},
			[]*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 0, Index: 1},
					TargetSpan: &Span{Begin: 1, End: 1, Index: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &IndexCluster{
				clusterables: tt.fields.clusterables,
				mapping:      tt.fields.mapping,
			}
			if gotSeeds := c.Lookup(tt.args.clusterables); !reflect.DeepEqual(gotSeeds, tt.wantSeeds) {
				t.Errorf("Lookup() = %v, want %v", gotSeeds, tt.wantSeeds)
			}
		})
	}
}

func TestIndexCluster_ToSeeds(t *testing.T) {
	type fields struct {
		clusterables []Clusterable
		mapping      map[uint64]*Seed
	}
	tests := []struct {
		name      string
		fields    fields
		wantSeeds []*Seed
	}{
		{
			"base test",
			fields{
				clusterables: []Clusterable{IndexClusterable{0, 0, 0}, IndexClusterable{1, 1, 1}},
				mapping: map[uint64]*Seed{
					0: {
						SourceSpan: &Span{Begin: 0, End: 0, Index: 0},
						TargetSpan: &Span{Begin: 1, End: 1, Index: 0},
					},
					1: {
						SourceSpan: &Span{Begin: 0, End: 0, Index: 1},
						TargetSpan: &Span{Begin: 1, End: 1, Index: 1},
					},
				},
			},
			[]*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 0, Index: 0},
					TargetSpan: &Span{Begin: 1, End: 1, Index: 0},
				},
				{
					SourceSpan: &Span{Begin: 0, End: 0, Index: 1},
					TargetSpan: &Span{Begin: 1, End: 1, Index: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &IndexCluster{
				clusterables: tt.fields.clusterables,
				mapping:      tt.fields.mapping,
			}
			if gotSeeds := c.ToSeeds(); !reflect.DeepEqual(gotSeeds, tt.wantSeeds) {
				t.Errorf("ToSeeds() = %v, want %v", gotSeeds, tt.wantSeeds)
			}
		})
	}
}

func TestIndexClusterable_Distance(t *testing.T) {
	type fields struct {
		x  uint32
		y  uint32
		id uint64
	}
	type args struct {
		other Clusterable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			"identical coordinates",
			fields{
				x:  1,
				y:  1,
				id: 0,
			},
			args{other: IndexClusterable{x: 1, y: 1, id: 1}},
			0,
			false,
		},
		{
			"arbitrary coordinates",
			fields{
				x:  5,
				y:  6,
				id: 0,
			},
			args{other: IndexClusterable{x: 1, y: 3, id: 1}},
			5,
			false,
		},
		{
			"symmetrical coordinates",
			fields{
				x:  0,
				y:  1,
				id: 0,
			},
			args{other: IndexClusterable{x: 1, y: 0, id: 1}},
			math.Sqrt(2),
			false,
		},
		{
			"shared x coordinate",
			fields{
				x:  0,
				y:  0,
				id: 0,
			},
			args{other: IndexClusterable{x: 0, y: 4, id: 1}},
			4,
			false,
		},
		{
			"shared y coordinate",
			fields{
				x:  0,
				y:  0,
				id: 0,
			},
			args{other: IndexClusterable{x: 4, y: 0, id: 1}},
			4,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &IndexClusterable{
				x:  tt.fields.x,
				y:  tt.fields.y,
				id: tt.fields.id,
			}
			got, err := c.Distance(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Distance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Distance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexClusterable_ID(t *testing.T) {
	type fields struct {
		x  uint32
		y  uint32
		id uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			"base test",
			fields{
				x:  0,
				y:  1,
				id: 1954278510971257,
			},
			1954278510971257,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &IndexClusterable{
				x:  tt.fields.x,
				y:  tt.fields.y,
				id: tt.fields.id,
			}
			if got := c.ID(); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIndexCluster(t *testing.T) {
	tests := []struct {
		name string
		want *IndexCluster
	}{
		{
			"generator test",
			&IndexCluster{clusterables: make([]Clusterable, 0, 0), mapping: make(map[uint64]*Seed)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIndexCluster(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIndexCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
