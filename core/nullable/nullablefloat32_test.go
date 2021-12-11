package nullable

import (
	"fmt"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

func TestNullableFloat32GetPointerCopy(t *testing.T) {
	nf := NullableFloat32{}
	nf.Set(1.23)
	nfCopy := nf.GetPointerCopy()
	if *nfCopy != nf.Value {
		t.Error("\tnfCopy value should be the same as nf Value", nf, nfCopy)
	}
	if &nf.Value == nfCopy {
		t.Error("\tthe address of nf.Value and nfCopy should be different", &nf.Value, &nfCopy)
	}
	nf.Unset()
	nfCopy = nf.GetPointerCopy()
	if nfCopy != nil {
		t.Error("\tnfCopy should be nil because nf HasValue is false", nf, nfCopy)
	}
}

func TestNullableFloat32SetUnset(t *testing.T) {
	nf := NullableFloat32{}
	testValue := float32(1.23)
	nf.Set(testValue)
	if nf.HasValue != true || nf.Value != testValue {
		t.Error("\tnullable struct in invalid state after Set call", nf)
	}
	nf.Unset()
	if nf.HasValue || nf.Value != defaultFloat32Value {
		t.Error("\tnullable struct in invalid state after Unset call", nf)
	}
}

func TestNullableFloat32Scan(t *testing.T) {
	ns := NullableFloat32{}
	err := ns.Scan(nil)
	if err != nil {
		t.Error("\tFailed to scan nil into NullableFloat32", err, ns)
	}
	if ns.Value != 0 || ns.HasValue != false {
		t.Error("\tNullable float32 has wrong value after scanning nil", ns)
	}
	testValue := float32(1.2)
	err = ns.Scan(testValue)
	if err != nil {
		t.Error("\tFailed to scan nil into NullableFloat32", err, ns)
	}
	if ns.Value != testValue || ns.HasValue != true {
		errMsg := fmt.Sprintf("Nullable float32 has wrong value after scanning %f", testValue)
		t.Error(errMsg, ns)
	}
	testNumber := 3
	err = ns.Scan(testNumber)
	if err != nil && err.(errors.RichError).GetErrorCode() != coreerrors.ErrCodeWrongType {
		t.Error("\tExpected error to be of type WrongTypeError", err)
	}
	if ns.Value != 0 || ns.HasValue != false {
		errMsg := fmt.Sprintf("Nullable float32 has wrong value after scanning %d", testNumber)
		t.Error(errMsg, ns)
	}
}

func TestNullableFloat32MarshalJson(t *testing.T) {
	ns := NullableFloat32{
		Value:    0,
		HasValue: false,
	}
	data, err := ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "null" {
		t.Error("\tdata from marshal was not null when underlaying nullable float32 was nil", data)
	}
	ns = NullableFloat32{
		Value:    1.2,
		HasValue: true,
	}
	data, err = ns.MarshalJSON()
	if err != nil {
		t.Error("\tFailed to marshal null to JSON.", err)
	}
	if string(data) != "1.2" {
		t.Error("\tdata from marshal was not what was expected", data, ns)
	}
}

func TestNullableFloat32UnmarshalJson(t *testing.T) {
	testString := "null"
	ns := NullableFloat32{}
	err := ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal null", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("\tUnmarshaling null should result in a nullable float32 with an empty value and is null being true", ns)
	}
	testString = "1.2"
	err = ns.UnmarshalJSON([]byte(testString))
	if err != nil {
		t.Error("\tFailed to unmarshal \"Test\"", err)
	}
	if ns.HasValue != true || ns.Value != 1.2 {
		t.Error("\tUnmarshaling 1.2 should result in a nullable float32 with a value of 1.2 and is null being false", ns)
	}
	testString = "false"
	err = ns.UnmarshalJSON([]byte(testString))
	if err == nil {
		t.Error("\texpected an error", err)
	}
	if ns.HasValue != false || ns.Value != 0 {
		t.Error("\tUnmarshaling false should result in a nullable float32 with an empty value and is null being true", ns)
	}
}
