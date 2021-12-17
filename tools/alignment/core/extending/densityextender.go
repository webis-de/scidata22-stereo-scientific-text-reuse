/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package extending

import (
	"math"
	. "picapica.org/picapica4/alignment/core/types"
	"strconv"
)

const (
	UNKNOWN = "unk"
	NOISE   = "noise"
)

// Extender with value-based hierarchical clustering
// Based on the DBScan algorithm
type DensityExtender struct {
	minPoints int
	epsilon   float64
}

func NewDensityExtender(minPoints int, epsilon float64) *DensityExtender {
	return &DensityExtender{minPoints: minPoints, epsilon: epsilon}
}

func (e *DensityExtender) Extend(seeds []*Seed) (extension []*Seed, err error) {

	cluster := IndexCluster{}
	cluster.FromSeeds(seeds)

	densityClusterer := NewDBScan(e.minPoints, e.epsilon)
	clusters := densityClusterer.Cluster(cluster.Clusterables())
	// Convert clusters back into seeds
	for _, c := range clusters {
		// Find enclosing seed of cluster
		var (
			beginSrc uint32 = math.MaxUint32
			endSrc   uint32 = 0
			indexSrc uint32 = math.MaxUint32
			beginTgt uint32 = math.MaxUint32
			endTgt   uint32 = 0
			indexTgt uint32 = math.MaxUint32
		)
		for _, seed := range cluster.Lookup(c...) {
			if seed.SourceSpan.Begin < beginSrc {
				beginSrc = seed.SourceSpan.Begin
			}
			if seed.SourceSpan.End > endSrc {
				endSrc = seed.SourceSpan.End
			}
			if seed.SourceSpan.Index < indexSrc {
				indexSrc = seed.SourceSpan.Index
			}
			if seed.TargetSpan.Begin < beginTgt {
				beginTgt = seed.TargetSpan.Begin
			}
			if seed.TargetSpan.End > endTgt {
				endTgt = seed.TargetSpan.End
			}
			if seed.TargetSpan.Index < indexTgt {
				indexTgt = seed.TargetSpan.Index
			}
		}
		extension = append(extension, &Seed{
			SourceSpan: &Span{Begin: beginSrc, End: endSrc, Index: indexSrc},
			TargetSpan: &Span{Begin: beginTgt, End: endTgt, Index: indexTgt},
		})
	}
	return extension, nil
}

type DBScan struct {
	clusterables []Clusterable                        // Saves the clusterables
	minPts       int                                  // Minimum points a cluster needs to have (density)
	eps          float64                              // Maximum distance two points can be apart to still be considered within the same cluster
	size         int                                  // Number of clusterables
	reachable    map[Clusterable]map[Clusterable]bool // Pairwise distance matrix of clusterables
}

// Constructor
func NewDBScan(minPts int, eps float64) *DBScan {
	return &DBScan{
		minPts: minPts,
		eps:    eps,
	}
}

// Clusters a given set of clusterables using the parameters specified in the DBScan struct
func (e *DBScan) Cluster(clusterables []Clusterable) (clusters [][]Clusterable) {
	// Setup
	e.clusterables = clusterables
	e.size = len(e.clusterables)
	e.reachable = e.DistanceMatrix(e.clusterables)
	clusterMap := make(map[string][]Clusterable)
	// Compute clusters & join clusterables with the same cluster ID into a slice
	for clusterable, clusterID := range e.computeClusters() {
		if v, ok := clusterMap[clusterID]; ok {
			clusterMap[clusterID] = append(v, clusterable)
		} else {
			clusterMap[clusterID] = []Clusterable{clusterable}
		}
	}
	// Get clusters from map, dropping noise clusterables
	for clusterID, cluster := range clusterMap {
		if clusterID != NOISE {
			clusters = append(clusters, cluster)
		}
	}
	return
}

// Calculates the symmetric pairwise distance matrix for a given set of clusterables.
// The matrix is symmetric since clusterables are required to implement a distance
// function for which d(q,p) == d(p,q) is true and
func (e *DBScan) DistanceMatrix(clusterables []Clusterable) map[Clusterable]map[Clusterable]bool {
	// Initialize matrix & prefill diagonal
	distances := make(map[Clusterable]map[Clusterable]bool)
	for _, x := range clusterables {
		distances[x] = make(map[Clusterable]bool)
		distances[x][x] = true
	}
	// Fill matrix; copies same value to each mirrored entry [i][j] and [j][i]
	for i, x := range clusterables {
		for j := i + 1; j < e.size; j++ {
			d, _ := x.Distance(clusterables[j])
			distances[x][clusterables[j]] = d <= e.eps
			distances[clusterables[j]][x] = d <= e.eps
		}
	}
	return distances
}

// Calculates the epsilon-neighborhood of a clusterable q to all other clusterables,
// i.e. a slice of all other clusterables that are reachable within epsilon distance
func (e *DBScan) Neighborhood(q Clusterable) []Clusterable {
	var neighbors []Clusterable
	for _, p := range e.clusterables {
		if e.reachable[q][p] || e.reachable[p][q] {
			neighbors = append(neighbors, p)
		}
	}
	return neighbors
}

func (e *DBScan) computeClusters() map[Clusterable]string {
	// Map containing information about the state of each clusterable
	clusterLabels := make(map[Clusterable]string)
	// Populate label map
	for _, clusterable := range e.clusterables {
		clusterLabels[clusterable] = UNKNOWN
	}

	c := 0 // Cluster counter
	for _, clusterable := range e.clusterables {
		// Continue if clusterable was previously processed in inner loop
		if clusterLabels[clusterable] != UNKNOWN {
			continue
		}
		// Find neighborhood of current clusterable
		neighbors := e.Neighborhood(clusterable)
		// If density of neighborhood is insufficient, label point as noise and continue
		if len(neighbors) < e.minPts {
			clusterLabels[clusterable] = NOISE
			continue
		}
		// increase cluster label
		c++
		// Label initial point
		clusterLabels[clusterable] = strconv.Itoa(c)

		// Process each clusterable in neighborhood
		seeds := neighbors
		for i := 0; i < len(seeds); i++ {
			// If same, continue
			if seeds[i].ID() == clusterable.ID() {
				continue
			}
			// Check if was previously considered, and is not noise
			switch clusterLabels[seeds[i]] {
			case NOISE: // If clusterable is marked as noise, add to current cluster
				clusterLabels[seeds[i]] = strconv.Itoa(c)
			case UNKNOWN: // If clusterable was previously unconsidered, proceed with extended neighborhood
				// Label current seed as part of the current cluster
				clusterLabels[seeds[i]] = strconv.Itoa(c)
				// Check density of the neighborhood of the current seed; if valid, add
				// all clusterables in the seed neighborhood to the seed set
				if extension := e.Neighborhood(seeds[i]); len(extension) >= e.minPts {
					for _, k := range extension {
						seeds = append(seeds, k)
					}
				}
			}
		}
	}
	return clusterLabels
}
