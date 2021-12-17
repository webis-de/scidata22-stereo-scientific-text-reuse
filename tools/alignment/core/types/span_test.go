/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import "testing"

func TestSpan_String(t *testing.T) {
	type fields struct {
		Begin uint32
		End   uint32
		Index uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"base test",
			fields{
				Begin: 10,
				End:   92,
				Index: 5,
			},
			"Begin: 10, End: 92, Index: 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Span{
				Begin: tt.fields.Begin,
				End:   tt.fields.End,
				Index: tt.fields.Index,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpan_Length(t *testing.T) {
	type fields struct {
		Begin uint32
		End   uint32
		Index uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			"length test zero",
			fields{
				Begin: 0,
				End:   0,
				Index: 0,
			},
			0,
		},
		{
			"length test positive",
			fields{
				Begin: 51,
				End:   56,
				Index: 0,
			},
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Span{
				Begin: tt.fields.Begin,
				End:   tt.fields.End,
				Index: tt.fields.Index,
			}
			if got := m.Length(); got != tt.want {
				t.Errorf("Length() = %v, want %v", got, tt.want)
			}
		})
	}
}
