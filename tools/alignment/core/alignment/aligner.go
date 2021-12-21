/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package alignment

import (
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	. "picapica.org/picapica4/alignment/core/types"
)

type Aligner interface {
	Align([]*Element, []*Element) ([]*Seed, error)
}

type BaseAligner struct {
	Seeder   seeding.Seeder     `json:"seeder"`
	Extender extending.Extender `json:"extender"`
	Reducer  reducing.Reducer   `json:"reducer"`
}

func NewBaseAligner(seeder seeding.Seeder, extender extending.Extender, reducer reducing.Reducer) *BaseAligner {
	return &BaseAligner{
		Seeder:   seeder,
		Extender: extender,
		Reducer:  reducer,
	}
}

func (a BaseAligner) Align(source []*Element, target []*Element) ([]*Seed, error) {
	seeds, err := a.Seeder.Seed(source, target)
	if err != nil {
		return nil, err
	}
	seeds, err = a.Extender.Extend(seeds)
	if err != nil {
		return nil, err
	}
	seeds, err = a.Reducer.Process(seeds)
	if err != nil {
		return nil, err
	}
	return seeds, nil
}
