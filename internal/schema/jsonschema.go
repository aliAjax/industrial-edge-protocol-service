package schema

import (
	"fmt"
	"reflect"
)

type Field struct {
	Name     string
	Type     string
	Required bool
}
type Schema struct {
	Name   string
	Fields []Field
}

func (s Schema) Validate(value map[string]any) error {
	for _, f := range s.Fields {
		v, ok := value[f.Name]
		if !ok {
			if f.Required {
				return fmt.Errorf("missing field %s", f.Name)
			}
			continue
		}
		if !compatible(v, f.Type) {
			return fmt.Errorf("field %s expects %s", f.Name, f.Type)
		}
	}
	return nil
}
func compatible(v any, want string) bool {
	kind := reflect.TypeOf(v)
	if kind == nil {
		return false
	}
	switch want {
	case "string":
		return kind.Kind() == reflect.String
	case "number":
		return kind.Kind() >= reflect.Int && kind.Kind() <= reflect.Float64
	case "boolean":
		return kind.Kind() == reflect.Bool
	case "object":
		return kind.Kind() == reflect.Map
	default:
		return false
	}
}
func (s Schema) RequiredFields() []string {
	out := []string{}
	for _, f := range s.Fields {
		if f.Required {
			out = append(out, f.Name)
		}
	}
	return out
}
