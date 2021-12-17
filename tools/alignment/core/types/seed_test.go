/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import "testing"

func TestSeed_String(t *testing.T) {
	type fields struct {
		Span1 *Span
		Span2 *Span
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"base test",
			fields{
				Span1: &Span{
					Begin: 10,
					End:   20,
					Index: 3,
				},
				Span2: &Span{
					Begin: 9,
					End:   144,
					Index: 6,
				},
			},
			"Span 1: [Begin: 10, End: 20, Index: 3], Span 2: [Begin: 9, End: 144, Index: 6]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Seed{
				SourceSpan: tt.fields.Span1,
				TargetSpan: tt.fields.Span2,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
