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
	"#3a": TestCaseExampleData{
		`{
			"id": 1,
			"name": "A green door",
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
		false,
	},
	"#3b": TestCaseExampleData{
		`{
			"price": 12.50
		}`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product",
			"description": "A product from Acme's catalog",
			"type": "object",
			"properties": {
				"price": {
					"type": "number",
					"minimum": 0,
					"exclusiveMinimum": true
				}
			}
		}`,
		true,
	},
	"#3c": TestCaseExampleData{
		`{
			"price": 12.50
		}`,
		`{
			"type": "object",
			"properties": {
				"price": {
					"type": "number",
					"minimum": 15,
					"exclusiveMinimum": true
				}
			}
		}`,
		false,
	},
	"#3d": TestCaseExampleData{
		`{
			"price": 12
		}`,
		`{
			"type": "object",
			"properties": {
				"price": {
					"type": "number",
					"minimum": 12,
					"exclusiveMinimum": true
				}
			}
		}`,
		false,
	},
	"#3e": TestCaseExampleData{
		`{
			"price": 12
		}`,
		`{
			"type": "object",
			"properties": {
				"price": {
					"type": "number",
					"minimum": 12,
					"exclusiveMinimum": false
				}
			}
		}`,
		true,
	},
	"#4": TestCaseExampleData{
		`[
			{
				"id": 2,
				"name": "An ice sculpture",
				"price": 12.50,
				"tags": ["cold", "ice"],
				"dimensions": {
					"length": 7.0,
					"width": 12.0,
					"height": 9.5
				},
				"warehouseLocation": {
					"latitude": -78.75,
					"longitude": 20.4
				}
			},
			{
				"id": 3,
				"name": "A blue mouse",
				"price": 25.50,
				"dimensions": {
					"length": 3.1,
					"width": 1.0,
					"height": 1.0
				},
				"warehouseLocation": {
					"latitude": 54.4,
					"longitude": -32.7
				}
			}
		]`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product set",
			"type": "array",
			"items": {
				"title": "Product",
				"type": "object",
				"properties": {
					"id": {
						"description": "The unique identifier for a product",
						"type": "number"
					},
					"name": {
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
					},
					"dimensions": {
						"type": "object",
						"properties": {
							"length": {"type": "number"},
							"width": {"type": "number"},
							"height": {"type": "number"}
						},
						"required": ["length", "width", "height"]
					},
					"warehouseLocation": {
						"description": "Coordinates of the warehouse with the product",
						"$ref": "http://json-schema.org/geo"
					}
				},
				"required": ["id", "name", "price"]
			}
		}`,
		true,
	},
	"#4a": TestCaseExampleData{
		`[
			{
				"id": 2,
				"name": "An ice sculpture",
				"price": 12.50,
				"tags": ["cold", "ice"],
				"dimensions": {
					"length": 7.0,
					"width": 12.0,
					"height": 9.5
				},
				"warehouseLocation": {
					"latitude": -78.75,
					"longitude": 20.4
				}
			},
			{
				"id": 3,
				"name": "A blue mouse",
				"price": 25.50,
				"dimensions": {
					"length": 3.1,
					"width": 1.0,
					"height": 1.0
				},
				"warehouseLocation": {
					"latitude": 54.4,
					"longitude": -32.7
				}
			}
		]`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product set",
			"type": "array",
			"items": {
				"title": "Product",
				"type": "object",
				"properties": {
					"id": {
						"description": "The unique identifier for a product",
						"type": "number"
					},
					"name": {
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
					},
					"dimensions": {
						"type": "object",
						"properties": {
							"length": {"type": "number"},
							"width": {"type": "number"},
							"height": {"type": "number"}
						},
						"required": ["length", "width", "height"]
					},
					"warehouseLocation": {
						"description": "Coordinates of the warehouse with the product",
						"$ref": "http://json-schema.org/geo"
					}
				},
				"required": ["id", "name", "price"]
			}
		}`,
		true,
	},
	"#4b": TestCaseExampleData{
		`[
			{
				"warehouseLocation": {
					"latitude": "hello",
					"longitude": -32.7
				}
			}
		]`,
		`{
			"$schema": "http://json-schema.org/draft-04/schema#",
			"title": "Product set",
			"type": "array",
			"items": {
				"title": "Product",
				"type": "object",
				"properties": {
					"warehouseLocation": {
						"description": "Coordinates of the warehouse with the product",
						"$ref": "http://json-schema.org/geo"
					}
				}
			}
		}`,
		false,
	},
	"#5": TestCaseExampleData{
		`{
			"storage": {
				"type": "disk",
				"device": "/dev/sda1"
			},
			"fstype": "btrfs",
			"readonly": true
		}`,
		`{
			"type": "object",
			"required": [ "storage" ],
			"properties": {
				"storage": {
					"$ref": "#/definitions/diskDevice"
				}
			},
			"definitions": {
				"diskDevice": {
					"type" : "object",
					"properties": {
						"type": { "enum": [ "disk" ] },
						"device": {
							"type": "string",
							"pattern": "^/dev/[^/]+(/[^/]+)*$"
						}
					},
					"required": [ "type", "device" ]
				}
			}
		}`,
		true,
	},
	"#5a": TestCaseExampleData{
		`{
			"storage": {
				"type": "disk",
				"device": 1
			},
			"fstype": "btrfs",
			"readonly": true
		}`,
		`{
			"type": "object",
			"required": [ "storage" ],
			"properties": {
				"storage": {
					"$ref": "#/definitions/diskDevice"
				}
			},
			"definitions": {
				"diskDevice": {
					"type" : "object",
					"properties": {
						"type": { "enum": [ "disk" ] },
						"device": {
							"type": "string",
							"pattern": "^/dev/[^/]+(/[^/]+)*$"
						}
					},
					"required": [ "type", "device" ]
				}
			}
		}`,
		false,
	},
	"#6": TestCaseExampleData{
		`{
			"storage": {
				"type": "disk",
				"device": 1
			},
			"fstype": "btrfs",
			"readonly": true
		}`,
		`{
			"type": "object",
			"required": [ "storage" ],
			"properties": {
				"storage": {
					"$ref": "#/definitions/diskDevice"
				}
			},
			"definitions": {
				"diskDevice": {
					"type" : "object",
					"properties": {
						"type": { "enum": [ "disk" ] },
						"device": {
							"type": "string",
							"pattern": "^/dev/[^/]+(/[^/]+)*$"
						}
					},
					"required": [ "type", "device" ]
				}
			}
		}`,
		false,
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
