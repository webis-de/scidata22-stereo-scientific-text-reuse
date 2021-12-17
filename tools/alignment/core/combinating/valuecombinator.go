/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package combinating

type ValueCombinator interface {
	Combine([]interface{}) (interface{}, error)
}
