package engine

import (
	"log"
	"strconv"
)

// ValueInterface - Represent all supported types of data that can be inserted into Table
type ValueInterface interface {
	ToString() string
	GetType() SupportedTypes
	IsEqual(valueInterface ValueInterface) bool
	isSmallerThan(valueInterface ValueInterface) bool
	isGreaterThan(valueInterface ValueInterface) bool
}

type SupportedTypes int

const (
	IntType = iota
	StringType
)

// IntegerValue - Implementation of ValueInterface that is containing integer values
type IntegerValue struct {
	Value int
}

// StringValue - Implementation of ValueInterface that is containing string values
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

// isSmallerThan implementations
func (value IntegerValue) isSmallerThan(secondValue ValueInterface) bool {
	secondValueAsInteger, isInteger := secondValue.(IntegerValue)

	if !isInteger {
		log.Fatal("Can't compare Integer with other type")
	}

	return value.Value < secondValueAsInteger.Value
}
func (value StringValue) isSmallerThan(secondValue ValueInterface) bool {
	secondValueAsString, isString := secondValue.(StringValue)

	if !isString {
		log.Fatal("Can't compare String with other type")
	}

	return value.Value < secondValueAsString.Value
}

// isGreaterThan implementations
func (value IntegerValue) isGreaterThan(secondValue ValueInterface) bool {
	secondValueAsInteger, isInteger := secondValue.(IntegerValue)

	if !isInteger {
		log.Fatal("Can't compare Integer with other type")
	}

	return value.Value > secondValueAsInteger.Value
}
func (value StringValue) isGreaterThan(secondValue ValueInterface) bool {
	secondValueAsString, isString := secondValue.(StringValue)

	if !isString {
		log.Fatal("Can't compare String with other type")
	}

	return value.Value > secondValueAsString.Value
}

func areEqual(first ValueInterface, second ValueInterface) bool {
	return first.GetType() == second.GetType() && first.ToString() == second.ToString()
}
