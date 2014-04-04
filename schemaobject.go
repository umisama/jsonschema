package jsonschema

import (
	"regexp"
)

// schemaObject reprecents a jsonschema.
type schemaObject struct {
	recognized  *schemaProperty
	raw         map[string]interface{}
	refResolver *refResolver
}

func newSchemaObject(schema map[string]interface{}) (s *schemaObject, err error) {
	s = new(schemaObject)
	s.raw = schema

	s.refResolver, err = newRefResolver(schema)
	if err != nil {
		return
	}

	s.recognized, err = s.newSchemaProperty(schema)
	if err != nil {
		return
	}
	return
}

func (s *schemaObject) newSchemaProperty(schema map[string]interface{}) (prop *schemaProperty, err error) {
	prop = newSchemaProperty(nil, s, "#")

	err = prop.Recognize(schema)
	if err != nil {
		return
	}

	return
}

// schemaProperty reprecents a property of jsonschema.
type schemaProperty struct {
	mother       *schemaProperty
	schemaobject *schemaObject
	original     string

	// properties
	jsontype []JsonType

	properties                map[string]*schemaProperty
	patternProperties         map[string]*schemaProperty
	subprop_list              []schemaPropertySub
	additionalProperties      *schemaProperty
	allowAdditionalProperties bool

	isItemsOne           bool
	items                []*schemaProperty
	additionalItems      *schemaProperty
	allowAdditionalItems bool

	// validation
	checked []string
}

func newSchemaProperty(mother *schemaProperty, schema *schemaObject, original string) *schemaProperty {
	return &schemaProperty{
		jsontype:                  make([]JsonType, 0),
		properties:                make(map[string]*schemaProperty),
		patternProperties:         make(map[string]*schemaProperty),
		items:                     make([]*schemaProperty, 0),
		mother:                    mother,
		schemaobject:              schema,
		original:                  original,
		allowAdditionalProperties: true,
		allowAdditionalItems:      true,
		subprop_list:              make([]schemaPropertySub, 0),
	}
}

func (s *schemaProperty) NewChild() *schemaProperty {
	return newSchemaProperty(s, s.schemaobject, s.original)
}

func (s *schemaProperty) NewBrother() *schemaProperty {
	return newSchemaProperty(s.mother, s.schemaobject, s.original)
}

func (s *schemaProperty) Recognize(schema map[string]interface{}) error {
	fnlist := []func(map[string]interface{}) error{
		s.SetRef,
		s.SetJsonTypes,
		s.SetPatternProperties,
		s.SetItems,
		s.SetAdditionalProperties,
		s.SetAdditionalItems,
		s.SetSubProperties,
		s.SetProperties,
	}

	for _, fn := range fnlist {
		err := fn(schema)
		if err == errFoundReference {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (s *schemaProperty) SetSubProperties(schema map[string]interface{}) error {
	creater_list := []func(map[string]interface{}, *schemaProperty) (schemaPropertySub, error){
		newSubProp_maxProperties,
		newSubProp_minProperties,
		newSubProp_maximum,
		newSubProp_minimum,
		newSubProp_maxLength,
		newSubProp_minLength,
		newSubProp_maxItems,
		newSubProp_minItems,
		newSubProp_pattern,
		newSubProp_uniqueItem,
		newSubProp_required,
		newSubProp_dependency,
		newSubProp_enum,
		newSubProp_allOf,
		newSubProp_anyOf,
		newSubProp_oneOf,
		newSubProp_not,
		newSubProp_multipleOf,
	}

	for _, fn := range creater_list {
		obj, err := fn(schema, s)
		if err != nil {
			return err
		}

		if obj != nil {
			s.subprop_list = append(s.subprop_list, obj)
		}
	}

	return nil
}

func (s *schemaProperty) SetRef(schema map[string]interface{}) error {
	v, ok := schema["$ref"]
	if !ok {
		return nil
	}

	if path, ok := v.(string); ok {
		news := s.NewBrother()
		err := s.schemaobject.refResolver.GetReferencedObject(path, news)
		if err != nil {
			return err
		}
		*s = *news

		return errFoundReference
	}
	return nil
}

func (s *schemaProperty) SetJsonTypes(schema map[string]interface{}) error {
	v, ok := schema["type"]
	if !ok {
		s.jsontype = append(s.jsontype, JsonType_Any)
		return nil
	}

	switch typename := v.(type) {
	case string:
		type_raw, err := GetJsonType(typename)
		if err != nil {
			return err
		}

		s.jsontype = append(s.jsontype, type_raw)

	case []interface{}:
		for _, obj := range typename {
			str, ok := obj.(string)
			if !ok {
				return ErrInvalidSchemaFormat
			}

			type_raw, err := GetJsonType(str)
			if err != nil {
				return err
			}

			s.jsontype = append(s.jsontype, type_raw)
		}
	default:
		return ErrInvalidSchemaFormat
	}

	return nil
}

func (s *schemaProperty) SetProperties(schema map[string]interface{}) error {
	obj, ok := schema["properties"]
	if !ok {
		return nil
	}

	obj2, ok := obj.(map[string]interface{})
	if !ok {
		return ErrInvalidSchemaFormat
	}

	for k, v := range obj2 {
		obj3, ok := v.(map[string]interface{})
		if !ok {
			return ErrInvalidSchemaFormat
		}

		news := s.NewChild()
		err := news.Recognize(obj3)
		if err != nil {
			return err
		}

		s.properties[k] = news
	}

	return nil
}

func (s *schemaProperty) SetPatternProperties(schema map[string]interface{}) error {
	obj, ok := schema["patternProperties"]
	if !ok {
		return nil
	}

	obj2, ok := obj.(map[string]interface{})
	if !ok {
		return ErrInvalidSchemaFormat
	}

	for k, v := range obj2 {
		obj3, ok := v.(map[string]interface{})
		if !ok {
			return ErrInvalidSchemaFormat
		}

		news := s.NewChild()
		err := news.Recognize(obj3)
		if err != nil {
			return err
		}

		s.patternProperties[k] = news
	}

	return nil
}

func (s *schemaProperty) SetItems(schema map[string]interface{}) error {
	obj, ok := schema["items"]
	if !ok {
		return nil
	}
	if obj2, ok := obj.(map[string]interface{}); ok {
		news := s.NewChild()
		err := news.Recognize(obj2)
		if err != nil {
			return err
		}

		s.isItemsOne = true
		s.items = append(s.items, news)
	} else if obj2, ok := obj.([]interface{}); ok {
		for _, obj3 := range obj2 {
			obj4, ok := obj3.(map[string]interface{})
			if !ok {
				return ErrInvalidSchemaFormat
			}

			news := s.NewChild()
			err := news.Recognize(obj4)
			if err != nil {
				return err
			}

			s.items = append(s.items, news)
		}
	}

	return nil
}

func (s *schemaProperty) SetAdditionalProperties(schema map[string]interface{}) error {
	obj, ok := schema["additionalProperties"]
	if !ok {
		return nil
	}

	switch prop := obj.(type) {
	case map[string]interface{}:
		news := s.NewChild()
		err := news.Recognize(prop)
		if err != nil {
			return err
		}

		s.additionalProperties = news

	case bool:
		s.allowAdditionalProperties = prop
	}

	return nil
}

func (s *schemaProperty) SetAdditionalItems(schema map[string]interface{}) error {
	obj, ok := schema["additionalItems"]
	if !ok {
		return nil
	}

	switch prop := obj.(type) {
	case map[string]interface{}:
		news := s.NewChild()
		err := news.Recognize(prop)
		if err != nil {
			return err
		}

		s.additionalItems = news

	case bool:
		s.allowAdditionalItems = prop
	}

	return nil
}

// ==validators
func (s *schemaObject) IsValid(src interface{}) bool {
	return s.recognized.IsValid(src)
}

func (p *schemaProperty) IsValid(src interface{}) bool {
	p.checked = make([]string, 0)

	fnlist := []func(interface{}) bool{
		p.IsTypeValid,
		p.IsItemsValid,
		p.IsPatternPropertiesValid,
		p.IsSubPropertiesValid,

		// fixed order
		p.IsPropertiesValid,
		p.IsAdditionalPropertyValid,
		p.IsAdditionalItemsValid,
	}

	for _, fn := range fnlist {
		if !fn(src) {
			return false
		}
	}

	return true
}

func (p *schemaProperty) IsSubPropertiesValid(src interface{}) bool {
	for _, obj := range p.subprop_list {
		valid := obj.IsValid(src)
		if !valid {
			return false
		}
	}
	return true
}

func (p *schemaProperty) IsTypeValid(src interface{}) bool {
	for _, v := range p.jsontype {
		if v.IsMatched(src) {
			return true
		}
	}

	return false
}

//--
func (p *schemaProperty) IsPropertiesValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	for k, v := range p.properties {
		if prop, ok := obj[k]; ok {
			res := v.IsValid(prop)
			p.checked = append(p.checked, k)
			if !res {
				return false
			}
		}
	}

	return true
}

func (p *schemaProperty) IsItemsValid(src interface{}) bool {
	if len(p.items) == 0 {
		return true
	}

	if obj, ok := src.([]interface{}); ok {
		if !p.isItemsOne {
			for i := 0; i < len(obj); i++ {
				if len(p.items) <= i {
					break
				}
				if !p.items[i].IsValid(obj[i]) {
					return false
				}
			}
		} else {
			for _, v := range obj {
				if !p.items[0].IsValid(v) {
					return false
				}
			}
		}
	}

	return true
}

func (p *schemaProperty) IsPatternPropertiesValid(src interface{}) bool {
	obj, ok := src.(map[string]interface{})
	if !ok {
		return true
	}

	for pat, child := range p.patternProperties {
		re, err := regexp.Compile(pat)
		if err != nil {
			return false
		}

		for k, v := range obj {
			if re.MatchString(k) {
				res := child.IsValid(v)
				p.checked = append(p.checked, k)
				if !res {
					return false
				}
			}
		}
	}

	return true
}

func (p *schemaProperty) IsAdditionalPropertyValid(src interface{}) bool {
	if !p.allowAdditionalProperties {
		if obj, ok := src.(map[string]interface{}); ok {
			for k, _ := range obj {
				matched := false
				for _, checked := range p.checked {
					if checked == k {
						matched = true
					}
				}
				if !matched {
					return false
				}
			}
		}
	}

	if p.additionalProperties == nil {
		return true
	}

	if obj, ok := src.(map[string]interface{}); ok {
		for k, v := range obj {
			matched := false
			for _, checked := range p.checked {
				if checked == k {
					matched = true
				}
			}
			if !matched {
				res := p.additionalProperties.IsValid(v)
				if !res {
					return false
				}
			}
		}
	}

	return true
}

func (s *schemaProperty) IsAdditionalItemsValid(src interface{}) bool {
	if len(s.items) == 0 {
		return true
	}

	if s.isItemsOne {
		return true
	}

	if (s.additionalItems == nil) && (s.allowAdditionalItems) {
		return true
	}

	if obj, ok := src.([]interface{}); ok {
		if !s.allowAdditionalItems {
			if len(obj) > len(s.items) {
				return false
			}
		} else {
			for i := len(s.items) + 1; i < len(obj); i++ {
				if !s.additionalItems.IsValid(obj[i]) {
					return false
				}
			}
		}
	}

	return true
}
