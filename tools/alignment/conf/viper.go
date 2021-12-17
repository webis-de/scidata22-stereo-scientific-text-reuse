/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package conf

import (
	"github.com/spf13/viper"
	"picapica.org/picapica4/alignment/core/alignment"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/normalizing"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	"picapica.org/picapica4/alignment/core/similarity"
)

func AlignerFromViper(cfg *viper.Viper) alignment.Aligner {
	return alignment.NewBaseAligner(
		seederFromViper(cfg),
		extenderFromViper(cfg),
		reducerFromViper(cfg),
	)
}

func modeFromViper(cfg *viper.Viper) string {
	// Base case if config des not specify anything
	if !cfg.IsSet("mode") {
		return BaseConfig.Mode
	}
	// Parse config into object
	switch cfg.GetString("mode") {
	case "STRING":
		return "string"
	case "VECTOR":
		return "vector"
	default:
		return BaseConfig.Mode
	}
}

func seederFromViper(cfg *viper.Viper) seeding.Seeder {
	// Base case if config does not specify anything
	if !cfg.IsSet("seeder") {
		return BaseConfig.Seeder
	}
	// Parse config into object
	if cfg.Sub("seeder").IsSet("ngram") {
		return seeding.NewNGramSeeder(
			cfg.Sub("seeder").Sub("ngram").GetUint32("n"),
			cfg.Sub("seeder").Sub("ngram").GetUint32("overlap"),
			similarityFromViper(cfg),
			combinatorFromViper(cfg),
		)
	} else if cfg.Sub("seeder").IsSet("hash") {
		return seeding.NewHashSeeder(
			cfg.Sub("seeder").Sub("hash").GetUint32("n"),
			cfg.Sub("seeder").Sub("hash").GetUint32("overlap"),
			normalizerFromViper(cfg),
			combinatorFromViper(cfg),
		)
	} else if cfg.Sub("seeder").IsSet("smithwaterman") {
		return seeding.NewSmithWatermanSeeder(
			similarityFromViper(cfg),
			cfg.Sub("seeder").Sub("smithwaterman").GetFloat64("mismatchcost"),
			cfg.Sub("seeder").Sub("smithwaterman").GetFloat64("gapcost"),
			cfg.Sub("seeder").Sub("smithwaterman").GetFloat64("matchaward"),
		)
	} else if cfg.Sub("seeder").IsSet("needlemanwunsch") {
		return seeding.NewNeedlemanWunschSeeder(
			similarityFromViper(cfg),
			cfg.Sub("seeder").Sub("needlemanwunsch").GetFloat64("mismatchcost"),
			cfg.Sub("seeder").Sub("needlemanwunsch").GetFloat64("gapcost"),
			cfg.Sub("seeder").Sub("needlemanwunsch").GetFloat64("matchaward"),
		)
	} else {
		return BaseConfig.Seeder
	}
}

func extenderFromViper(cfg *viper.Viper) extending.Extender {
	// Base case if config does not specify anything
	if !cfg.IsSet("extender") {
		return BaseConfig.Extender
	}
	// Parse config into object
	if cfg.Sub("extender").IsSet("density") {
		return extending.NewDensityExtender(
			cfg.Sub("extender").Sub("density").GetInt("minpoints"),
			cfg.Sub("extender").Sub("density").GetFloat64("epsilon"),
		)
	} else if cfg.Sub("extender").IsSet("range") {
		return extending.NewRangeExtender(cfg.Sub("extender").Sub("range").GetUint32("theta"))
	} else {
		return BaseConfig.Extender
	}
}

func reducerFromViper(cfg *viper.Viper) reducing.Reducer {
	// Base case if config does not specify anything
	if !cfg.IsSet("reducer") {
		return BaseConfig.Reducer
	}
	// Parse config into objects
	if cfg.Sub("reducer").IsSet("passagelength") {
		return reducing.NewPassageLengthReducer(cfg.Sub("reducer").Sub("passagelength").GetUint32("minlength"))
	} else {
		return BaseConfig.Reducer
	}
}

func similarityFromViper(cfg *viper.Viper) similarity.Similarity {
	// Base case if config does not specify anything
	if !cfg.IsSet("similarity") {
		return BaseConfig.Similarity
	}
	// Parse config into object
	if cfg.Sub("similarity").IsSet("identity") {
		return similarity.NewIdentitySimilarity()
	} else if cfg.Sub("similarity").IsSet("cosine") {
		return similarity.NewCosineVectorSimilarity(cfg.Sub("similarity").Sub("cosine").GetFloat64("threshold"))
	} else if cfg.Sub("similarity").IsSet("levenshtein") {
		return similarity.NewLevenshteinStringSimilarity(cfg.Sub("similarity").Sub("levenshtein").GetFloat64("threshold"))
	} else if cfg.Sub("similarity").IsSet("jaccard") {
		return similarity.NewJaccardStringSimilarity(cfg.Sub("similarity").Sub("jaccard").GetFloat64("threshold"), cfg.Sub("similarity").Sub("jaccard").GetString("separator"))
	} else {
		return BaseConfig.Similarity
	}
}

func normalizerFromViper(cfg *viper.Viper) normalizing.Normalizer {
	// Base case if config does not specify anything
	if !cfg.IsSet("normalizer") {
		return BaseConfig.Normalizer
	}
	// Parse config into object
	if cfg.Sub("normalizer").IsSet("sorting") {
		return normalizing.NewSortingNormalizer()
	} else {
		return BaseConfig.Normalizer
	}
}

func combinatorFromViper(cfg *viper.Viper) combinating.ValueCombinator {
	// Base case if config does not specify anything
	if !cfg.IsSet("combinator") {
		return BaseConfig.Combinator
	}
	// Parse config into object
	if cfg.Sub("combinator").IsSet("stringconcat") {
		return combinating.NewStringValueCombinator(cfg.Sub("combinator").Sub("stringconcat").GetString("sep"))
	} else if cfg.Sub("combinator").IsSet("vectorconcat") {
		return combinating.NewVectorValueCombinator()
	} else {
		return BaseConfig.Combinator
	}
}
