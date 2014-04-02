package jsonschema

import (
	"encoding/json"
	"errors"
	"github.com/umisama/jsonptr"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

var (
	ErrInvalidTypeName      = errors.New("jsonschema: invalid type name")
	ErrInvalidSchemaVersion = errors.New("jsonschema: invalid schema version")
	ErrInvalidSchemaFormat  = errors.New("jsonschema: invalid schema format")
	errFoundReference       = errors.New("notify found reference")
)

type SchemaType string

const (
	SchemaType_Draft3   = "http://json-schema.org/draft-03/schema#"
	SchemaType_Draft4   = "http://json-schema.org/draft-04/schema#"
	SchemaType_Standard = "http://json-schema.org/schema#"
	SchemaType_Unknown  = "unknown"
)

func GetSchemaType(typestr string) (t SchemaType) {
	types := []SchemaType{
		SchemaType_Draft3,
		SchemaType_Draft4,
		SchemaType_Standard,
	}

	for _, v := range types {
		if typestr == v.String() {
			return v
		}
	}
	return SchemaType_Unknown
}

func (s SchemaType) String() string {
	return string(s)
}

// JsonType reprecents json schema's primitive types.
type JsonType string

const (
	JsonType_Any     = JsonType("")
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
	types := []JsonType{
		JsonType_Bool,
		JsonType_Number,
		JsonType_Integer,
		JsonType_String,
		JsonType_Array,
		JsonType_Object,
		JsonType_Null,
	}

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

func (j JsonType) IsMatched(v interface{}) (ret bool) {
	switch j {
	case JsonType_Any:
		ret = true
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
		ret = (v == nil)
	case JsonType_Object:
		_, ret = v.(map[string]interface{})
	case JsonType_String:
		_, ret = v.(string)
	default:
		ret = false
	}

	return
}

type refResolver struct {
	originals         map[string]map[string]interface{}
	cached            map[string]*schemaProperty
	outherfile_schema map[string]*schemaProperty
}

func newRefResolver(schema map[string]interface{}) (*refResolver, error) {
	return &refResolver{
		originals:         map[string]map[string]interface{}{"#": schema},
		cached:            make(map[string]*schemaProperty),
		outherfile_schema: make(map[string]*schemaProperty),
	}, nil
}

func (r *refResolver) GetReferencedObject(path string, dst *schemaProperty) error {
	if obj, ok := r.cached[path]; ok {
		*dst = *obj
		return nil
	}

	raw := make(map[string]interface{})
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		raw = r.getReferenceObjectViaHttp(path)
		dst.original = path
	} else {
		raw = r.getReferecneObjectViaJsonPtr(path, dst.original)
	}

	r.cached[path] = dst
	err := dst.Recognize(raw)
	if err != nil {
		return err
	}
	return nil
}

func (r *refResolver) getReferecneObjectViaJsonPtr(path string, original string) map[string]interface{} {
	buf, _ := json.Marshal(r.originals[original])
	ret_buf, _ := jsonptr.Find(buf, path)

	ret := make(map[string]interface{})
	json.Unmarshal(ret_buf, &ret)
	return ret
}

func (r *refResolver) getReferenceObjectViaHttp(path string) map[string]interface{} {
	resp, err := http.Get(path)
	if err != nil {
		return nil
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	ret := make(map[string]interface{})
	json.Unmarshal(buf, &ret)
	r.originals[path] = ret
	return ret
}
