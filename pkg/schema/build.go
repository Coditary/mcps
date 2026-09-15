package schema

import (
	"encoding/json"

	"github.com/coditary/mcps/pkg/config"
)

const confirmField = "confirm"

// BuildInputSchema converts YAML input definitions to JSON Schema.
func BuildInputSchema(tool config.ToolSpec) map[string]any {
	properties := map[string]any{}
	required := make([]string, 0)

	for name, field := range tool.Inputs {
		prop := map[string]any{
			"type":        normalizeType(field.Type),
			"description": field.Description,
		}
		if len(field.Enum) > 0 {
			prop["enum"] = field.Enum
		}
		if field.Default != nil {
			prop["default"] = field.Default
		}
		properties[name] = prop
		if field.Required {
			required = append(required, name)
		}
	}

	if tool.Mutating {
		properties[confirmField] = map[string]any{
			"type":        "boolean",
			"description": "Set true to run without dry-run. Defaults to false.",
			"default":     false,
		}
	}

	schema := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// BuildInputSchemaJSON returns schema as json.RawMessage.
func BuildInputSchemaJSON(tool config.ToolSpec) (json.RawMessage, error) {
	schema := BuildInputSchema(tool)
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func normalizeType(t string) string {
	if t == "" {
		return "string"
	}
	return t
}
