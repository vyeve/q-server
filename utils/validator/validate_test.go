package validator

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"go.uber.org/fx"
)

func TestValidateJSON(t *testing.T) {
	var (
		err    error
		p      json.RawMessage
		ctx    = context.Background()
		schema ValidatorJSON
	)
	app := fx.New(
		fx.Provide(New),
		fx.Populate(&schema),
	)
	if err = app.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer app.Stop(ctx) // nolint: errcheck
	p, err = ioutil.ReadFile("../assets/sample1.json")
	if err != nil {
		t.Fatal(err)
	}
	err = schema.Validate(p)
	if err != nil {
		t.Error(err)
	}
	p, err = ioutil.ReadFile("../assets/sample2.json")
	if err != nil {
		t.Fatal(err)
	}
	err = schema.Validate(p)
	if err == nil {
		t.Error("expected not <nil> error")
	}
}
