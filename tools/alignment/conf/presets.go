/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package conf

import (
	"encoding/json"
	"github.com/spf13/viper"
	"picapica.org/picapica4/alignment/core/alignment"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/normalizing"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	"picapica.org/picapica4/alignment/core/similarity"
)

type Preset struct {
	Mode       string                `json:"mode"`
	Similarity similarity.Similarity `json:"similarity"`
	Seeder     seeding.Seeder        `json:"seeder"`
	Extender   extending.Extender    `json:"extender"`
	Reducer    reducing.Reducer      `json:"reducer"`
}

func (p *Preset) Json() ([]byte, error) {
	marshalled, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return marshalled, nil
}

var DefaultPreset = Preset{
	"string",
	similarity.NewIdentitySimilarity(),
	seeding.NewHashSeeder(2, 1, normalizing.NewIdentityNormalizer(), combinating.NewStringValueCombinator(" ")),
	extending.NewIdentityExtender(),
	reducing.NewIdentityReducer(),
}

func AlignerFromPreset(preset Preset) alignment.Aligner {
	return alignment.NewBaseAligner(
		preset.Seeder,
		preset.Extender,
		preset.Reducer,
	)
}

func PresetFromViper(cfg *viper.Viper) Preset {
	return Preset{
		Mode:       modeFromViper(cfg),
		Similarity: similarityFromViper(cfg),
		Seeder:     seederFromViper(cfg),
		Extender:   extenderFromViper(cfg),
		Reducer:    reducerFromViper(cfg),
	}
}
