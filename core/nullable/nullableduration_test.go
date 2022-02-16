package nullable

import (
	"fmt"
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func TestNullableDurationGetPointerCopy(t *testing.T) {
	ni := NullableDuration{}
	ni.Set(3)
	niCopy := ni.GetPointerCopy()
	if *niCopy != ni.Value {
		t.Error("\tniCopy value should be the same as ni Value", ni, niCopy)
	}
	if &ni.Value == niCopy {
		t.Error("\tthe address of ni.Value and niCopy should be different", &ni.Value, &niCopy)
	}
	ni.Unset()
	niCopy = ni.GetPointerCopy()
	if niCopy != nil {
		t.Error("\tniCopy should be nil because ni HasValue is false", ni, niCopy)
	}
}

func TestNullableDurationSetUnset(t *testing.T) {
	ni := NullableDuration{}
	testValue := time.Duration(1)
	ni.Set(testValue)
	if ni.HasValue != true || ni.Value != testValue {
		t.Error("\tnullable struct in invalid state after Set call", ni)
	}
	ni.Unset()
	if ni.HasValue || ni.Value != defaultDurationValue {
		t.Error("\tnullable struct in invalid state after Unset call", ni)
	}
}

func TestNullableDurationScan(t *testing.T) {
	ns := NullableDuration{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("\tFailed to scan nil into NullableDuration", err, ns)
	}
	if ns.Value != 0 || ns.HasValue != false {
		t.Error("\tNullable int has wrong value after scanning nil", ns)
	}
	testValue := time.Duration(2)
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("\tFailed to scan value into NullableDuration", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testValue)
		t.Error(errMsg, ns)
	}
	testString := "abc"
	err = ns.Scan(testString)
	if err != nil && err.(errors.RichError).GetErrorCode() != coreerrors.ErrCodeWrongType {
		t.Error("\tExpected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable int has wrong value after scanning %v", testString)
		t.Error(errMsg, ns)
	}
}

func TestNullableDurationMarshalJson(t *testing.T) {
	ns := NullableDuration{
		Value:    0,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("\tdata from marshal was not null when underlaying nullable int was nil", data)
	}
	ns = NullableDuration{
		Value:    -2,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "-2" {
		t.Error("\tdata from marshal was not what was expected", data, ns)
	}
}

func TestNullableDurationUnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableDuration{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("\tUnmarshaling null should result in a nullable int with an empty value and is null being true", ns)
	}
	testString = "5"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != 5 {
		t.Error("\tUnmarshaling 1.2 should result in a nullable int with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("\texpected an error", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("\tUnmarshaling false should result in a nullable int with an empty value and is null being true", ns)
	}
}
