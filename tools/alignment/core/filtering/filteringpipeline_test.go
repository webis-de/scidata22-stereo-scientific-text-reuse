/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package filtering

import (
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func TestFilterPipeline_All(t *testing.T) {
	type fields struct {
		filters []Filter
	}
	type args struct {
		element *Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"identity filter test",
			fields{filters: []Filter{
				IdentityFilter,
				IdentityFilter,
			}},
			args{
				element: &Element{Span: &Span{Begin: 0, End: 3, Index: 0}, Value: "test"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterPipeline{}
			f.Register(tt.fields.filters...)
			if got := f.All(tt.args.element); got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterPipeline_Reduce(t *testing.T) {
	type fields struct {
		pipeline FilterPipeline
	}
	type args struct {
		elements []*Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Element
	}{
		{
			"identity pipeline",
			fields{
				pipeline: FilterPipeline{IdentityFilter},
			},
			args{elements: []*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
				{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
			}},
			[]*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
				{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
			},
		},
		{
			"no-scalar pipeline",
			fields{
				pipeline: FilterPipeline{NoScalarValueFilter},
			},
			args{elements: []*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
				{Span: &Span{Begin: 9, End: 9, Index: 2}, Value: "a"},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
			}},
			[]*Element{
				{Span: &Span{Begin: 1, End: 4, Index: 0}, Value: "This"},
				{Span: &Span{Begin: 6, End: 7, Index: 1}, Value: "is"},
				{Span: &Span{Begin: 11, End: 14, Index: 3}, Value: "test"},
				{Span: &Span{Begin: 16, End: 21, Index: 4}, Value: "string"},
			},
		},
		{
			"empty sequence",
			fields{
				pipeline: FilterPipeline{IdentityFilter},
			},
			args{elements: []*Element{}},
			[]*Element{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.fields.pipeline
			if got := f.Reduce(tt.args.elements...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFilterPipeline(t *testing.T) {
	tests := []struct {
		name string
		want FilterPipeline
	}{
		{
			"constructor test",
			make([]Filter, 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFilterPipeline(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilterPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}
