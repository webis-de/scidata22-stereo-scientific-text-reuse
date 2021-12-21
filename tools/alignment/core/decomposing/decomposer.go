/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package decomposing

import (
	. "picapica.org/picapica4/alignment/core/types"
)

type Decomposer interface {
	Decompose(interface{}) ([]*Element, error)
}
