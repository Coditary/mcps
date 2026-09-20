# Coditary MCP Servers

Go-based [Model Context Protocol](https://modelcontextprotocol.io/) servers that wrap Coditary CLI tools. Each server is a separate binary with a focused tool surface for AI clients such as Cursor.

Wuji is intentionally **not** included here. Wuji orchestrates MCP servers itself and would be circular to expose as an MCP target.

## Architecture

```
plugins/mcps/
├── pkg/                 # shared wrapper library (config, exec, schema, server)
├── reqpack/             # declarative MCP (tools.yaml + main.go)
├── parser-cli/
├── ipmc/
├── tempify/
├── prebyte/
├── beez/
├── ycallr/              # custom MCP (dynamic API profiles)
├── go.work
└── Makefile
```

Most servers are **declarative**: `tools.yaml` describes how CLI arguments map to MCP tools. The shared `pkg/server` package loads the YAML, registers tools, executes the wrapped binary, and returns stdout or parsed JSON.

`ycallr` is the exception: it reads API profiles from disk and exposes `list_apis`, `describe_api`, and `call_command`.

## Available MCP servers

| MCP binary | Wrapped CLI | Env override | Notes |
|------------|-------------|--------------|-------|
| `mcp-reqpack` | `rqp` | `REQPACK_BIN` | Package management, audit, SBOM |
| `mcp-parser-cli` | `parser-cli` | `PARSER_CLI_BIN` | Tree-sitter AST parsing |
| `mcp-ipmc` | `ipmc` | `IPMC_BIN` | Impact map rendering |
| `mcp-tempify` | `tempify` | `TEMPIFY_BIN` | Template scaffolding |
| `mcp-prebyte` | `prebyte` | `PREBYTE_BIN` | Template preprocessing |
| `mcp-beez` | `beez` | `BEEZ_BIN` | Build orchestration |
| `mcp-ycallr` | `ycallr` | `YCALLR_BIN` | API profile execution |

### Mutating tools

Tools marked as mutating in `tools.yaml` default to the configured `dry_run_flag` (usually `--dry-run`). Pass `confirm: true` in the tool arguments to run the real command.

## Build

```bash
cd plugins/mcps
make build
```

Binaries are written to `plugins/mcps/bin/`.

## Cursor configuration example

```json
{
  "mcpServers": {
    "reqpack": {
      "command": "/absolute/path/to/Coditary/plugins/mcps/bin/mcp-reqpack"
    },
    "parser-cli": {
      "command": "/absolute/path/to/Coditary/plugins/mcps/bin/mcp-parser-cli"
    }
  }
}
```

Optional override config file:

```bash
mcp-reqpack -config /path/to/custom-tools.yaml
```

## Adding a new declarative MCP

1. Create `plugins/mcps/<name>/` with `go.mod`, `main.go`, and `tools.yaml`.
2. Add the module to `go.work`.
3. Add the name to `SERVERS` in `Makefile`.
4. Define tools in `tools.yaml` using the schema in `pkg/config/config.go`.

Use a custom Go `main.go` when tool discovery or response handling cannot be expressed in YAML (see `ycallr/`).

## Prerequisites

Each MCP server expects its wrapped CLI binary on `PATH`, or set the corresponding `*_BIN` environment variable in the MCP server config.

Tempify and Prebyte MCPs call the standalone `tempify` and `prebyte` binaries. Initialize the app submodules if they are not checked out yet:

```bash
git submodule update --init apps/tempify apps/prebyte
```

## License

Follows the license of the Coditary monorepo and wrapped tools.
