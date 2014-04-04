package jsonschema

import (
	"reflect"
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

func newSubProp_minimum(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_maximum(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_minProperties(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_maxProperties(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_maxLength(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_minLength(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_maxItems(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_minItems(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

func newSubProp_pattern(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
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

// defined at 5.3.4
type schemaPropertySub_uniqueItem struct {
	value bool
}

func newSubProp_uniqueItem(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	prop_raw, exist := schema["uniqueItems"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_uniqueItem)
	prop_b, ok := prop_raw.(bool)
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop_b
	return s, nil
}

func (s *schemaPropertySub_uniqueItem) IsValid(src interface{}) bool {
	val, ok := src.([]interface{})
	if !ok {
		return true
	}

	if !s.value {
		// always success if the keyword has false.
		return true
	}

	for k1, v1 := range val {
		for k2, v2 := range val {
			if reflect.DeepEqual(v1, v2) && k1 != k2 {
				return false
			}
		}
	}

	return true
}

// defined at 5.4.3
type schemaPropertySub_required struct {
	value []string
}

func newSubProp_required(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	prop_raw, exist := schema["required"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_required)
	prop_a, ok := prop_raw.([]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	prop_s := convInterfaceArrayToStringArray(prop_a)
	if prop_s == nil {
		return nil, ErrInvalidSchemaFormat
	}

	// values are must unique.
	for k1, v1 := range prop_a {
		for k2, v2 := range prop_a {
			if v1 == v2 && k1 != k2 {
				return nil, ErrInvalidSchemaFormat
			}
		}
	}

	s.value = prop_s
	return s, nil
}

func (s *schemaPropertySub_required) IsValid(src interface{}) bool {
	val, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	for _, v := range s.value {
		_, ok := val[v]
		if !ok {
			// elements not found
			return false
		}
	}

	return true
}

// defined at 5.4.5
type schemaPropertySub_dependency struct {
	elementname map[string][]string
	validation  map[string]*schemaProperty
}

func newSubProp_dependency(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	dep, ok := schema["dependencies"]
	if !ok {
		return nil, nil
	}

	depobjs, ok := dep.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s := &schemaPropertySub_dependency{
		elementname: make(map[string][]string),
		validation:  make(map[string]*schemaProperty, 0),
	}
	for name, value := range depobjs {
		switch depobj := value.(type) {
		case []interface{}:
			val := convInterfaceArrayToStringArray(depobj)
			if val == nil {
				return nil, ErrInvalidSchemaFormat
			}
			s.elementname[name] = val

		case map[string]interface{}:
			news := m.NewChild()
			err := news.Recognize(depobj)
			if err != nil {
				return nil, ErrInvalidSchemaFormat
			}
			s.validation[name] = news

		default:
			return nil, ErrInvalidSchemaFormat
		}
	}

	return s, nil
}

func (s *schemaPropertySub_dependency) IsValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	// keyname
	for name, deps := range s.elementname {
		if _, ok := obj[name]; !ok {
			// specified element was not found
			continue
		}

		for _, dep := range deps {
			// is depenedant keys exist?
			if _, ok := obj[dep]; !ok {
				return false
			}
		}
	}

	// element schema
	for name, dep := range s.validation {
		if _, ok := obj[name]; !ok {
			continue
		}

		if !dep.IsValid(obj) {
			return false
		}
	}

	return true
}

// defined at 5.5.1
type schemaPropertySub_enum struct {
	value []interface{}
}

func newSubProp_enum(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	prop_raw, exist := schema["enum"]
	if !exist {
		return nil, nil
	}

	s := new(schemaPropertySub_enum)
	prop, ok := prop_raw.([]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s.value = prop
	return s, nil
}

func (s *schemaPropertySub_enum) IsValid(src interface{}) bool {
	for _, v := range s.value {
		if reflect.DeepEqual(v, src) {
			return true
		}
	}
	return false
}

// defined at 5.5.3
type schemaPropertySub_allOf struct {
	value []*schemaProperty
}

func newSubProp_allOf(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	props_raw, exist := schema["allOf"]
	if !exist {
		return nil, nil
	}

	props, ok := props_raw.([]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s := &schemaPropertySub_allOf{
		value: make([]*schemaProperty, 0),
	}

	for _, prop := range props {
		news := m.NewBrother()

		prop_map, ok := prop.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidSchemaFormat
		}

		err := news.Recognize(prop_map)
		if err != nil {
			return nil, err
		}
		s.value = append(s.value, news)
	}

	return s, nil
}

func (s *schemaPropertySub_allOf) IsValid(src interface{}) bool {
	for _, v := range s.value {
		if !v.IsValid(src) {
			return false
		}
	}
	return true
}

// defined at 5.5.4
type schemaPropertySub_anyOf struct {
	value []*schemaProperty
}

func newSubProp_anyOf(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	props_raw, exist := schema["anyOf"]
	if !exist {
		return nil, nil
	}

	props, ok := props_raw.([]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s := &schemaPropertySub_anyOf{
		value: make([]*schemaProperty, 0),
	}

	for _, prop := range props {
		news := m.NewBrother()

		prop_map, ok := prop.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidSchemaFormat
		}

		err := news.Recognize(prop_map)
		if err != nil {
			return nil, err
		}
		s.value = append(s.value, news)
	}

	return s, nil
}

func (s *schemaPropertySub_anyOf) IsValid(src interface{}) bool {
	for _, sub := range s.value {
		if sub.IsValid(src) {
			return true
		}
	}

	return false
}

// defined at 5.5.5
type schemaPropertySub_oneOf struct {
	value []*schemaProperty
}

func newSubProp_oneOf(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	props_raw, exist := schema["oneOf"]
	if !exist {
		return nil, nil
	}

	props, ok := props_raw.([]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s := &schemaPropertySub_oneOf{
		value: make([]*schemaProperty, 0),
	}

	for _, prop := range props {
		prop_map, ok := prop.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidSchemaFormat
		}

		news := m.NewBrother()
		err := news.Recognize(prop_map)
		if err != nil {
			return nil, err
		}
		s.value = append(s.value, news)
	}

	return s, nil
}

func (s *schemaPropertySub_oneOf) IsValid(src interface{}) bool {
	cond := false
	for _, v := range s.value {
		if !v.IsValid(src) {
			if cond {
				return false
			}
			cond = true
		}
	}

	return cond
}

// defined at 5.5.5
type schemaPropertySub_not struct {
	value *schemaProperty
}

func newSubProp_not(schema map[string]interface{}, m *schemaProperty) (schemaPropertySub, error) {
	prop_raw, exist := schema["not"]
	if !exist {
		return nil, nil
	}

	prop, ok := prop_raw.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidSchemaFormat
	}

	s := new(schemaPropertySub_not)
	news := m.NewBrother()
	err := news.Recognize(prop)
	if err != nil {
		return nil, err
	}

	s.value = news
	return s, nil
}

func (s *schemaPropertySub_not) IsValid(src interface{}) bool {
	return !s.value.IsValid(src)
}
