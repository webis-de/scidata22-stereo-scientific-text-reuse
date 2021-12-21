/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package decomposing

import (
	"fmt"
	"picapica.org/picapica4/alignment/core/errors"
	. "picapica.org/picapica4/alignment/core/types"
)

// Simple Tokenizer that assumes no annotations.
// Text is given as string with every word delimited with the string given as wordDelimiter,
// Sentences are delimited with the string given as sentenceDelimiter
type StringDecomposer struct {
	wordDelimiter        string
	lenWordDelimiter     int
	sentenceDelimiter    string
	lenSentenceDelimiter int
}

func (t StringDecomposer) Decompose(input interface{}) ([]*Element, error) {
	switch v := input.(type) {
	case string:
		var (
			seq            = make([]*Element, 0, 0)
			err     error  = nil
			word    string = ""
			index   uint32 = 0
			endLast int    = 0
		)

		for pos, char := range input.(string) {
			if fmt.Sprintf("%c", char) == t.sentenceDelimiter || fmt.Sprintf("%c", char) == t.wordDelimiter {
				if word != "" {
					seq = append(seq, &Element{
						Span: &Span{
							Begin: uint32(endLast),
							End:   uint32(pos),
							Index: index,
						},
						Value: word,
					})
					word = ""
					endLast = pos
					index++
				}
			} else {
				word = fmt.Sprintf("%s%c", word, char)
			}
		}
		return seq, err
	default:
		return make([]*Element, 0, 0), errors.NewTypeError(v)
	}
}

func NewStringDecomposer(wordDelimiter string, sentenceDelimiter string) *StringDecomposer {
	return &StringDecomposer{wordDelimiter: wordDelimiter, lenWordDelimiter: len(wordDelimiter), sentenceDelimiter: sentenceDelimiter, lenSentenceDelimiter: len(sentenceDelimiter)}
}
