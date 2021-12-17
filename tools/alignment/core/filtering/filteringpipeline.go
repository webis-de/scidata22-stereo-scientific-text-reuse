/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package filtering

import . "picapica.org/picapica4/alignment/core/types"

// Function Interface for Filters
// A filter has a single function that returns true or false for a given element
type Filter func(element *Element) bool

// A FilterPipeline is a set of filters which are successively applied to elements
type FilterPipeline []Filter

// Returns a new FilterPipeline with an initialized, but empty set of filters
func NewFilterPipeline() FilterPipeline { return make([]Filter, 0, 0) }

// Returns true iff all registered filters return true for any given element.
func (fp FilterPipeline) All(element *Element) bool {
	for _, filter := range fp {
		if !filter(element) {
			return false
		}
	}
	return true
}

// Applies all registered filters to the given elements and returns only those
// for which all filters are true.
func (fp FilterPipeline) Reduce(elements ...*Element) []*Element {
	if len(elements) == 0 {
		return elements
	} else {
		filtered := make([]*Element, 0, len(elements))
		for _, element := range elements {
			if fp.All(element) {
				filtered = append(filtered, element)
			}
		}
		return filtered
	}
}

// Adds a one or more new filters to the FilterPipeline
func (fp *FilterPipeline) Register(filters ...Filter) {
	for _, filter := range filters {
		*fp = append(*fp, filter)
	}
}
