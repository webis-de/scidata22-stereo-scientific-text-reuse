/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package errors

import "fmt"

// Error if a method or function is not defined for the supplied type
type TypeError struct {
	t interface{}
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Unexpected type %T", e.t)
}

func NewTypeError(t interface{}) TypeError {
	return TypeError{t: t}
}

// Error if a specified preset name is not present in a preset file
type PresetNotFoundError struct {
	presetName     string
	presetFilePath string
}

func (e PresetNotFoundError) Error() string {
	return fmt.Sprintf("Cannot find preset with name %s in preset file %s", e.presetName, e.presetFilePath)
}

func NewPresetNotFoundError(presetName, presetFilePath string) PresetNotFoundError {
	return PresetNotFoundError{
		presetName:     presetName,
		presetFilePath: presetFilePath,
	}
}
