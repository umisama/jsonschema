package jsonschema

import (
	"math"
)

type schemaPropertySub interface {
	IsValid(interface{}) bool
}

// defined 5.1.3.(@Validation)
type schemaPropertySub_minimum struct {
	minimum          float64
	exclusiveMinimum bool
}

func newSubProp_minimum(schema map[string]interface{}) (schemaPropertySub, error) {
	min_raw, min_exist := schema["minimum"]
	excMin_raw, excMin_exist := schema["exclusiveMinimum"]

	if !min_exist {
		if excMin_exist {
			// If "exclusiveMaximum" is present, "maximum" MUST also be present.
			return nil, ErrInvalidSchemaFormat
		} else {
			return nil, nil
		}
	}

	s := new(schemaPropertySub_minimum)

	ok := false
	s.minimum, ok = min_raw.(float64)
	if !ok {
		// must JSON number.
		return nil, ErrInvalidSchemaFormat
	}

	if excMin_exist {
		s.exclusiveMinimum, ok = excMin_raw.(bool)
		if !ok {
			// must boolean.
			return nil, ErrInvalidSchemaFormat
		}
	} else {
		s.exclusiveMinimum = false
	}

	return s, nil
}

func (s *schemaPropertySub_minimum) IsValid(src interface{}) bool {
	val, ok := src.(float64)
	if !ok {
		return true
	}

	switch s.exclusiveMinimum {
	case true:
		return val > s.minimum
	case false:
		return val >= s.minimum
	}

	return false
}

// defined 5.1.2.(@Validation)
type schemaPropertySub_maximum struct {
	maximum          float64
	exclusiveMaximum bool
}

func newSubProp_maximum(schema map[string]interface{}) (schemaPropertySub, error) {
	max_raw, max_exist := schema["maximum"]
	excMax_raw, excMax_exist := schema["exclusiveMaximum"]

	if !max_exist {
		if excMax_exist {
			// If "exclusiveMaximum" is present, "maximum" MUST also be present.
			return nil, ErrInvalidSchemaFormat
		} else {
			return nil, nil
		}
	}

	s := new(schemaPropertySub_maximum)

	ok := false
	s.maximum, ok = max_raw.(float64)
	if !ok {
		// must JSON number
		return nil, ErrInvalidSchemaFormat
	}

	if excMax_exist {
		s.exclusiveMaximum, ok = excMax_raw.(bool)
		if !ok {
			// must boolean
			return nil, ErrInvalidSchemaFormat
		}
	} else {
		s.exclusiveMaximum = false
	}

	return s, nil
}

func (s *schemaPropertySub_maximum) IsValid(src interface{}) bool {
	val, ok := src.(float64)
	if !ok {
		return true
	}

	switch s.exclusiveMaximum {
	case true:
		return val < s.maximum
	case false:
		return val <= s.maximum
	}

	return false
}

// defined 5.4.1. (@Validation)
type schemaPropertySub_maxProperties struct {
	value int
}

func newSubProp_maxProperties(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, prop_exist := schema["maxProperties"]
	if !prop_exist {
		return nil, nil
	}

	s := new(schemaPropertySub_maxProperties)
	val_float, ok := prop_raw.(float64)
	if !ok {
		// must number
		return nil, ErrInvalidSchemaFormat
	}

	if math.Mod(val_float, 1) != 0 {
		// must number
		return nil, ErrInvalidSchemaFormat
	}

	s.value = int(val_float)
	return s, nil
}

func (s *schemaPropertySub_maxProperties) IsValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	return len(obj) <= s.value
}
