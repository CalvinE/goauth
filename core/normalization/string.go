package normalization

import (
	"errors"
	"reflect"
	"strings"
)

func NormalizeStringValue(value interface{}) (string, error) {
	switch svt := value.(type) {
	case *string:
		stringValue := reflect.ValueOf(svt).Elem().Interface()
		return NormalizeStringValue(stringValue)
	case string:
		return strings.ToUpper(svt), nil
	case []byte:
		sValue := string(svt)
		return strings.ToUpper(sValue), nil
	default:
		// TODO: specific error here?
		return "", errors.New("")
	}
}
