package jsonschema

import (
	"math"
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
	case JsonType_Number:
		_, ret = v.(float64)
	case JsonType_Integer:
		// integer only(check to without floating point)
		if smpl, ok := v.(float64); ok {
			ret = math.Mod(smpl, 1) == 0
		}
	case JsonType_Null:
		ret = false
	case JsonType_Object:
		_, ret = v.(map[string]interface{})
	case JsonType_String:
		_, ret = v.(string)
	default:
		ret = false
	}

	return
}
