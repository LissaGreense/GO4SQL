package engine

import "strconv"

type ValueInterface interface {
	ToString() string
	GetType() SupportedTypes
	IsEqual(valueInterface ValueInterface) bool
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

// IsEqual implementations
func (value IntegerValue) IsEqual(valueInterface ValueInterface) bool {
	return areEqual(value, valueInterface)
}
func (value StringValue) IsEqual(valueInterface ValueInterface) bool {
	return areEqual(value, valueInterface)
}

func areEqual(first ValueInterface, second ValueInterface) bool {
	return first.GetType() == second.GetType() && first.ToString() == second.ToString()
}
