/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package types

import (
	"fmt"
)

type Seed struct {
	SourceSpan *Span `json:"source,omitempty"`
	TargetSpan *Span `json:"target,omitempty"`
}

func (m *Seed) String() string {
	return fmt.Sprintf("Span 1: [%s], Span 2: [%s]", m.SourceSpan.String(), m.TargetSpan.String())
}
