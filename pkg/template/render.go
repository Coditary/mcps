package template

import (
	"fmt"
	"regexp"
	"strings"
)

var placeholder = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)

// Render replaces {{key}} placeholders using stringified values from vars.
func Render(text string, vars map[string]string) string {
	return placeholder.ReplaceAllStringFunc(text, func(match string) string {
		key := match[2 : len(match)-2]
		if value, ok := vars[key]; ok {
			return value
		}
		return match
	})
}

// RenderArgs renders each argument template and drops empty optional segments.
func RenderArgs(args []string, vars map[string]string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		rendered := strings.TrimSpace(Render(arg, vars))
		if rendered == "" {
			continue
		}
		out = append(out, rendered)
	}
	return out
}

// Stringify converts tool argument values to strings for templating.
func Stringify(args map[string]any) (map[string]string, error) {
	out := make(map[string]string, len(args))
	for key, value := range args {
		switch v := value.(type) {
		case nil:
			out[key] = ""
		case string:
			out[key] = v
		case bool:
			if v {
				out[key] = "true"
			} else {
				out[key] = ""
			}
		case float64:
			if v == float64(int64(v)) {
				out[key] = fmt.Sprintf("%d", int64(v))
			} else {
				out[key] = fmt.Sprintf("%v", v)
			}
		case int:
			out[key] = fmt.Sprintf("%d", v)
		case []any:
			parts := make([]string, 0, len(v))
			for _, item := range v {
				parts = append(parts, fmt.Sprint(item))
			}
			out[key] = strings.Join(parts, " ")
		default:
			out[key] = fmt.Sprint(v)
		}
	}
	return out, nil
}
