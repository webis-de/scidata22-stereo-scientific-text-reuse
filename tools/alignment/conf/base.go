/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package conf

import (
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/normalizing"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	"picapica.org/picapica4/alignment/core/similarity"
)

// Provides structured access to the base cases of each alignment component
// Should be used to fill unspecified components in a conf
var BaseConfig = struct {
	Mode       string
	Similarity similarity.Similarity
	Seeder     seeding.Seeder
	Extender   extending.Extender
	Reducer    reducing.Reducer
	Normalizer normalizing.Normalizer
	Combinator combinating.ValueCombinator
}{
	"string",
	similarity.NewIdentitySimilarity(),
	seeding.NewIdentitySeeder(),
	extending.NewIdentityExtender(),
	reducing.NewIdentityReducer(),
	normalizing.NewIdentityNormalizer(),
	combinating.NewStringValueCombinator(""),
}
