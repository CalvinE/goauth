package normalization

import (
	"errors"
	"reflect"
)

var (
	boolStringValues = map[string]bool{
		"TRUE":  true,
		"T":     true,
		"YES":   true,
		"Y":     true,
		"1":     true,
		"FALSE": false,
		"F":     false,
		"NO":    false,
		"N":     false,
		"0":     false,
	}

	boolIntValues = map[int64]bool{
		1: true,
		0: false,
	}

	boolUintValues = map[uint64]bool{
		1: true,
		0: false,
	}

	boolFloatValues = map[float64]bool{
		1: true,
		0: false,
	}
)

// ReadBoolValue attempts to read a value provided as a boolean value.
// So T / F or 1 / 0 or TRUE / FALSE or true / false will all be parsed into the appropriate boolean value.
// Also int uint and float values will be converted where 1 is true and 0 is false
// All strings are run through strings.ToUpper first
func ReadBoolValue(v interface{}, defaultToFalse bool) (bool, error) {
	switch cv := v.(type) {
	case *string:
		if cv == nil {
			if defaultToFalse {
				return false, nil
			}
			// TODO: specific error here?
			return false, errors.New("")
		}
		return stringToBool(*cv, defaultToFalse)
	case string:
		return stringToBool(cv, defaultToFalse)
	case *int8, *int16, *int, *int32, *int64:
		if cv == nil {
			if defaultToFalse {
				return false, nil
			}
			// TODO: specific error here?
			return false, errors.New("")
		}
		value := reflect.ValueOf(cv).Elem().Interface()
		return intToBool(value, defaultToFalse)
	case int8, int16, int32, int, int64:
		return intToBool(cv, defaultToFalse)
	case *uint8, *uint16, *uint32, *uint, *uint64:
		if cv == nil {
			if defaultToFalse {
				return false, nil
			}
			// TODO: specific error here?
			return false, errors.New("")
		}
		value := reflect.ValueOf(cv).Elem().Interface()
		return uintToBool(value, defaultToFalse)
	case uint8, uint16, uint32, uint, uint64:
		return uintToBool(cv, defaultToFalse)
	case *float32, *float64:
		if cv == nil {
			if defaultToFalse {
				return false, nil
			}
			// TODO: specific error here?
			return false, errors.New("")
		}
		value := reflect.ValueOf(cv).Elem().Interface()
		return floatToBool(value, defaultToFalse)
	case float32, float64:
		return floatToBool(cv, defaultToFalse)
	case *bool:
		if cv == nil {
			if defaultToFalse {
				return false, nil
			}
			// TODO: specific error here?
			return false, errors.New("")
		}
		value := reflect.ValueOf(cv).Elem().Bool()
		return value, nil
	case bool:
		return cv, nil
	case nil:
		if defaultToFalse {
			return false, nil
		}
		// TODO: specific error here?
		return false, errors.New("")
	}
	// TODO: specific error here?
	return false, errors.New("")
}

func stringToBool(value string, defaultToFalse bool) (bool, error) {
	normalizedString, err := NormalizeStringValue(value)
	if err != nil {
		return false, err
	}
	boolValue, containsValue := boolStringValues[normalizedString]
	if containsValue {
		return boolValue, nil
	}
	if defaultToFalse {
		return false, nil
	}
	// TODO: specific error here?
	return false, errors.New("")
}

func intToBool(value interface{}, defaultToFalse bool) (bool, error) {
	intValue, err := NormalizeIntValue(value)
	if err != nil {
		return false, err
	}
	boolValue, containsValue := boolIntValues[intValue]
	if containsValue {
		return boolValue, nil
	}
	if defaultToFalse {
		return false, nil
	}
	// TODO: specific error here?
	return false, errors.New("")
}

func uintToBool(value interface{}, defaultToFalse bool) (bool, error) {
	uintValue, err := NormalizeUintValue(value)
	if err != nil {
		return false, err
	}
	boolValue, containsValue := boolUintValues[uintValue]
	if containsValue {
		return boolValue, nil
	}
	if defaultToFalse {
		return false, nil
	}
	// TODO: specific error here?
	return false, errors.New("")
}

func floatToBool(value interface{}, defaultToFalse bool) (bool, error) {
	floatValue, err := NormalizeFloatValue(value)
	if err != nil {
		return false, err
	}
	boolValue, containsValue := boolFloatValues[floatValue]
	if containsValue {
		return boolValue, nil
	}
	if defaultToFalse {
		return false, nil
	}
	// TODO: specific error here?
	return false, errors.New("")
}
