package engine

import "strconv"

type ValueInterface interface {
	ToString() string
	GetType() SupportedTypes
}

type SupportedTypes int

const (
	IntType = iota
	StringType
)

type IntegerValue struct {
	Value int
}

type StringValue struct {
	Value string
}

// ToString implementations
func (value IntegerValue) ToString() string { return strconv.Itoa(value.Value) }
func (value StringValue) ToString() string  { return value.Value }

// GetType implementations
func (value IntegerValue) GetType() SupportedTypes { return IntType }
func (value StringValue) GetType() SupportedTypes  { return StringType }
