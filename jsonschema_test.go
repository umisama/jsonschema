package jsonschema

import (
	"testing"
)

type TestCaseJsonType struct {
	jsontype JsonType
	ret      string
}

var TestCasesJsonType_String = []TestCaseJsonType{
	TestCaseJsonType{JsonType_Bool, "boolean"},
	TestCaseJsonType{JsonType_Null, "null"},
	TestCaseJsonType{JsonType_Number, "number"},
	TestCaseJsonType{JsonType_Integer, "integer"},
	TestCaseJsonType{JsonType_Array, "array"},
	TestCaseJsonType{JsonType_Object, "object"},
	TestCaseJsonType{JsonType_String, "string"},
}

type TestCaseExampleData struct {
	input, schema string
	isvalid       bool
}

// testdatas from examples on http://json-schema.org/example1.html.
// ref) http://json-schema.org/example1.html
var TestCasesExampleData = map[string]TestCaseExampleData{
	"#1": TestCaseExampleData{
		`{
			"id": 1,
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object"
		}`,
		true,
	},
	"#1a": TestCaseExampleData{
		`{
			"id": 1,
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "array"
		}`,
		false,
	},
	"#2": TestCaseExampleData{
		`{
			"id": 1,
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object",
			"properties": {
				"id": {
					"description": "The unique identifier for a product",
					"type": "integer"
				}
			},
			"required": ["id"]
		}`,
		true,
	},
	"#2a": TestCaseExampleData{
		`{
			"id": 1.5,
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object",
			"properties": {
				"id": {
					"description": "The unique identifier for a product",
					"type": "integer"
				}
			},
			"required": ["id"]
		}`,
		false,
	},
	"#2b": TestCaseExampleData{
		`{
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object",
			"properties": {
				"id": {
					"description": "The unique identifier for a product",
					"type": "integer"
				}
			},
			"required": ["id"]
		}`,
		false,
	},
	"#3": TestCaseExampleData{
		`{
			"id": 1,
			"name": "A green door",
			"price": 12.50,
			"tags": ["home", "green"]
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object",
			"properties": {
				"id": {
					"description": "The unique identifier for a product",
					"type": "integer"
				},
				"name": {
					"description": "Name of the product",
					"type": "string"
				},
				"price": {
					"type": "number",
					"minimum": 0,
					"exclusiveMinimum": true
				},
				"tags": {
					"type": "array",
					"items": {
						"type": "string"
					},
					"minItems": 1,
					"uniqueItems": true
				}
			},
			"required": ["id", "name", "price"]
		}`,
		true,
	},
}

func Test_JsonType_String(t *testing.T) {
	for k, v := range TestCasesJsonType_String {
		if v.jsontype.String() != v.ret {
			t.Error("fail on", k)
			continue
		}
	}
	return
}

func Test_NewValidator(t *testing.T) {
	t.Skip()
	schema := []byte(`{}`)
	val, err := NewValidator(schema)
	println(val, err)
	return
}

func Test_ExampleDatas(t *testing.T) {
	for k, v := range TestCasesExampleData {
		validator, err := NewValidator([]byte(v.schema))
		if err != nil {
			t.Error("fail on", k, "with", err)
			continue
		}

		if validator.IsValid([]byte(v.input)) != v.isvalid {
			t.Error("fail on", k)
			continue
		}
	}
}
