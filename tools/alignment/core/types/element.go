/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import "fmt"

type Element struct {
	Span  *Span       `json:"span"`
	Value interface{} `json:"stringvalue"`
}

func (m *Element) String() string {
	return fmt.Sprintf("Span: [%s], Value: %v", m.Span.String(), m.Value)
}
