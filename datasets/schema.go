package datasets

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// JSONSchema was extracted from heedy to allow

type jsonSchema_prop struct {
	s        *gojsonschema.Schema
	required bool
}
type JSONSchema struct {
	Schema map[string]interface{}

	s *gojsonschema.Schema

	props           map[string]jsonSchema_prop
	additionalProps bool
}

func NewSchema(schema map[string]interface{}) (*JSONSchema, error) {
	objectMap := make(map[string]interface{})
	objectMap["type"] = "object"
	objectMap["additionalProperties"] = false

	if v, ok := schema["type"]; ok {
		if v != "object" {
			return nil, errors.New("Schema must have type 'object'")
		}
		objectMap = schema
	} else {
		// Treat the schema as a prop map
		propMap := make(map[string]interface{})
		for k, v := range schema {
			switch k {
			// Allow these modifiers to go directly to the underlying schema object
			case "additionalProperties", "required":
				objectMap[k] = v
			default:
				propMap[k] = v
			}
		}
		objectMap["properties"] = propMap
	}

	s, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(objectMap))

	// Now load the per-prop schema, so that each prop can be validated individually
	// for update queries
	props := make(map[string]jsonSchema_prop)
	additionalProps := false
	if err == nil {
		apropi, ok := objectMap["additionalProperties"]
		if ok {
			additionalProps, ok = apropi.(bool)
			if !ok {
				return nil, errors.New("schema: additionalProperties must be a boolean")
			}
		}

		// Load the individual properties
		ppsi, ok := objectMap["properties"]
		if ok {
			pps, ok := ppsi.(map[string]interface{})
			if !ok {
				return nil, errors.New("schema: properties should be an object")
			}
			for k, v := range pps {
				sp, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(v))
				if err != nil {
					return nil, err
				}
				props[k] = jsonSchema_prop{
					s:        sp,
					required: false,
				}
			}
		}

		// Now mark required properties as required
		rpi, ok := objectMap["required"]
		if ok {
			rp, ok := rpi.([]interface{})
			if !ok {
				return nil, errors.New("schema: required props must be an array")
			}
			for _, v := range rp {
				vs, ok := v.(string)
				if !ok {
					return nil, errors.New("schema: elements of required props array msut be strings")
				}
				jsp, ok := props[vs]
				if !ok {
					return nil, errors.New("A required prop has not associated schema")
				}
				jsp.required = true
				props[vs] = jsp
			}
		}

	}

	return &JSONSchema{
		Schema:          objectMap,
		s:               s,
		props:           props,
		additionalProps: additionalProps,
	}, err
}

// Validate ensures that the passed data conforms to the given schema
func (s *JSONSchema) Validate(data map[string]interface{}) error {
	res, err := s.s.Validate(gojsonschema.NewGoLoader(data))
	if err != nil {
		return err
	}
	if !res.Valid() {
		return errors.New(res.Errors()[0].String())
	}
	return nil
}

// ValidateWithDefaults both validates the given data, and inserts defaults for any missing
// values in the root object
func (s *JSONSchema) ValidateWithDefaults(data map[string]interface{}) (err error) {
	// The actual validation happens here
	defer func() {
		err = s.Validate(data)
	}()

	// Insert defaults into the object wherever the data is not provided
	propMapV, ok := s.Schema["properties"]
	if !ok {
		return
	}
	propMap, ok := propMapV.(map[string]interface{})
	if !ok {
		return
	}
	for k, v := range propMap {
		_, ok := data[k]
		if !ok {
			vmap, ok := v.(map[string]interface{})
			if ok {
				dval, ok := vmap["default"]
				if ok {
					data[k] = dval
				}
			}
		}
	}
	return
}

// ValidateUpdate checks an update struct for validity
func (s *JSONSchema) ValidateUpdate(data map[string]interface{}) (err error) {
	for k, v := range data {
		jsp, ok := s.props[k]
		if !ok {
			if !s.additionalProps {
				return fmt.Errorf("Property '%s' not permitted", k)
			}
		} else {
			if v == nil {
				if jsp.required {
					return fmt.Errorf("Property '%s' can't be deleted", k)
				}
			} else {
				res, err := jsp.s.Validate(gojsonschema.NewGoLoader(v))
				if err != nil {
					return err
				}
				if !res.Valid() {
					return errors.New(res.Errors()[0].String())
				}
			}
		}
	}
	return nil
}
