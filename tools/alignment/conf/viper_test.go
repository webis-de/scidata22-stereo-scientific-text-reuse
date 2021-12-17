/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"picapica.org/picapica4/alignment/core/alignment"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/normalizing"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	"picapica.org/picapica4/alignment/core/similarity"
	"reflect"
	"testing"
)

func Test_normalizerFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want normalizing.Normalizer
	}{
		{
			"base case",
			args{
				config: []byte(`{"normalizer":{"base":{}}"`),
			},
			normalizing.NewIdentityNormalizer(),
		},
		{
			"SortingNormalizer",
			args{
				config: []byte(`{"normalizer":{"sorting":{}}`),
			},
			normalizing.NewSortingNormalizer(),
		},
	}
	for _, tt := range tests {
		cfg := &viper.Viper{}
		cfg.SetConfigType("json")
		_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizerFromViper(cfg); &got == &tt.want {
				t.Errorf("SelectorFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_combinatorFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want combinating.ValueCombinator
	}{
		{
			"StringConcat",
			args{
				config: []byte(`{"combinator":{"stringconcat":{"sep": " "}}}`),
			},
			combinating.NewStringValueCombinator(" "),
		},
		{
			"VectorConcat",
			args{
				config: []byte(`{"combinator":{"vectorconcat":{}}}`),
			},
			combinating.NewVectorValueCombinator(),
		},
		{
			"base case",
			args{
				config: []byte(`{"combinator":{"base":{}}}`),
			},
			BaseConfig.Combinator,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
			if got := combinatorFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combinatorFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_similarityFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want similarity.Similarity
	}{
		{
			"no similarity specified",
			args{
				[]byte(`"similarity:{}"`),
			},
			BaseConfig.Similarity,
		},
		{
			"identity similarity",
			args{
				[]byte(`{"similarity":{"identity":{}}}`),
			},
			similarity.NewIdentitySimilarity(),
		},
		{
			"cosine similarity",
			args{
				[]byte(`{"similarity":{"cosine":{"threshold":0.2}}}`),
			},
			similarity.NewCosineVectorSimilarity(0.2),
		},
		{
			"levenshtein similarity",
			args{
				[]byte(`{"similarity":{"levenshtein":{"threshold":1.5}}}`),
			},
			similarity.NewLevenshteinStringSimilarity(1.5),
		},
		{
			"jaccard similarity",
			args{
				[]byte(`{"similarity":{"jaccard":{"threshold":0.9, "separator": " "}}}`),
			},
			similarity.NewJaccardStringSimilarity(0.9, " "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
			if got := similarityFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("similarityFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reducerFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want reducing.Reducer
	}{
		{
			"no reducer specified",
			args{
				config: []byte(`{"reducer":{}}`),
			},
			BaseConfig.Reducer,
		},
		{
			"PassageLength",
			args{
				config: []byte(`{"reducer":{"passagelength":{"minlength":100}}}`),
			},
			reducing.NewPassageLengthReducer(100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
			if got := reducerFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reducerFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extenderFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want extending.Extender
	}{
		{
			"no extender specified",
			args{
				config: []byte(`{"extender":{}}`),
			},
			BaseConfig.Extender,
		},
		{
			"Density",
			args{
				config: []byte(`{"extender":{"density":{"minpoints":5,"epsilon":20}}}`),
			},
			extending.NewDensityExtender(5, 20),
		},
		{
			"Range",
			args{
				config: []byte(`{"extender":{"range":{"theta":10}}}`),
			},
			extending.NewRangeExtender(10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
			if got := extenderFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extenderFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seederFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want seeding.Seeder
	}{
		{
			"no seeder specified",
			args{
				config: []byte(`{"seeder":{}}`),
			},
			BaseConfig.Seeder,
		},
		{
			"NeedlemanWunsch",
			args{
				config: []byte(`{"seeder":{"needlemanwunsch":{"mismatchcost":-1,"gapcost":-1, "matchaward":1}}, "similarity":{"identity":{}}}`),
			},
			seeding.NewNeedlemanWunschSeeder(
				similarity.NewIdentitySimilarity(),
				-1,
				-1,
				1,
			),
		},
		{
			"SmithWaterman",
			args{
				config: []byte(`{"seeder":{"smithwaterman":{"mismatchcost":-1,"gapcost":-1, "matchaward":1}}, "similarity":{"identity":{}}}`),
			},
			seeding.NewSmithWatermanSeeder(
				similarity.NewIdentitySimilarity(),
				-1,
				-1,
				1,
			),
		},
		{
			"NGram",
			args{
				config: []byte(`{"seeder":{"ngram":{"n":5,"overlap":4}},"similarity":{"identity":{}},"combinator":{"stringconcat":{"sep":" "}}}`),
			},
			seeding.NewNGramSeeder(
				5,
				4,
				similarity.NewIdentitySimilarity(),
				combinating.NewStringValueCombinator(" "),
			),
		},
		{
			"Hash",
			args{
				config: []byte(`{"seeder":{"hash":{"n":5,"overlap":4}},"similarity":{"identity":{}},"combinator":{"stringconcat":{"sep":" "}}}`),
			},
			seeding.NewHashSeeder(
				5,
				4,
				normalizing.NewIdentityNormalizer(),
				combinating.NewStringValueCombinator(" "),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
			if got := seederFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seederFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlignerFromViper(t *testing.T) {
	type args struct {
		config []byte
	}
	tests := []struct {
		name string
		args args
		want alignment.Aligner
	}{
		{
			"string aligner",
			args{
				config: []byte(`{"mode":"STRING","seeder":{"ngram":{"n":5,"overlap":4}},"extender":{"density":{"minPoints":5,"epsilon":20}},"reducer":{"passageLength":{"minLength":100}},"similarity":{"identity":{}},"combinator":{"stringConcat":{"sep":" "}}}`),
			},
			alignment.NewBaseAligner(
				seeding.NewNGramSeeder(5, 4, similarity.NewIdentitySimilarity(), combinating.NewStringValueCombinator(" ")),
				extending.NewDensityExtender(5, 20),
				reducing.NewPassageLengthReducer(100),
			),
		},
	}
	for _, tt := range tests {
		cfg := &viper.Viper{}
		cfg.SetConfigType("json")
		_ = cfg.ReadConfig(bytes.NewBuffer(tt.args.config))
		t.Run(tt.name, func(t *testing.T) {
			if got := AlignerFromViper(cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AlignerFromViper() = %v, want %v", got, tt.want)
			}
		})
	}
}
