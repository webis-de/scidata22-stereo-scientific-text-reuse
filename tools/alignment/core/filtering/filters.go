/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package filtering

import . "picapica.org/picapica4/alignment/core/types"

// Returns true for any element, and is intended for testing purposes
func IdentityFilter(_ *Element) bool {
	return true
}

// Filters out elements with scalar values
func NoScalarValueFilter(e *Element) bool {
	if (e.Span.End - e.Span.Begin) == 0 {
		return false
	}
	return true
}
