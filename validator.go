package jsonschema

import (
	"encoding/json"
)

type Validator struct {
	schema *schemaObject
}

func NewValidator(schema []byte) (*Validator, error) {
	jsonmap := make(map[string]interface{})
	err := json.Unmarshal(schema, &jsonmap)
	if err != nil {
		return nil, err
	}

	return newValidator(jsonmap)
}

func newValidator(schema map[string]interface{})(*Validator, error) {
	s, err := newSchemaObject(schema)
	if err != nil {
		return nil, err
	}

	return &Validator{
		schema: s,
	}, nil
}

func (v *Validator) IsValid(src []byte)(bool, error) {
	var obj interface{}
	err := json.Unmarshal(src, &obj)
	if err != nil {
		return false, err
	}

	return v.schema.IsValid(obj), nil
}

func (v *Validator) Unmarshal(src []byte, dst interface{}) error {
	return nil
}

func (v *Validator) Marshal(src interface{}) ([]byte, error) {
	return nil, nil
}

func IsValid() bool {
	return false
}
