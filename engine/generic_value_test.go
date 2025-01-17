package engine

import (
	"testing"
)

func TestIsGreaterThan(t *testing.T) {
	oneInt := IntegerValue{
		Value: 0,
	}
	twoInt := IntegerValue{
		Value: 2,
	}
	oneString := StringValue{
		Value: "aaa",
	}
	twoString := StringValue{
		Value: "aab",
	}
	oneNull := NullValue{}
	twoNull := NullValue{}

	if oneInt.isGreaterThan(twoInt) {
		t.Errorf("0 shouldn't be greater than 2")
	}

	if !twoInt.isGreaterThan(oneInt) {
		t.Errorf("0 shouldn't be greater than 2")
	}

	if oneString.isGreaterThan(twoString) {
		t.Errorf("aaa shouldn't be greater than aab")
	}

	if !twoString.isGreaterThan(oneString) {
		t.Errorf("1 shouldn't be greater than 2")
	}

	if twoNull.isGreaterThan(oneNull) {
		t.Errorf("null is not greater than null")
	}

	if !oneInt.isGreaterThan(oneNull) {
		t.Errorf("Any Int value cannot be smaller than null")
	}

	if !oneString.isGreaterThan(oneNull) {
		t.Errorf("Any String value cannot be smaller than null")
	}

	if oneNull.isGreaterThan(oneInt) {
		t.Errorf("Null cannot be greater than any int value")
	}

	if oneNull.isGreaterThan(oneString) {
		t.Errorf("Null cannot be greater than any string value")
	}
}

func TestIsSmallerThan(t *testing.T) {
	oneInt := IntegerValue{
		Value: 0,
	}
	twoInt := IntegerValue{
		Value: 2,
	}
	oneString := StringValue{
		Value: "aaa",
	}
	twoString := StringValue{
		Value: "aab",
	}
	oneNull := NullValue{}
	twoNull := NullValue{}

	if !oneInt.isSmallerThan(twoInt) {
		t.Errorf("0 should be smaller than 2")
	}

	if twoInt.isSmallerThan(oneInt) {
		t.Errorf("0 should be smaller than 2")
	}

	if !oneString.isSmallerThan(twoString) {
		t.Errorf("aaa should be smaller than aab")
	}

	if twoString.isSmallerThan(oneString) {
		t.Errorf("1 should be smaller than 2")
	}

	if twoNull.isSmallerThan(oneNull) {
		t.Errorf("null is not smaller than null")
	}

	if oneInt.isSmallerThan(oneNull) {
		t.Errorf("Any int value cannot be smaller than null")
	}

	if oneString.isSmallerThan(oneNull) {
		t.Errorf("Any string value cannot be smaller than null")
	}

	if !oneNull.isSmallerThan(oneInt) {
		t.Errorf("Null cannot be greater than any int value")
	}

	if !oneNull.isSmallerThan(oneString) {
		t.Errorf("Null cannot be greater than any string value")
	}
}

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
	oneNull := NullValue{}
	twoNull := NullValue{}

	shouldBeEqual(t, oneInt, oneInt)
	shouldBeEqual(t, oneString, oneString)
	shouldBeEqual(t, oneNull, twoNull)
	shouldNotBeEqual(t, oneInt, twoInt)
	shouldNotBeEqual(t, oneString, twoString)
	shouldNotBeEqual(t, oneString, oneInt)
	shouldNotBeEqual(t, oneNull, oneInt)
	shouldNotBeEqual(t, oneNull, oneString)
	shouldNotBeEqual(t, twoInt, twoNull)
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
