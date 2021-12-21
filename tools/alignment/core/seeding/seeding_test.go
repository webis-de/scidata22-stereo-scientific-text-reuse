/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"picapica.org/picapica4/alignment/core/filtering"
	"picapica.org/picapica4/alignment/core/similarity"
	"reflect"
	"testing"
)

func TestNewBaseSeeder(t *testing.T) {
	type args struct {
		similarity similarity.Similarity
	}
	tests := []struct {
		name string
		args args
		want BaseSeeder
	}{
		{
			"generator test",
			args{similarity: similarity.NewIdentitySimilarity()},
			BaseSeeder{
				PreFilterPipeline:  filtering.NewFilterPipeline(),
				SeedFilterPipeline: filtering.NewFilterPipeline(),
				Similarity:         similarity.NewIdentitySimilarity(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseSeeder(tt.args.similarity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseSeeder() = %v, want %v", got, tt.want)
			}
		})
	}
}
