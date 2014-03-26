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

	schema = &schemaObject{}
	var load func(shm *schemaObject, jsonobj map[string]interface{}) (err error)
	load = func(shm *schemaObject, jsonobj map[string]interface{}) (err error) {
		shm.child = make(map[string]*schemaObject)
		shm.title, _ = jsonobj["title"].(string)
		shm.schema, _ = jsonobj["$schema"].(string)
		shm.description, _ = jsonobj["description"].(string)
		shm.required, _ = jsonobj["required"].([]interface{})

		jsontype, _ := jsonobj["type"].(string)
		shm.jsontype, err = GetJsonType(jsontype)
		if err != nil {
			return
		}

		switch shm.jsontype {
		case JsonType_Object:
			if props, ok := jsonobj["properties"].(map[string]interface{}); ok {
				for k, p := range props {
					if ip, ok := p.(map[string]interface{}); ok {
						nextsch := &schemaObject{}
						load(nextsch, ip)
						shm.child[k] = nextsch
					}
				}
			}
		case JsonType_Array:
			if item, ok := json["item"].(map[string]interface{}); ok {
				nextsch := &schemaObject{}
				load(nextsch, item)
				shm.child["item"] = nextsch
			}
		}

		return
	}

	err = load(schema, json)
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
		if !sc.jsontype.IsValidType(obj) {
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
