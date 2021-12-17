/*******************************************************************************
 * Last modified: 02.09.21, 18:56
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package normalizing

import (
	"reflect"
	"testing"
)

func TestNewIdentityNormalizer(t *testing.T) {
	tests := []struct {
		name string
		want *IdentityNormalizer
	}{
		{
			"constructor test",
			&IdentityNormalizer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdentityNormalizer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIdentityNormalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSortingNormalizer(t *testing.T) {
	tests := []struct {
		name string
		want *SortingNormalizer
	}{
		{
			"constructor test",
			&SortingNormalizer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSortingNormalizer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSortingNormalizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortingNormalizer_Normalize(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"string sorting",
			args{
				input: []string{"B", "A", "D", "C"},
			},
			[]string{"A", "B", "C", "D"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &SortingNormalizer{}
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
			default:
				t.Errorf("unexpected type %T", inputType)
			}
			var ss []interface{}
			switch inputType := tt.args.input.(type) {
			case []string:
				tmp := make([]interface{}, len(tt.want.([]string)), len(tt.want.([]string)))
				for i, v := range tt.want.([]string) {
					tmp[i] = v
				}
				ss = tmp
			default:
				t.Errorf("unexpected type %T", inputType)
			}

			got, err := n.Normalize(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, ss) {
				t.Errorf("Normalize() got = %v, want %v", got, ss)
			}
		})
	}
}

func TestIdentityNormalizer_Normalize(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"string values",
			args{
				input: []string{"B", "A", "D", "C"},
			},
			[]string{"B", "A", "D", "C"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &IdentityNormalizer{}
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
			default:
				t.Errorf("unexpected type %T", inputType)
			}
			var ss []interface{}
			switch inputType := tt.args.input.(type) {
			case []string:
				tmp := make([]interface{}, len(tt.want.([]string)), len(tt.want.([]string)))
				for i, v := range tt.want.([]string) {
					tmp[i] = v
				}
				ss = tmp
			default:
				t.Errorf("unexpected type %T", inputType)
			}

			got, err := n.Normalize(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, ss) {
				t.Errorf("Normalize() got = %v, want %v", got, ss)
			}
		})
	}
}
