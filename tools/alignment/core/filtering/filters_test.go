/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package filtering

import (
	. "picapica.org/picapica4/alignment/core/types"
	"testing"
)

func TestIdentityFilter(t *testing.T) {
	type args struct {
		input *Element
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Base test",
			args{
				input: &Element{Span: &Span{Begin: 0, End: 3, Index: 0}, Value: "test"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IdentityFilter(tt.args.input); got != tt.want {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoScalarValueFilter(t *testing.T) {
	type args struct {
		input *Element
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"non-scalar element",
			args{
				input: &Element{Span: &Span{Begin: 0, End: 3, Index: 0}, Value: "test"},
			},
			true,
		},
		{
			"scalar element",
			args{
				input: &Element{Span: &Span{Begin: 0, End: 0, Index: 0}, Value: "t"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoScalarValueFilter(tt.args.input); got != tt.want {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
