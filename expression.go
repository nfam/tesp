package goextract

import (
	"encoding/json"
	"errors"
)

//go:generate go run -tags=dev generate.go

// Plugin represents custom string process right after slice.between.
type Plugin func(input, name string) (string, error)

// Convert represents custom coverter to process final value.
type Convert func(input, name string) (interface{}, error)

// Expression represents a simex handler.
type Expression struct {
	json.Marshaler
	json.Unmarshaler
	value jsonSlice
}

// New returns the parsed `Expression` from given json data.
func New(data []byte) (*Expression, error) {
	e := &Expression{}
	err := e.UnmarshalJSON(data)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Must returns the parsed `Expression` from given json data, panic if fail.
func Must(data []byte) *Expression {
	e, err := New(data)
	if err != nil {
		panic(err)
	}
	return e
}

// UnmarshalJSON implements `json.Unmarshaler` interface.
func (e *Expression) UnmarshalJSON(data []byte) error {
	container := make(map[string]interface{})
	if err := json.Unmarshal(data, &container); err != nil {
		return err
	}
	if container == nil {
		return errors.New("null is invalid for simex Expression")
	}
	if err := e.value.fromJSON(container, "expression"); err != nil {
		return err
	}
	return nil
}

// MarshalJSON implements `json.Marshaler` interface.
func (e *Expression) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.value.toJSON())
}

// Extract returns result from extracting givent string.
func (e *Expression) Extract(input string, plugin Plugin, convert Convert) (interface{}, *Error) {
	return e.value.process(input, plugin, convert)
}
