/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package reducing

import (
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func TestNewPassageLengthReducer(t *testing.T) {
	type args struct {
		minLength uint32
	}
	tests := []struct {
		name string
		args args
		want *PassageLengthReducer
	}{
		{
			"Constructor test 1",
			args{minLength: 5},
			&PassageLengthReducer{5},
		},
		{
			"Constructor test 2",
			args{minLength: 10},
			&PassageLengthReducer{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPassageLengthReducer(tt.args.minLength); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPassageLengthReducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassageLengthReducer_Process(t *testing.T) {
	type fields struct {
		minLength uint32
	}
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Seed
	}{
		{
			"minlength 5",
			fields{minLength: 5},
			args{seeds: []*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 4, Index: 0},
					TargetSpan: &Span{Begin: 0, End: 5, Index: 0},
				},
				{
					SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
					TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
				},
			}},
			[]*Seed{
				{
					SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
					TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
				},
			},
		},
		{
			"minlength 0",
			fields{minLength: 0},
			args{seeds: []*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 4, Index: 0},
					TargetSpan: &Span{Begin: 0, End: 5, Index: 0},
				},
				{
					SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
					TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
				},
			}},
			[]*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 4, Index: 0},
					TargetSpan: &Span{Begin: 0, End: 5, Index: 0},
				},
				{
					SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
					TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PassageLengthReducer{
				MinLength: tt.fields.minLength,
			}
			if got, _ := p.Process(tt.args.seeds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdentityReducer(t *testing.T) {
	tests := []struct {
		name string
		want *IdentityReducer
	}{
		{
			"constructor",
			&IdentityReducer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentityReducer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentityReducer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentityReducer_Process(t *testing.T) {
	type args struct {
		seeds []*Seed
	}
	tests := []struct {
		name    string
		args    args
		want    []*Seed
		wantErr bool
	}{
		{
			"base",
			args{
				[]*Seed{
					{
						SourceSpan: &Span{Begin: 0, End: 4, Index: 0},
						TargetSpan: &Span{Begin: 0, End: 5, Index: 0},
					},
					{
						SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
						TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
					},
				},
			},
			[]*Seed{
				{
					SourceSpan: &Span{Begin: 0, End: 4, Index: 0},
					TargetSpan: &Span{Begin: 0, End: 5, Index: 0},
				},
				{
					SourceSpan: &Span{Begin: 5, End: 10, Index: 0},
					TargetSpan: &Span{Begin: 6, End: 15, Index: 0},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &IdentityReducer{}
			got, err := p.Process(tt.args.seeds)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() = %v, want %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
