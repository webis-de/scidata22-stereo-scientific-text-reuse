/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package seeding

import (
	"picapica.org/picapica4/alignment/core/filtering"
	"picapica.org/picapica4/alignment/core/similarity"
	. "picapica.org/picapica4/alignment/core/types"
)

type Seeder interface {
	Seed([]*Element, []*Element) ([]*Seed, error)
}

type BaseSeeder struct {
	PreFilterPipeline  filtering.FilterPipeline `json:"filters,omitempty"`
	SeedFilterPipeline filtering.FilterPipeline `json:"seed_filters,omitempty"`
	Similarity         similarity.Similarity    `json:"similarity"`
}

func NewBaseSeeder(similarity similarity.Similarity) BaseSeeder {
	return BaseSeeder{
		PreFilterPipeline:  filtering.NewFilterPipeline(),
		SeedFilterPipeline: filtering.NewFilterPipeline(),
		Similarity:         similarity,
	}
}
