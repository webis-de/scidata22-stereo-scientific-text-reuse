/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	. "picapica.org/picapica4/alignment/core/types"
)

// An Extender transforms a given SeedSet
type Extender interface {
	Extend([]*Seed) ([]*Seed, error)
}

// IdentityExtender
type IdentityExtender struct{}

func NewIdentityExtender() *IdentityExtender {
	return &IdentityExtender{}
}

func (e *IdentityExtender) Extend(seeds []*Seed) ([]*Seed, error) { return seeds, nil }
