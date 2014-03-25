package json

type JsonType int64
const (
	JsonType_Bool = JsonType(iota)
	JsonType_Number
	JsonType_String
	JsonType_Array
	JsonType_Object
	JsonType_Nil
)

func (j JsonType)String() string {
	switch j {
	case JsonType_Bool:
		return "booleans"
	case JsonType_Number:
		return "numbers"
	case JsonType_String:
		return "strings"
	case JsonType_Array:
		return "arrays"
	case JsonType_Object:
		return "objects"
	case JsonType_Nil:
		return "null"
	default:
		return "unknown"
	}
}
