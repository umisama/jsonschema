package jsonschema

type schemaObject struct {
	schema, title, ref, description string
	jsontype                        JsonType
	required                        []interface{}
	child                           map[string]*schemaObject
}

func NewSchemaObject() *schemaObject {
	return &schemaObject{
		child:    make(map[string]*schemaObject),
		required: make([]interface{}, 0),
	}
}

func (s *schemaObject) ParseJsonSchema(json map[string]interface{}) (err error) {
	loaders := []func(map[string]interface{}) error{
		s.setTitle,
		s.setSchema,
		s.setRef,
		s.setDescription,
		s.setJsonType,
		s.setRequired,
		s.setChilds,
	}

	for _, fn := range loaders {
		err = fn(json)
		if err != nil {
			return
		}
	}

	return
}

func (s *schemaObject) isRequired(name string) bool {
	for _, v := range s.required {
		if v == name {
			return true
		}
	}

	return false
}

func (s *schemaObject) setSchema(obj map[string]interface{}) error {
	s.schema, _ = obj["$schema"].(string)
	return nil
}

func (s *schemaObject) setTitle(obj map[string]interface{}) error {
	s.title, _ = obj["title"].(string)
	return nil
}

func (s *schemaObject) setRef(obj map[string]interface{}) error {
	s.ref, _ = obj["ref"].(string)
	return nil
}

func (s *schemaObject) setDescription(obj map[string]interface{}) error {
	s.description, _ = obj["description"].(string)
	return nil
}

func (s *schemaObject) setJsonType(obj map[string]interface{}) (err error) {
	typestr, ok := obj["type"].(string)
	if ok {
		s.jsontype, err = GetJsonType(typestr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *schemaObject) setRequired(obj map[string]interface{}) error {
	s.required, _ = obj["required"].([]interface{})
	return nil
}

func (s *schemaObject) setChilds(obj map[string]interface{}) error {
	switch s.jsontype {
	case JsonType_Object:
		if props, ok := obj["properties"].(map[string]interface{}); ok {
			for k, p := range props {
				if ip, ok := p.(map[string]interface{}); ok {
					news := NewSchemaObject()
					news.ParseJsonSchema(ip)
					s.child[k] = news
				}
			}
		}
	case JsonType_Array:
		if item, ok := obj["item"].(map[string]interface{}); ok {
			news := NewSchemaObject()
			news.ParseJsonSchema(item)
			s.child["item"] = news
		}
	}
	return nil
}
