package jsonschema

type schemaObject struct {
	schema, title, ref, description string
	jsontype                        JsonType
	isarray                         bool
	required                        []interface{}
	child                           map[string]*schemaObject
}

func (s *schemaObject) isRequired(name string) bool {
	for _, v := range s.required {
		if v == name {
			return true
		}
	}

	return false
}
