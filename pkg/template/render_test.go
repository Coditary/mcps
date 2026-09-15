package template

import "testing"

func TestRenderArgs(t *testing.T) {
	vars := map[string]string{"name": "world"}
	args := RenderArgs([]string{"hello", "{{name}}"}, vars)
	if len(args) != 2 || args[1] != "world" {
		t.Fatalf("unexpected args: %#v", args)
	}
}
