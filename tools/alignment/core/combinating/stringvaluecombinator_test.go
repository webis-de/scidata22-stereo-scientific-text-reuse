/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package combinating

import (
	"reflect"
	"testing"
)

func TestStringValueCombinator_Combine(t *testing.T) {
	type fields struct {
		sep string
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"no error",
			fields{" "},
			args{[]string{"This", "is", "a", "test", "text"}},
			"This is a test text",
			false,
		},
		{
			"error",
			fields{" "},
			args{[]uint32{0, 2, 5, 4, 3, 3, 5}},
			make([]string, 7, 7),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &StringValueCombinator{
				sep: tt.fields.sep,
			}
			// Create interface slice to get type conversion working correctly
			// https://golang.org/doc/faq#convert_slice_of_interface
			var s []interface{}
			switch inputType := tt.args.input.(type) {
			case []string:
				tmp := make([]interface{}, len(tt.args.input.([]string)), len(tt.args.input.([]string)))
				for i, v := range tt.args.input.([]string) {
					tmp[i] = v
				}
				s = tmp
			case []uint32:
				tmp := make([]interface{}, len(tt.args.input.([]uint32)), len(tt.args.input.([]uint32)))
				for i, v := range tt.args.input.([]uint32) {
					tmp[i] = v
				}
				s = tmp
			default:
				t.Errorf("unexpected type %T", inputType)
			}
			got, err := c.Combine(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Combine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Combine() got = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestNewStringValueCombinator(t *testing.T) {
	type args struct {
		sep string
	}
	tests := []struct {
		name string
		args args
		want *StringValueCombinator
	}{
		{
			"constructor test",
			args{sep: " "},
			&StringValueCombinator{" "},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringValueCombinator(tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStringValueCombinator() = %v, want %v", got, tt.want)
			}
		})
	}
}
