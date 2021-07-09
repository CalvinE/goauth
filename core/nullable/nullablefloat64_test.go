package nullable

import (
	"fmt"
	"testing"

	goautherrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/errors/codes"
)

func TestNullableFloat64GetPointerCopy(t *testing.T) {
	nf := NullableFloat64{}
	nf.Set(1.23)
	nfCopy := nf.GetPointerCopy()
	if *nfCopy != nf.Value {
		t.Error("nfCopy value should be the same as nf Value", nf, nfCopy)
	}
	if &nf.Value == nfCopy {
		t.Error("the address of nf.Value and nfCopy should be different", &nf.Value, &nfCopy)
	}
	nf.Unset()
	nfCopy = nf.GetPointerCopy()
	if nfCopy != nil {
		t.Error("nfCopy should be nil because nf HasValue is false", nf, nfCopy)
	}
}

func TestNullableFloat64SetUnset(t *testing.T) {
	nf := NullableFloat64{}
	testValue := float64(1.23)
	nf.Set(testValue)
	if nf.HasValue != true || nf.Value != testValue {
		t.Error("nullable struct in invalid state after Set call", nf)
	}
	nf.Unset()
	if nf.HasValue || nf.Value != defaultFloat64Value {
		t.Error("nullable struct in invalid state after Unset call", nf)
	}
}

func TestNullableFloat64Scan(t *testing.T) {
	ns := NullableFloat64{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("Failed to scan nil into NullableFloat64", err, ns)
	}
	if ns.Value != 0 || ns.HasValue != false {
		t.Error("Nullable float64 has wrong value after scanning nil", ns)
	}
	testValue := 1.2
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("Failed to scan nil into NullableFloat64", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable float64 has wrong value after scanning %f", testValue)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)

	if err != nil && err.(goautherrors.RichError).ErrCode != codes.ErrCodeWrongType {
		t.Error("Expected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable float64 has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableFloat64MarshalJson(t *testing.T) {
	ns := NullableFloat64{
		Value:    0,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("data from marshal was not null when underlaying nullable float64 was nil", data)
	}
	ns = NullableFloat64{
		Value:    1.2,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("Failed to marshal null to JSON.", err)
	}
	if string(data) != "1.2" {
		t.Error("data from marshal was not what was expected", data, ns)
	}
}

func TestNullableFloat64UnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableFloat64{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("Unmarshaling null should result in a nullable float64 with an empty value and is null being true", ns)
	}
	testString = "1.2"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("Failed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != 1.2 {
		t.Error("Unmarshaling 1.2 should result in a nullable float64 with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("expected an error", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("Unmarshaling false should result in a nullable float64 with an empty value and is null being true", ns)
	}
}
