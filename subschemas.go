package jsonschema

import (
	"regexp"
)

type schemaPropertySub interface {
	IsValid(interface{}) bool
}

// defined at 5.1.3.(@Validation)
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

// defined at 5.1.2.(@Validation)
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

// defined at 5.4.2. (@Validation)
type schemaPropertySub_minProperties struct {
	value int
}

func newSubProp_minProperties(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, prop_exist := schema["minProperties"]
	if !prop_exist {
		return nil, nil
	}

	s := new(schemaPropertySub_minProperties)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_minProperties) IsValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	return len(obj) >= s.value
}

// defined at 5.4.1. (@Validation)
type schemaPropertySub_maxProperties struct {
	value int
}

func newSubProp_maxProperties(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, prop_exist := schema["maxProperties"]
	if !prop_exist {
		return nil, nil
	}

	s := new(schemaPropertySub_maxProperties)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_maxProperties) IsValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	return len(obj) <= s.value
}

// defined at 5.4.1. (@Validation)
type schemaPropertySub_maxLength struct {
	value int
}

func newSubProp_maxLength(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, exist := schema["maxLength"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_maxLength)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_maxLength) IsValid(src interface{}) bool {
	src_s, ok := src.(string)
	if !ok {
		return true
	}

	return len(src_s) <= s.value
}

// defined at 5.4.2. (@Validation)
type schemaPropertySub_minLength struct {
	value int
}

func newSubProp_minLength(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, exist := schema["minLength"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_minLength)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_minLength) IsValid(src interface{}) bool {
	src_s, ok := src.(string)
	if !ok {
		return true
	}

	return len(src_s) >= s.value
}

// defined at 5.3.2. (@Validation)
type schemaPropertySub_maxItems struct {
	value int
}

func newSubProp_maxItems(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, exist := schema["maxItems"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_maxItems)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_maxItems) IsValid(src interface{}) bool {
	src_a, ok := src.([]interface{})
	if !ok {
		return true
	}

	return len(src_a) <= s.value
}

// defined at 5.3.3. (@Validation)
type schemaPropertySub_minItems struct {
	value int
}

func newSubProp_minItems(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, exist := schema["minItems"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_minItems)
	prop_i, ok := getInteger(prop_raw)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_i
	return s, nil
}

func (s *schemaPropertySub_minItems) IsValid(src interface{}) bool {
	src_a, ok := src.([]interface{})
	if !ok {
		return true
	}

	return len(src_a) >= s.value
}

// defined at 5.2.3. (@Validation)
type schemaPropertySub_pattern struct {
	value *regexp.Regexp
}

func newSubProp_pattern(schema map[string]interface{}) (schemaPropertySub, error) {
	prop_raw, exist := schema["pattern"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_pattern)
	prop_s, ok := prop_raw.(string)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	// FIXME: re2 is not compatible with ECMA-262.
	exp, err := regexp.CompilePOSIX(prop_s)
	if err != nil {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = exp
	return s, nil
}

func (s *schemaPropertySub_pattern) IsValid(src interface{}) bool {
	val, ok := src.(string)
	if !ok {
		return true
	}

	return s.value.Match([]byte(val))
}
