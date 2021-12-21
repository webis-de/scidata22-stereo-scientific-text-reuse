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

func TestIdentityExtender_Extend(t *testing.T) {
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
			"base test",
			args{seeds: []*Seed{
				{
					&Span{
						Begin: 0,
						End:   1,
						Index: 0,
					},
					&Span{
						Begin: 0,
						End:   1,
						Index: 0,
					},
				},
				{
					&Span{
						Begin: 1,
						End:   2,
						Index: 1,
					},
					&Span{
						Begin: 1,
						End:   2,
						Index: 1,
					},
				},
			}},
			[]*Seed{
				{
					&Span{
						Begin: 0,
						End:   1,
						Index: 0,
					},
					&Span{
						Begin: 0,
						End:   1,
						Index: 0,
					},
				},
				{
					&Span{
						Begin: 1,
						End:   2,
						Index: 1,
					},
					&Span{
						Begin: 1,
						End:   2,
						Index: 1,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &IdentityExtender{}
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

func TestNewIdentityExtender(t *testing.T) {
	tests := []struct {
		name string
		want *IdentityExtender
	}{
		{
			"constructor test",
			&IdentityExtender{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentityExtender(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentityExtender() = %v, want %v", got, tt.want)
			}
		})
	}
}
