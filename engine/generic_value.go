package engine

import (
	"errors"
	"fmt"
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
	NullType
)

// IntegerValue - Implementation of ValueInterface that is containing integer values
type IntegerValue struct {
	Value int
}

// StringValue - Implementation of ValueInterface that is containing string values
type StringValue struct {
	Value string
}

// NullValue - Implementation of ValueInterface that is containing null
type NullValue struct {
}

// HandleValue - Function to take an instance of ValueInterface and cast to a specific implementation
func CastValueInterface(v ValueInterface) {
	switch value := v.(type) {
	case IntegerValue:
		fmt.Printf("IntegerValue with Value: %d\n", value.Value)
	case StringValue:
		fmt.Printf("StringValue with Value: %s\n", value.Value)
	case NullValue:
		fmt.Println("NullValue (no value)")
	default:
		fmt.Println("Unknown type")
	}
}

// ToString implementations
func (value IntegerValue) ToString() string { return strconv.Itoa(value.Value) }
func (value StringValue) ToString() string  { return value.Value }
func (value NullValue) ToString() string    { return "NULL" }

// GetType implementations
func (value IntegerValue) GetType() SupportedTypes { return IntType }
func (value StringValue) GetType() SupportedTypes  { return StringType }
func (value NullValue) GetType() SupportedTypes    { return NullType }

// IsEqual implementations
func (value IntegerValue) IsEqual(valueInterface ValueInterface) bool {
	return areEqual(value, valueInterface)
}
func (value StringValue) IsEqual(valueInterface ValueInterface) bool {
	return areEqual(value, valueInterface)
}
func (value NullValue) IsEqual(valueInterface ValueInterface) bool {
	return areEqual(value, valueInterface)
}

// isSmallerThan implementations
func (value IntegerValue) isSmallerThan(secondValue ValueInterface) bool {
	nullValue, isNull := secondValue.(NullValue)
	if isNull {
		return nullValue.isGreaterThan(value)
	}

	secondValueAsInteger, isInteger := secondValue.(IntegerValue)
	if !isInteger {
		log.Fatal("Can't compare Integer with other type")
	}

	return value.Value < secondValueAsInteger.Value
}

func (value StringValue) isSmallerThan(secondValue ValueInterface) bool {
	nullValue, isNull := secondValue.(NullValue)
	if isNull {
		return nullValue.isGreaterThan(value)
	}

	secondValueAsString, isString := secondValue.(StringValue)
	if !isString {
		log.Fatal("Can't compare String with other type")
	}

	return value.Value < secondValueAsString.Value
}

func (value NullValue) isSmallerThan(secondValue ValueInterface) bool {
	_, isNull := secondValue.(NullValue)

	if isNull {
		return false
	}

	return true
}

// isGreaterThan implementations
func (value IntegerValue) isGreaterThan(secondValue ValueInterface) bool {
	nullValue, isNull := secondValue.(NullValue)
	if isNull {
		return nullValue.isSmallerThan(value)
	}

	secondValueAsInteger, isInteger := secondValue.(IntegerValue)
	if !isInteger {
		log.Fatal("Can't compare Integer with other type")
	}

	return value.Value > secondValueAsInteger.Value
}
func (value StringValue) isGreaterThan(secondValue ValueInterface) bool {
	nullValue, isNull := secondValue.(NullValue)
	if isNull {
		return nullValue.isSmallerThan(value)
	}

	secondValueAsString, isString := secondValue.(StringValue)
	if !isString {
		log.Fatal("Can't compare String with other type")
	}

	return value.Value > secondValueAsString.Value
}

func (value NullValue) isGreaterThan(_ ValueInterface) bool {
	return false
}

func areEqual(first ValueInterface, second ValueInterface) bool {
	return first.GetType() == second.GetType() && first.ToString() == second.ToString()
}

func getMin(values []ValueInterface) (ValueInterface, error) {
	if len(values) == 0 {
		return nil, errors.New("can't extract min from empty array")
	}
	minValue := values[0]

	for _, value := range values[1:] {
		if value.isSmallerThan(minValue) {
			minValue = value
		}
	}
	return minValue, nil
}

func getMax(values []ValueInterface) (ValueInterface, error) {
	if len(values) == 0 {
		return nil, errors.New("can't extract max from empty array")
	}

	maxValue := values[0]
	for _, value := range values[1:] {
		if value.isGreaterThan(maxValue) {
			maxValue = value
		}
	}

	return maxValue, nil
}
