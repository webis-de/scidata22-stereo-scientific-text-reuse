/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	"fmt"
	"hash/maphash"
	"math"
	"picapica.org/picapica4/alignment/core/errors"
	. "picapica.org/picapica4/alignment/core/types"
)

// Common container for mapping seeds to a numerical space for clustering
type Clusterable interface {
	Distance(c Clusterable) (float64, error)
	ID() uint64
}

type Cluster interface {
	FromSeeds(seeds []*Seed)
	ToSeeds() []*Seed
	Clusterables() []Clusterable
	Lookup(clusterables ...Clusterable) (seeds []*Seed)
}

// Clusterable type holding a numerical value, for testing purposes
type ValueClusterable struct {
	value float64
	id    uint64
}

func (c ValueClusterable) ID() uint64 {
	return c.id
}

// Implements a distance function as value difference between two
func (c ValueClusterable) Distance(other Clusterable) (float64, error) {
	switch t := other.(type) {
	case ValueClusterable:
		return math.Abs(other.(ValueClusterable).value - c.value), nil
	default:
		return 0, errors.NewTypeError(t)
	}
}

// Clusterable type based on the indices of a seeds' spans
type IndexClusterable struct {
	x  uint32
	y  uint32
	id uint64
}

func (c IndexClusterable) ID() uint64 {
	return c.id
}

// Implements Euclidean Distance between two IndexClusterables
func (c IndexClusterable) Distance(other Clusterable) (float64, error) {
	switch t := other.(type) {
	case IndexClusterable:
		dx := float64(c.x) - float64(other.(IndexClusterable).x)
		dy := float64(c.y) - float64(other.(IndexClusterable).y)
		return math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2)), nil
	default:
		return 0, errors.NewTypeError(t)
	}

}

// Cluster type that is based on the IndexClusterables
type IndexCluster struct {
	clusterables []Clusterable
	mapping      map[uint64]*Seed
}

// Constructor for IndexCluster type
func NewIndexCluster() *IndexCluster {
	return &IndexCluster{
		clusterables: make([]Clusterable, 0, 0),
		mapping:      make(map[uint64]*Seed),
	}
}

// Initializes the clusters' elements from a list of seeds, creating IndexClusterables based on the supplied seeds' indices
func (c *IndexCluster) FromSeeds(seeds []*Seed) {
	c.mapping = make(map[uint64]*Seed)
	hash := maphash.Hash{}
	for _, seed := range seeds {
		_, _ = hash.Write([]byte(fmt.Sprintf("%v", seed)))
		key := hash.Sum64()
		hash.Reset()
		c.mapping[key] = seed
	}

	for key, seed := range c.mapping {
		c.clusterables = append(c.clusterables, IndexClusterable{
			x:  seed.SourceSpan.Index,
			y:  seed.TargetSpan.Index,
			id: key,
		})
	}
}

// Returns the clusters elements as seeds
func (c *IndexCluster) ToSeeds() (seeds []*Seed) {
	for _, clusterable := range c.clusterables {
		seeds = append(seeds, c.mapping[clusterable.ID()])
	}
	return
}

// Returns the corresponding seeds from the clusters internal mapping of a given list of clusterables
func (c *IndexCluster) Lookup(clusterables ...Clusterable) (seeds []*Seed) {
	for _, clusterable := range clusterables {
		seeds = append(seeds, c.mapping[clusterable.ID()])
	}
	return
}

func (c *IndexCluster) Clusterables() []Clusterable {
	return c.clusterables
}
