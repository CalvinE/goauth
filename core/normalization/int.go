package normalization

import (
	"reflect"

	"github.com/calvine/goauth/core/errors"
)

func NormalizeIntValue(intValue interface{}) (int64, error) {
	switch ivt := intValue.(type) {
	case *int8, *int16, *int, *int32, *int64:
		intValue := reflect.ValueOf(ivt).Elem().Interface()
		return NormalizeIntValue(intValue)
	case int8:
		return int64(ivt), nil
	case int16:
		return int64(ivt), nil
	case int32:
		return int64(ivt), nil
	case int:
		return int64(ivt), nil
	case int64:
		return ivt, nil
	default:
		return 0, errors.NewInvalidTypeError(reflect.TypeOf(intValue).String(), true)
	}
}
