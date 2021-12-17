/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import "testing"

func TestElement_String(t *testing.T) {
	type fields struct {
		Span  *Span
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"base test",
			fields{
				Span: &Span{
					Begin: 0,
					End:   9,
					Index: 0,
				},
				Value: "teststring",
			},
			"Span: [Begin: 0, End: 9, Index: 0], Value: teststring",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Element{
				Span:  tt.fields.Span,
				Value: tt.fields.Value,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
