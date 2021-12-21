/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import (
	"fmt"
)

type Span struct {
	Begin uint32 `json:"begin"`
	End   uint32 `json:"end"`
	Index uint32 `json:"index"`
}

func (m *Span) String() string {
	return fmt.Sprintf("Begin: %d, End: %d, Index: %d", m.Begin, m.End, m.Index)
}

func (m *Span) Length() uint32 {
	return m.End - m.Begin
}
