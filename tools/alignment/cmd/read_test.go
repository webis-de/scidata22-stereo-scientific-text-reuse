/*******************************************************************************
 * Last modified: 27.01.21, 10:45
 * Author: Lukas Gienapp
 * Copyright (c) 2021 picapica.org
 */

package cmd

import (
	"picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func Test_generateResults(t *testing.T) {
	type args struct {
		seeds         []*types.Seed
		source        string
		target        string
		returnText    bool
		returnContext bool
	}
	tests := []struct {
		name string
		args args
		want []Result
	}{
		{
			"no text, no context",
			args{
				seeds: []*types.Seed{
					{
						SourceSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
						TargetSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
					},
				},
				source:        "This is the source. Short overlapping text. After source context.",
				target:        "This is the target. Short overlapping text. After target context.",
				returnText:    false,
				returnContext: false,
			},
			[]Result{
				{
					Source: Partial{
						Begin:  20,
						End:    44,
						Text:   "",
						Before: "",
						After:  "",
					},
					Target: Partial{
						Begin:  20,
						End:    44,
						Text:   "",
						Before: "",
						After:  "",
					},
				},
			},
		},
		{
			"text, no context",
			args{
				seeds: []*types.Seed{
					{
						SourceSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
						TargetSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
					},
				},
				source:        "This is the source. Short overlapping text. After source context.",
				target:        "This is the target. Short overlapping text. After target context.",
				returnText:    true,
				returnContext: false,
			},
			[]Result{
				{
					Source: Partial{
						Begin:  20,
						End:    44,
						Text:   "Short overlapping text. ",
						Before: "",
						After:  "",
					},
					Target: Partial{
						Begin:  20,
						End:    44,
						Text:   "Short overlapping text. ",
						Before: "",
						After:  "",
					},
				},
			},
		},
		{
			"text, context",
			args{
				seeds: []*types.Seed{
					{
						SourceSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
						TargetSpan: &types.Span{
							Begin: 20,
							End:   44,
							Index: 5,
						},
					},
				},
				source:        "This is the source. Short overlapping text. After source context.",
				target:        "This is the target. Short overlapping text. After target context.",
				returnText:    true,
				returnContext: true,
			},
			[]Result{
				{
					Source: Partial{
						Begin:  20,
						End:    44,
						Text:   "Short overlapping text. ",
						Before: "This is the source. ",
						After:  "After source context.",
					},
					Target: Partial{
						Begin:  20,
						End:    44,
						Text:   "Short overlapping text. ",
						Before: "This is the target. ",
						After:  "After target context.",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateResults(tt.args.seeds, tt.args.source, tt.args.target, tt.args.returnText, tt.args.returnContext); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateResults() = %v, want %v", got, tt.want)
			}
		})
	}
}
