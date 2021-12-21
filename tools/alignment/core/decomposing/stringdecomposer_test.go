/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package decomposing

import (
	. "picapica.org/picapica4/alignment/core/types"
	"reflect"
	"testing"
)

func TestNewStringDecomposer(t *testing.T) {
	type args struct {
		wordDelimiter     string
		sentenceDelimiter string
	}
	tests := []struct {
		name string
		args args
		want *StringDecomposer
	}{
		{
			"Test for equality",
			args{
				" ",
				".",
			},
			&StringDecomposer{
				wordDelimiter:        " ",
				lenWordDelimiter:     1,
				sentenceDelimiter:    ".",
				lenSentenceDelimiter: 1,
			},
		},
		{
			"Test for longer separators",
			args{
				"  ",
				"<SEP>",
			},
			&StringDecomposer{
				wordDelimiter:        "  ",
				lenWordDelimiter:     2,
				sentenceDelimiter:    "<SEP>",
				lenSentenceDelimiter: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringDecomposer(tt.args.wordDelimiter, tt.args.sentenceDelimiter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStringDecomposer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDecomposer_Decompose(t1 *testing.T) {
	type fields struct {
		wordDelimiter        string
		lenWordDelimiter     int
		sentenceDelimiter    string
		lenSentenceDelimiter int
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSeq []*Element
		wantErr bool
	}{
		{
			"Test for two sentences",
			fields{
				wordDelimiter:        " ",
				lenWordDelimiter:     1,
				sentenceDelimiter:    ".",
				lenSentenceDelimiter: 1,
			},
			args{input: "This is. a test string."},
			[]*Element{
				{
					Span: &Span{
						Begin: 0,
						End:   4,
						Index: 0,
					},
					Value: "This",
				},
				{
					Span: &Span{
						Begin: 4,
						End:   7,
						Index: 1,
					},
					Value: "is",
				},
				{
					Span: &Span{
						Begin: 7,
						End:   10,
						Index: 2,
					},
					Value: "a",
				},
				{
					Span: &Span{
						Begin: 10,
						End:   15,
						Index: 3,
					},
					Value: "test",
				},
				{
					Span: &Span{
						Begin: 15,
						End:   22,
						Index: 4,
					},
					Value: "string",
				},
			},
			false,
		},
		{
			"Test for one sentences",
			fields{
				wordDelimiter:        " ",
				lenWordDelimiter:     1,
				sentenceDelimiter:    ".",
				lenSentenceDelimiter: 1,
			},
			args{input: "This is a test string."},
			[]*Element{
				{
					Span: &Span{
						Begin: 0,
						End:   4,
						Index: 0,
					},
					Value: "This",
				},
				{
					Span: &Span{
						Begin: 4,
						End:   7,
						Index: 1,
					},
					Value: "is",
				},
				{
					Span: &Span{
						Begin: 7,
						End:   9,
						Index: 2,
					},
					Value: "a",
				},
				{
					Span: &Span{
						Begin: 9,
						End:   14,
						Index: 3,
					},
					Value: "test",
				},
				{
					Span: &Span{
						Begin: 14,
						End:   21,
						Index: 4,
					},
					Value: "string",
				},
			},
			false,
		},
		{
			"Test for single word",
			fields{
				wordDelimiter:        " ",
				lenWordDelimiter:     1,
				sentenceDelimiter:    ".",
				lenSentenceDelimiter: 1,
			},
			args{input: "This."},
			[]*Element{
				{
					Span: &Span{
						Begin: 0,
						End:   4,
						Index: 0,
					},
					Value: "This",
				},
			},
			false,
		},
		{
			"Test for empty string",
			fields{
				wordDelimiter:        " ",
				lenWordDelimiter:     1,
				sentenceDelimiter:    ".",
				lenSentenceDelimiter: 1,
			},
			args{input: ""},
			[]*Element{},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := StringDecomposer{
				wordDelimiter:        tt.fields.wordDelimiter,
				lenWordDelimiter:     tt.fields.lenWordDelimiter,
				sentenceDelimiter:    tt.fields.sentenceDelimiter,
				lenSentenceDelimiter: tt.fields.lenSentenceDelimiter,
			}
			gotSeq, err := t.Decompose(tt.args.input)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Decompose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSeq, tt.wantSeq) {
				t1.Errorf("Decompose() gotSeq = %v, want %v", gotSeq, tt.wantSeq)
			}
		})
	}
}
