package validator

import (
	_ "embed"
	"encoding/json"
	"errors"

	jsonSchema "github.com/xeipuuv/gojsonschema"
	multiErr "go.uber.org/multierr"
)

//go:embed schema.json
var schemaJSON []byte

// ValidatorJSON validates request to defined json schema
type ValidatorJSON interface {
	Validate(raw json.RawMessage) error
}
type validatorImpl struct {
	schema *jsonSchema.Schema
}

// New returns new instance of ValidatorJSON
func New() (ValidatorJSON, error) {
	s, err := jsonSchema.NewSchema(jsonSchema.NewBytesLoader(schemaJSON))
	if err != nil {
		return nil, err
	}
	return &validatorImpl{
		schema: s,
	}, nil
}

// Validate implements ValidatorJSON interface
func (s *validatorImpl) Validate(raw json.RawMessage) error {
	res, err := s.schema.Validate(jsonSchema.NewBytesLoader(raw))
	if err != nil {
		return err
	}
	if res.Valid() {
		return nil
	}
	for _, e := range res.Errors() {
		err = multiErr.Append(err, errors.New(e.Field()+": "+e.Description()))
	}
	return err
}
