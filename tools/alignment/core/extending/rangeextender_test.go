/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func TestNewRangeExtender(t *testing.T) {
	type args struct {
		theta uint32
	}
	tests := []struct {
		name string
		args args
		want *RangeExtender
	}{
		{
			"constructor",
			args{theta: 100},
			&RangeExtender{theta: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRangeExtender(tt.args.theta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRangeExtender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeExtender_Extend(t *testing.T) {
	type fields struct {
		theta uint32
	}
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Seed
		wantErr bool
	}{
		{
			"base",
			fields{theta: 3},
			args{seeds: []*Seed{
				{&Span{0, 4, 0}, &Span{0, 4, 0}},
				{&Span{5, 6, 2}, &Span{17, 18, 5}},
				{&Span{6, 9, 3}, &Span{19, 20, 6}},
				{&Span{12, 13, 5}, &Span{28, 30, 10}},
				{&Span{14, 15, 6}, &Span{31, 32, 11}},
				{&Span{21, 22, 10}, &Span{54, 55, 22}},
			}},
			[]*Seed{
				{&Span{0, 4, 0}, &Span{0, 4, 0}},
				{&Span{5, 9, 2}, &Span{17, 20, 5}},
				{&Span{12, 15, 5}, &Span{28, 32, 10}},
				{&Span{21, 22, 10}, &Span{54, 55, 22}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RangeExtender{
				theta: tt.fields.theta,
			}
			got, err := e.Extend(tt.args.seeds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extend() = %v, want %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeExtender_candidatePassages(t *testing.T) {
	type fields struct {
		theta uint32
	}
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPassages [][]*Seed
	}{
		{
			"single elements at start and end",
			fields{theta: 3},
			args{seeds: []*Seed{
				{&Span{0, 0, 0}, &Span{0, 0, 0}},
				{&Span{0, 0, 5}, &Span{0, 0, 8}},
				{&Span{0, 0, 6}, &Span{0, 0, 5}},
				{&Span{0, 0, 9}, &Span{0, 0, 3}},
				{&Span{0, 0, 10}, &Span{0, 0, 11}},
				{&Span{0, 0, 22}, &Span{0, 0, 20}},
			}},
			[][]*Seed{
				{
					{&Span{0, 0, 0}, &Span{0, 0, 0}},
				},
				{
					{&Span{0, 0, 5}, &Span{0, 0, 8}},
					{&Span{0, 0, 6}, &Span{0, 0, 5}},
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
				{
					{&Span{0, 0, 22}, &Span{0, 0, 20}},
				},
			},
		},
		{
			"single elements at start and end, unsorted",
			fields{theta: 3},
			args{seeds: []*Seed{
				{&Span{0, 0, 10}, &Span{0, 0, 11}},
				{&Span{0, 0, 0}, &Span{0, 0, 0}},
				{&Span{0, 0, 5}, &Span{0, 0, 4}},
				{&Span{0, 0, 22}, &Span{0, 0, 20}},
				{&Span{0, 0, 6}, &Span{0, 0, 5}},
				{&Span{0, 0, 9}, &Span{0, 0, 3}},
			}},
			[][]*Seed{
				{
					{&Span{0, 0, 0}, &Span{0, 0, 0}},
				},
				{
					{&Span{0, 0, 5}, &Span{0, 0, 4}},
					{&Span{0, 0, 6}, &Span{0, 0, 5}},
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
				{
					{&Span{0, 0, 22}, &Span{0, 0, 20}},
				},
			},
		},
		{
			"one group only",
			fields{theta: 3},
			args{seeds: []*Seed{
				{&Span{0, 0, 5}, &Span{0, 0, 4}},
				{&Span{0, 0, 5}, &Span{0, 0, 8}},
				{&Span{0, 0, 6}, &Span{0, 0, 5}},
				{&Span{0, 0, 9}, &Span{0, 0, 3}},
				{&Span{0, 0, 10}, &Span{0, 0, 11}},
			}},
			[][]*Seed{
				{
					{&Span{0, 0, 5}, &Span{0, 0, 4}},
					{&Span{0, 0, 5}, &Span{0, 0, 8}},
					{&Span{0, 0, 6}, &Span{0, 0, 5}},
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
			},
		},
		{
			"each seed in own group",
			fields{theta: 1},
			args{seeds: []*Seed{
				{&Span{0, 0, 5}, &Span{0, 0, 4}},
				{&Span{0, 0, 7}, &Span{0, 0, 8}},
				{&Span{0, 0, 9}, &Span{0, 0, 5}},
				{&Span{0, 0, 11}, &Span{0, 0, 3}},
				{&Span{0, 0, 13}, &Span{0, 0, 11}},
			}},
			[][]*Seed{
				{{&Span{0, 0, 5}, &Span{0, 0, 4}}},
				{{&Span{0, 0, 7}, &Span{0, 0, 8}}},
				{{&Span{0, 0, 9}, &Span{0, 0, 5}}},
				{{&Span{0, 0, 11}, &Span{0, 0, 3}}},
				{{&Span{0, 0, 13}, &Span{0, 0, 11}}},
			},
		},
		{
			"empty input",
			fields{theta: 3},
			args{seeds: []*Seed{}},
			[][]*Seed{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RangeExtender{
				theta: tt.fields.theta,
			}
			if gotPassages := e.candidatePassages(tt.args.seeds); !reflect.DeepEqual(gotPassages, tt.wantPassages) {
				t.Errorf("candidatePassages() = %v, want %v", gotPassages, tt.wantPassages)
			}
		})
	}
}

func TestRangeExtender_concatenatePassage(t *testing.T) {
	type fields struct {
		theta uint32
	}
	type args struct {
		sourcePassages [][]*Seed
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPassages []*Seed
	}{
		{
			"pair concantenation, one group",
			fields{theta: 3},
			args{sourcePassages: [][]*Seed{
				{
					{&Span{3, 5, 5}, &Span{3, 4, 1}},
					{&Span{6, 7, 7}, &Span{1, 2, 0}},
				},
			}},
			[]*Seed{
				{&Span{3, 7, 5}, &Span{1, 4, 0}},
			},
		},
		{
			name:   "pair concantenation, multiple groups",
			fields: fields{theta: 3},
			args: args{sourcePassages: [][]*Seed{
				{
					{&Span{3, 5, 5}, &Span{3, 4, 1}},
					{&Span{6, 7, 7}, &Span{1, 2, 0}},
				},
				{
					{&Span{3, 5, 5}, &Span{3, 4, 1}},
					{&Span{6, 7, 7}, &Span{1, 2, 0}},
				},
				{
					{&Span{3, 5, 5}, &Span{3, 4, 1}},
					{&Span{6, 7, 7}, &Span{1, 2, 0}},
				},
			}},
			wantPassages: []*Seed{
				{&Span{3, 7, 5}, &Span{1, 4, 0}},
				{&Span{3, 7, 5}, &Span{1, 4, 0}},
				{&Span{3, 7, 5}, &Span{1, 4, 0}},
			},
		},
		{
			"multiple concantenation",
			fields{theta: 3},
			args{sourcePassages: [][]*Seed{
				{
					{&Span{3, 5, 5}, &Span{3, 4, 1}},
					{&Span{6, 7, 7}, &Span{1, 2, 0}},
					{&Span{8, 10, 9}, &Span{6, 9, 4}},
					{&Span{13, 16, 11}, &Span{5, 6, 3}},
				},
			}},
			[]*Seed{
				{&Span{3, 16, 5}, &Span{1, 9, 0}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RangeExtender{
				theta: tt.fields.theta,
			}
			if gotPassages := e.concatenatePassage(tt.args.sourcePassages); !reflect.DeepEqual(gotPassages, tt.wantPassages) {
				t.Errorf("concatenatePassage() = %v, want %v", gotPassages, tt.wantPassages)
			}
		})
	}
}

func TestRangeExtender_splitCandidates(t *testing.T) {
	type fields struct {
		theta uint32
	}
	type args struct {
		candidates [][]*Seed
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantPassages [][]*Seed
	}{
		{
			"single element groups",
			fields{theta: 3},
			args{candidates: [][]*Seed{
				{{&Span{0, 0, 0}, &Span{0, 0, 0}}},
				{{&Span{0, 0, 5}, &Span{0, 0, 4}}},
				{{&Span{0, 0, 22}, &Span{0, 0, 20}}},
			}},
			[][]*Seed{
				{{&Span{0, 0, 0}, &Span{0, 0, 0}}},
				{{&Span{0, 0, 5}, &Span{0, 0, 4}}},
				{{&Span{0, 0, 22}, &Span{0, 0, 20}}},
			},
		},
		{
			"split group",
			fields{theta: 3},
			args{candidates: [][]*Seed{
				{
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 5}, &Span{0, 0, 4}},
					{&Span{0, 0, 6}, &Span{0, 0, 5}},
					{&Span{0, 0, 5}, &Span{0, 0, 9}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
			}},
			[][]*Seed{
				{
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 5}, &Span{0, 0, 4}},
					{&Span{0, 0, 6}, &Span{0, 0, 5}},
				},
				{
					{&Span{0, 0, 5}, &Span{0, 0, 9}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
			},
		},
		{
			"no splits",
			fields{theta: 3},
			args{candidates: [][]*Seed{
				{
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 5}, &Span{0, 0, 5}},
					{&Span{0, 0, 5}, &Span{0, 0, 8}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
			}},
			[][]*Seed{
				{
					{&Span{0, 0, 9}, &Span{0, 0, 3}},
					{&Span{0, 0, 5}, &Span{0, 0, 5}},
					{&Span{0, 0, 5}, &Span{0, 0, 8}},
					{&Span{0, 0, 10}, &Span{0, 0, 11}},
				},
			},
		},
		{
			"empty input",
			fields{theta: 3},
			args{candidates: [][]*Seed{}},
			[][]*Seed{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RangeExtender{
				theta: tt.fields.theta,
			}
			if gotPassages := e.splitCandidates(tt.args.candidates); !reflect.DeepEqual(gotPassages, tt.wantPassages) {
				t.Errorf("splitCandidates() = %v, want %v", gotPassages, tt.wantPassages)
			}
		})
	}
}
