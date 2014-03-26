package jsonschema

type subschema_minimum struct {
	enable           bool
	exclusiveMinimum bool
	value            float64
}

type schemaObject struct {
	schema, title, ref, description string
	jsontype                        JsonType
	required                        []string
	child                           map[string]*schemaObject
	minimum                         *subschema_minimum
}

func NewSchemaObject() *schemaObject {
	return &schemaObject{
		child:    make(map[string]*schemaObject),
		required: make([]string, 0),
		jsontype: SchemaType_Unknown,
		minimum:  &subschema_minimum{false, false, 0},
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
		s.setMinimum,
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
	req_raw, _ := obj["required"].([]interface{})
	for _, v := range req_raw {
		if req, ok := v.(string); ok {
			s.required = append(s.required, req)
		}
	}
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

func (s *schemaObject) setMinimum(obj map[string]interface{}) error {
	if v, ok := obj["minimum"].(float64); ok {
		s.minimum.enable = true
		s.minimum.value = v
	} else {
		return nil
	}

	if v, ok := obj["exclusiveMinimum"].(bool); ok {
		s.minimum.exclusiveMinimum = v
	}

	return nil
}
