package jsonschema

import (
	"math"
	"reflect"
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
	jsontype      []JsonType
	childs        map[string]*schemaProperty
	patternChilds map[string]*schemaProperty
	items         []*schemaProperty
	isItemsOne    bool
	subprop_list  []schemaPropertySub

	isSetMaxLength bool
	maxLength      int

	isSetMinLength bool
	minLength      int

	isSetMinProperties bool
	minProperties      int

	isSetMaxItems bool
	maxItems      int

	isSetMinItems bool
	minItems      int

	additionalProperties      *schemaProperty
	allowAdditionalProperties bool

	additionalItems      *schemaProperty
	allowAdditionalItems bool

	allOf []*schemaProperty
	anyOf []*schemaProperty
	oneOf []*schemaProperty

	required []string

	not *schemaProperty

	enum []interface{}

	uniqueItems bool

	dependency       map[string][]string
	dependencySchema map[string]*schemaProperty

	pattern string

	multipleOf float64

	// validation
	checked []string
}

func newSchemaProperty(mother *schemaProperty, schema *schemaObject, original string) *schemaProperty {
	return &schemaProperty{
		jsontype:                  make([]JsonType, 0),
		childs:                    make(map[string]*schemaProperty),
		patternChilds:             make(map[string]*schemaProperty),
		items:                     make([]*schemaProperty, 0),
		mother:                    mother,
		schemaobject:              schema,
		original:                  original,
		allowAdditionalProperties: true,
		allowAdditionalItems:      true,
		allOf:                     make([]*schemaProperty, 0),
		anyOf:                     make([]*schemaProperty, 0),
		oneOf:                     make([]*schemaProperty, 0),
		required:                  make([]string, 0),
		enum:                      make([]interface{}, 0),
		dependency:                make(map[string][]string),
		dependencySchema:          make(map[string]*schemaProperty),
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
		s.SetMaxItems,
		s.SetMaxLength,
		s.SetMinItems,
		s.SetMinLength,
		s.SetMinProperties,
		s.SetPatternChilds,
		s.SetItems,
		s.SetAdditionalProperties,
		s.SetAdditionalItems,
		s.SetAllOf,
		s.SetAnyOf,
		s.SetOneOf,
		s.SetRequired,
		s.SetNot,
		s.SetEnum,
		s.SetUniqueItems,
		s.SetDependency,
		s.SetPattern,
		s.SetSubProperties,
		s.SetChilds,
		s.SetMultipleOf,
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
	creater_list := []func(map[string]interface{}) (schemaPropertySub, error){
		newSubProp_maxProperties,
		newSubProp_maximum,
		newSubProp_minimum,
	}

	for _, fn := range creater_list {
		obj, err := fn(schema)
		if err != nil {
			return err
		}

		if obj != nil {
			s.subprop_list = append(s.subprop_list, obj)
		}
	}

	return nil
}

func (s *schemaProperty) SetMultipleOf(schema map[string]interface{}) error {
	if v, ok := schema["multipleOf"]; ok {
		if val, ok := v.(float64); ok {
			s.multipleOf = val
		}
	}
	return nil
}

func (s *schemaProperty) SetPattern(schema map[string]interface{}) error {
	if v, ok := schema["pattern"]; ok {
		if val, ok := v.(string); ok {
			s.pattern = val
		}
	}
	return nil
}

func (s *schemaProperty) SetDependency(schema map[string]interface{}) error {
	if v, ok := schema["dependencies"]; ok {
		if depkey, ok := v.(map[string]interface{}); ok {
			for k, v := range depkey {
				dep := make([]string, 0)
				if val, ok := v.([]interface{}); ok {
					for _, v := range val {
						dep = append(dep, v.(string))
					}
					s.dependency[k] = dep
				} else if val, ok := v.(map[string]interface{}); ok {
					news := s.NewChild()
					err := news.Recognize(val)
					if err != nil {
						return err
					}
					s.dependencySchema[k] = news
				}
			}
		}
	}

	return nil
}

func (s *schemaProperty) SetUniqueItems(schema map[string]interface{}) error {
	if v, ok := schema["uniqueItems"]; ok {
		if yes, ok := v.(bool); ok {
			s.uniqueItems = yes
		}
	}
	return nil
}

func (s *schemaProperty) SetEnum(schema map[string]interface{}) error {
	if v, ok := schema["enum"]; ok {
		if obj, ok := v.([]interface{}); ok {
			s.enum = append(s.enum, obj...)
		}
	}

	return nil
}

func (s *schemaProperty) SetRef(schema map[string]interface{}) error {
	if v, ok := schema["$ref"]; ok {
		if path, ok := v.(string); ok {
			news := s.NewBrother()
			err := s.schemaobject.refResolver.GetReferencedObject(path, news)
			if err != nil {
				return err
			}
			*s = *news

			return errFoundReference
		}
	}
	return nil
}

func (s *schemaProperty) SetJsonTypes(schema map[string]interface{}) error {
	if v, ok := schema["type"]; ok {
		switch typename := v.(type) {
		case string:
			type_raw, err := GetJsonType(typename)
			if err != nil {
				return err
			}

			s.jsontype = append(s.jsontype, type_raw)
		case []interface{}:
			for _, obj := range typename {
				if str, ok := obj.(string); ok {
					type_raw, err := GetJsonType(str)
					if err != nil {
						return err
					}

					s.jsontype = append(s.jsontype, type_raw)
				}
			}
		default:
			return ErrInvalidSchemaFormat
		}
	} else {
		s.jsontype = append(s.jsontype, JsonType_Any)
	}

	return nil
}

func (s *schemaProperty) SetChilds(schema map[string]interface{}) error {
	if obj, ok := schema["properties"]; ok {
		if obj2, ok := obj.(map[string]interface{}); ok {
			for k, v := range obj2 {
				if obj3, ok := v.(map[string]interface{}); ok {
					news := s.NewChild()
					err := news.Recognize(obj3)
					if err != nil {
						return err
					}
					s.childs[k] = news
				}
			}
		}
	}

	return nil
}

func (s *schemaProperty) SetPatternChilds(schema map[string]interface{}) error {
	if obj, ok := schema["patternProperties"]; ok {
		if obj2, ok := obj.(map[string]interface{}); ok {
			for k, v := range obj2 {
				if obj3, ok := v.(map[string]interface{}); ok {
					news := s.NewChild()
					err := news.Recognize(obj3)
					if err != nil {
						return err
					}
					s.patternChilds[k] = news
				}
			}
		}
	}

	return nil
}

func (s *schemaProperty) SetItems(schema map[string]interface{}) error {
	if obj, ok := schema["items"]; ok {
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
				if obj4, ok := obj3.(map[string]interface{}); ok {
					news := s.NewChild()
					err := news.Recognize(obj4)
					if err != nil {
						return err
					}
					s.items = append(s.items, news)
				}
			}
		}
	}
	return nil
}

func (s *schemaProperty) SetMaxItems(schema map[string]interface{}) error {
	if obj, ok := schema["maxItems"]; ok {
		if num, ok := obj.(float64); ok {
			s.isSetMaxItems = true
			s.maxItems = int(num)
		}
	}

	return nil
}

func (s *schemaProperty) SetMaxLength(schema map[string]interface{}) error {
	if obj, ok := schema["maxLength"]; ok {
		if f, ok := obj.(float64); ok {
			s.isSetMaxLength = true
			s.maxLength = int(f)
		}
	}

	return nil
}

func (s *schemaProperty) SetMinItems(schema map[string]interface{}) error {
	if obj, ok := schema["minItems"]; ok {
		if num, ok := obj.(float64); ok {
			s.isSetMinItems = true
			s.minItems = int(num)
		}
	}

	return nil
}

func (s *schemaProperty) SetMinLength(schema map[string]interface{}) error {
	if obj, ok := schema["minLength"]; ok {
		if f, ok := obj.(float64); ok {
			s.isSetMinLength = true
			s.minLength = int(f)
		}
	}

	return nil
}

func (s *schemaProperty) SetMinProperties(schema map[string]interface{}) error {
	if obj, ok := schema["minProperties"]; ok {
		if f, ok := obj.(float64); ok {
			s.isSetMinProperties = true
			s.minProperties = int(f)
		}
	}

	return nil
}

func (s *schemaProperty) SetAdditionalProperties(schema map[string]interface{}) error {
	if obj, ok := schema["additionalProperties"]; ok {
		if prop, ok := obj.(map[string]interface{}); ok {
			news := s.NewChild()
			err := news.Recognize(prop)
			if err != nil {
				return err
			}
			s.additionalProperties = news
		} else if allow, ok := obj.(bool); ok {
			s.allowAdditionalProperties = allow
		}
	}
	return nil
}

func (s *schemaProperty) SetAdditionalItems(schema map[string]interface{}) error {
	if obj, ok := schema["additionalItems"]; ok {
		if prop, ok := obj.(map[string]interface{}); ok {
			news := s.NewChild()
			err := news.Recognize(prop)
			if err != nil {
				return err
			}
			s.additionalItems = news
		} else if allow, ok := obj.(bool); ok {
			s.allowAdditionalItems = allow
		}
	}
	return nil
}

func (s *schemaProperty) SetAllOf(schema map[string]interface{}) error {
	if obj, ok := schema["allOf"]; ok {
		if prop, ok := obj.([]interface{}); ok {
			for _, v := range prop {
				if obj2, ok := v.(map[string]interface{}); ok {
					news := s.NewBrother()
					err := news.Recognize(obj2)
					if err != nil {
						return err
					}
					s.allOf = append(s.allOf, news)
				}
			}
		}
	}
	return nil
}

func (s *schemaProperty) SetAnyOf(schema map[string]interface{}) error {
	if obj, ok := schema["anyOf"]; ok {
		if prop, ok := obj.([]interface{}); ok {
			for _, v := range prop {
				if obj2, ok := v.(map[string]interface{}); ok {
					news := s.NewBrother()
					err := news.Recognize(obj2)
					if err != nil {
						return err
					}
					s.anyOf = append(s.anyOf, news)
				}
			}
		}
	}
	return nil
}

func (s *schemaProperty) SetOneOf(schema map[string]interface{}) error {
	if obj, ok := schema["oneOf"]; ok {
		if prop, ok := obj.([]interface{}); ok {
			for _, v := range prop {
				if obj2, ok := v.(map[string]interface{}); ok {
					news := s.NewBrother()
					err := news.Recognize(obj2)
					if err != nil {
						return err
					}
					s.oneOf = append(s.oneOf, news)
				}
			}
		}
	}
	return nil
}

func (s *schemaProperty) SetRequired(schema map[string]interface{}) error {
	if obj, ok := schema["required"]; ok {
		if prop, ok := obj.([]interface{}); ok {
			for _, v := range prop {
				if obj2, ok := v.(string); ok {
					s.required = append(s.required, obj2)
				}
			}
		}
	}
	return nil
}

func (s *schemaProperty) SetNot(schema map[string]interface{}) error {
	if obj, ok := schema["not"]; ok {
		if prop, ok := obj.(map[string]interface{}); ok {
			news := s.NewBrother()
			err := news.Recognize(prop)
			if err != nil {
				return err
			}
			s.not = news
		}
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
		p.IsMaxItemsValid,
		p.IsMaxLengthValid,
		p.IsMinItemsValid,
		p.IsMinLengthValid,
		p.IsMinPropertiesValid,
		p.IsPatternChildsValid,
		p.IsItemsValid,
		p.IsAllOfValid,
		p.IsAnyOfValid,
		p.IsOneOfValid,
		p.IsRequiredValid,
		p.IsNotValid,
		p.IsEnumValid,
		p.IsUniqueItemsValid,
		p.IsDependencyValid,
		p.IsPatternValid,
		p.IsMultipleOfValid,
		p.IsSubPropertiesValid,

		// fixed order
		p.IsChildsValid,
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

func (p *schemaProperty) IsMaxItemsValid(src interface{}) bool {
	if !p.isSetMaxItems {
		return true
	}

	if v, ok := src.([]interface{}); ok {
		return len(v) <= p.maxItems
	}

	return true
}

func (p *schemaProperty) IsMaxLengthValid(src interface{}) bool {
	if !p.isSetMaxLength {
		return true
	}

	if v, ok := src.(string); ok {
		return len(v) <= p.maxLength
	}

	return true
}

//--
func (p *schemaProperty) IsMinItemsValid(src interface{}) bool {
	if !p.isSetMinItems {
		return true
	}

	if v, ok := src.([]interface{}); ok {
		return len(v) >= p.minItems
	}

	return true
}

func (p *schemaProperty) IsMinLengthValid(src interface{}) bool {
	if !p.isSetMinLength {
		return true
	}

	if v, ok := src.(string); ok {
		return len(v) >= p.minLength
	}

	return true
}

func (p *schemaProperty) IsMinPropertiesValid(src interface{}) bool {
	if !p.isSetMinProperties {
		return true
	}

	if v, ok := src.(map[string]interface{}); ok {
		return len(v) >= p.minProperties
	}

	return true
}

func (p *schemaProperty) IsChildsValid(src interface{}) bool {
	if obj, ok := src.(map[string]interface{}); ok {
		for k, v := range p.childs {
			if prop, ok := obj[k]; ok {
				res := v.IsValid(prop)
				p.checked = append(p.checked, k)
				if !res {
					return false
				}
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

func (p *schemaProperty) IsPatternChildsValid(src interface{}) bool {
	if obj, ok := src.(map[string]interface{}); ok {
		for pat, child := range p.patternChilds {
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

func (s *schemaProperty) IsAllOfValid(src interface{}) bool {
	for _, subs := range s.allOf {
		if !subs.IsValid(src) {
			return false
		}
	}
	return true
}

func (s *schemaProperty) IsAnyOfValid(src interface{}) bool {
	if len(s.anyOf) == 0 {
		return true
	}

	for _, subs := range s.anyOf {
		if subs.IsValid(src) {
			return true
		}
	}
	return false
}

func (s *schemaProperty) IsOneOfValid(src interface{}) bool {
	if len(s.oneOf) == 0 {
		return true
	}

	cond := false
	for _, subs := range s.oneOf {
		if !subs.IsValid(src) {
			if cond {
				return false
			}
			cond = true
		}
	}
	return cond
}

func (s *schemaProperty) IsRequiredValid(src interface{}) bool {
	if obj, ok := src.(map[string]interface{}); ok {
		for _, v := range s.required {
			if _, ok := obj[v]; !ok {
				return false
			}
		}
	}
	return true
}

func (s *schemaProperty) IsNotValid(src interface{}) bool {
	if s.not == nil {
		return true
	}

	return !s.not.IsValid(src)
}

func (s *schemaProperty) IsEnumValid(src interface{}) bool {
	if len(s.enum) == 0 {
		return true
	}

	for _, v := range s.enum {
		if reflect.DeepEqual(v, src) {
			return true
		}
	}
	return false
}

func (s *schemaProperty) IsUniqueItemsValid(src interface{}) bool {
	if !s.uniqueItems {
		return true
	}

	if v, ok := src.([]interface{}); ok {
		for i, iv := range v {
			for j, jv := range v {
				if i == j {
					continue
				}

				if reflect.DeepEqual(iv, jv) {
					return false
				}
			}
		}
	}

	return true
}

func (s *schemaProperty) IsDependencyValid(src interface{}) bool {
	if obj, ok := src.(map[string]interface{}); ok {
		if len(s.dependency) != 0 {
			for k, v := range s.dependency {
				if _, ok := obj[k]; ok {
					for _, v2 := range v {
						if _, ok := obj[v2]; !ok {
							return false
						}
					}
				}
			}
		}

		if len(s.dependencySchema) != 0 {
			for k, v := range s.dependencySchema {
				if _, ok := obj[k]; ok {
					if !v.IsValid(obj) {
						return false
					}
				}
			}
		}
	}

	return true
}

func (s *schemaProperty) IsPatternValid(src interface{}) bool {
	if s.pattern != "" {
		if val, ok := src.(string); ok {
			re, err := regexp.Compile(s.pattern)
			if err != nil {
				return false
			}

			return re.MatchString(val)
		}
	}
	return true
}

func (s *schemaProperty) IsMultipleOfValid(src interface{}) bool {
	if s.multipleOf != 0 {
		if val, ok := src.(float64); ok {
			return (math.Mod(val*10e10, s.multipleOf*10e10) == 0)
		}
	}
	return true
}
