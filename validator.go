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

	result := true
	var check func(*schemaObject, interface{})
	check = func(sc *schemaObject, obj interface{}) {
		// validation type
		if !sc.jsontype.IsMatched(obj) {
			result = false
			return
		}

		// check required
		for _, v := range sc.required {
			prop, _ := obj.(map[string]interface{})
			_, ok := prop[v.(string)]
			if !ok {
				result = false
				return
			}
		}

		// go next if array or object
		switch sc.jsontype {
		case JsonType_Object:
			prop, _ := obj.(map[string]interface{})
			for k, v := range prop {
				c, ok := sc.child[k]
				if ok {
					check(c, v)
				}
			}
		case JsonType_Array:
			array, _ := obj.([]interface{})
			for _, v := range array {
				c, ok := sc.child["item"]
				if ok {
					check(c, v)
				}
			}
		}
		return
	}
	check(v.schema, jsonobj)

	return result
}

func (v *Validator) isTypeValid() {
	return
}
