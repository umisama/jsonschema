package jsonschema

import (
	"encoding/json"
)

type Validator struct {
	schema *schemaObject
}

func NewValidator(schema []byte) (v *Validator, err error) {
	v = new(Validator)

	schema_raw := make(map[string]interface{})
	err = json.Unmarshal(schema, &schema_raw)
	if err != nil {
		return
	}

	v.schema, err = v.parseSchema(schema_raw)

	return
}

func (v *Validator) parseSchema(json map[string]interface{}) (schema *schemaObject, err error) {
	schema = NewSchemaObject()
	err = schema.ParseJsonSchema(json)
	if err != nil {
		return
	}

	return
}

func (v *Validator) IsValid(jsonstr []byte) bool {
	var jsonobj interface{}
	err := json.Unmarshal(jsonstr, &jsonobj)
	if err != nil {
		return false
	}

	return v.isValid(v.schema, jsonobj)
}

func (v *Validator) isValid(sc *schemaObject, obj interface{}) bool {
	validators := []func(*schemaObject, interface{}) bool{
		v.isTypeValid,
		v.isRequiredValid,
		v.isMinimumValueValid,
		v.isChildsValid,
	}

	for _, fn := range validators {
		if !fn(sc, obj) {
			return false
		}
	}

	return true
}

func (v *Validator) isTypeValid(sc *schemaObject, obj interface{}) bool {
	return sc.jsontype.IsMatched(obj)
}

func (v *Validator) isRequiredValid(sc *schemaObject, obj interface{}) bool {
	switch sc.jsontype {
	case JsonType_Object:
		prop, _ := obj.(map[string]interface{})
		for _, v := range sc.required {
			_, ok := prop[v]
			if !ok {
				return false
			}
		}
	}

	return true
}

func (v *Validator) isChildsValid(sc *schemaObject, obj interface{}) bool {
	switch sc.jsontype {
	case JsonType_Object:
		prop, _ := obj.(map[string]interface{})
		for k, item := range prop {
			c, ok := sc.child[k]
			if ok {
				if !v.isValid(c, item) {
					return false
				}
			}
		}
	case JsonType_Array:
		array, _ := obj.([]interface{})
		for _, item := range array {
			c, ok := sc.child["item"]
			if ok {
				if !v.isValid(c, item) {
					return false
				}
			}
		}
	}
	return true
}

func (v *Validator) isMinimumValueValid(sc *schemaObject, obj interface{}) bool {
	if sc.minimum.enable {
		if val, ok := obj.(float64); ok {
			if sc.minimum.exclusiveMinimum {
				return (sc.minimum.value < val)
			}else {
				return (sc.minimum.value <= val)
			}
		} else {
			return false
		}
	}
	return true
}
