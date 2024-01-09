package engine

import (
	"testing"
)

func TestEquals(t *testing.T) {

	oneInt := IntegerValue{
		Value: 1,
	}
	twoInt := IntegerValue{
		Value: 2,
	}
	oneString := StringValue{
		Value: "one",
	}
	twoString := StringValue{
		Value: "two",
	}

	shouldBeEqual(t, oneInt, oneInt)
	shouldBeEqual(t, oneString, oneString)
	shouldNotBeEqual(t, oneInt, twoInt)
	shouldNotBeEqual(t, oneString, twoString)
	shouldNotBeEqual(t, oneString, oneInt)
}

func shouldBeEqual(t *testing.T, valueOne ValueInterface, valueTwo ValueInterface) {
	const ErrorMsgShouldBeEqual = "%s(type: %v) is not equal %s(type: %v), but is should be"

	if !valueOne.IsEqual(valueTwo) {
		t.Errorf(ErrorMsgShouldBeEqual, valueOne.ToString(), valueOne.GetType(), valueTwo.ToString(), valueTwo.GetType())
	}
}
func shouldNotBeEqual(t *testing.T, valueOne ValueInterface, valueTwo ValueInterface) {
	const ErrorMsgShouldNotBeEqual = "%s(type: %v) is equal %s(type: %v), but is shouldn't be"

	if valueOne.IsEqual(valueTwo) {
		t.Errorf(ErrorMsgShouldNotBeEqual, valueOne.ToString(), valueOne.GetType(), valueTwo.ToString(), valueTwo.GetType())
	}
}
