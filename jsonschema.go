package json

import (
	"encoding/json"
	"errors"
)

// JsonType reprecents json schema's primitive types.
// defined at
type JsonType string

const (
	JsonType_Bool    = JsonType("boolean")
	JsonType_Number  = JsonType("number")
	JsonType_Integer = JsonType("integer")
	JsonType_String  = JsonType("string")
	JsonType_Array   = JsonType("array")
	JsonType_Object  = JsonType("object")
	JsonType_Null    = JsonType("null")
	JsonType_INVALID = JsonType("system-Invalid")
)

var (
	ErrInvalidTypeName = errors.New("jsonschema: invalid type name")
)

func GetJsonType(typestr string) (t JsonType, err error) {
	types := []JsonType{JsonType_Bool, JsonType_Number, JsonType_Integer, JsonType_String, JsonType_Array, JsonType_Object, JsonType_Null}
	for _, v := range types {
		if v.String() == typestr {
			return v, nil
		}
	}

	return JsonType_INVALID, ErrInvalidTypeName
}

func (j JsonType) String() string {
	return string(j)
}

func (j JsonType) IsValidType(v interface{}) (ret bool) {
	switch j {
	case JsonType_Array:
		_, ret = v.([]interface{})
	case JsonType_Bool:
		_, ret = v.(bool)
	case JsonType_Integer:
		_, ret = v.(int)
	case JsonType_Null:
		// TODO
		ret = false
	case JsonType_Number:
		_, ret = v.(float64)
	case JsonType_Object:
		_, ret = v.(map[string]interface{})
	case JsonType_String:
		_, ret = v.(string)
	default:
		ret = false
	}

	return
}

type schemaRoot struct {
	schema, title, ref, description string
	jsontype                        JsonType
	required, isarray               bool
	child                           map[string]schemaProperties
}

type schemaProperties struct {
	val map[string]schemaProperties
}

type Validator struct {
	schema *schemaRoot
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

func (v *Validator) parseSchema(json map[string]interface{}) (schema *schemaRoot, err error) {
	schema = &schemaRoot{}

	schema.title, _ = json["title"].(string)
	schema.schema, _ = json["$schema"].(string)
	schema.description, _ = json["description"].(string)

	jsontype, _ := json["type"].(string)
	schema.jsontype, err = GetJsonType(jsontype)
	if err != nil {
		return
	}

	return
}

func (v *Validator) IsValid(jsonstr []byte) bool {
	var jsonobj interface{}
	err := json.Unmarshal(jsonstr, &jsonobj)
	if err != nil {
		println(err.Error())
		return false
	}

	result := true
	check := func(v *schemaRoot, obj interface{}) {
		if !v.jsontype.IsValidType(obj) {
			result = false
			return
		}

		// go next if array or object
		switch v.jsontype {
		case JsonType_Array:
			//TODO
		case JsonType_Object:
			//TODO
		}
	}
	check(v.schema, jsonobj)

	return result
}
