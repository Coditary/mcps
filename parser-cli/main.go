package main

import (
	_ "embed"

	"github.com/coditary/mcps/pkg/run"
)

//go:embed tools.yaml
var toolsYAML []byte

func main() {
	run.Main(toolsYAML, "tools.yaml")
}
