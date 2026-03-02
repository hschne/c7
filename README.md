# c7 — Context7 CLI in Go

A lightweight, zero-dependency CLI for [Context7](https://context7.com) — fetch up-to-date library documentation from the terminal. No Node.js required, single binary, instant startup.

## Install

```bash
# Build from source (requires Go 1.21+)
go build -o c7 .
sudo mv c7 /usr/local/bin/

# Or with make
make install
```

## Usage

```bash
# Set your API key for higher rate limits (optional but recommended)
export CONTEXT7_API_KEY="your-key-here"  # get at context7.com/dashboard
```

### `c7 get` — one-shot lookup (easiest)

Resolves the library name, then fetches docs in one command.

```bash
c7 get rails "active record scopes"
c7 get hotwire "form submission with turbo frames"
c7 get kamal "deploy with secrets" --tokens 8000
c7 get tailwindcss "responsive grid layout"
```

### `c7 search` — find a library ID

Useful when you want to know the exact library ID before querying.

```bash
c7 search rails
c7 search "ruby on rails" "active record"
```

Output:

```
ID                             TRUST    NAME
──────────────────────────────────────────────────────────────────────
/rails/rails                   95       Rails
                                        Full-stack web framework for Ruby
/heartcombo/devise             88       Devise
                                        Authentication solution for Rails
```

### `c7 docs` — fetch docs by known library ID

When you already know the ID (from `c7 search` or context7.com).

```bash
c7 docs /rails/rails "how to write custom validations"
c7 docs /vercel/next.js "middleware" --topic routing --page 2
c7 docs /hotwire-dev/turbo "streams" --tokens 10000
```

**Flags:**

- `--tokens N` — max tokens to return (default 5000, min 1000)
- `--page N` — pagination, 1–10 (default 1)
- `--topic TOPIC` — focus docs on a specific topic

## How it works

Context7 exposes a REST API at `context7.com/api/v2`. This CLI wraps two endpoints directly — no MCP protocol overhead, no Node.js process:

| Endpoint                  | What it does                 |
| ------------------------- | ---------------------------- |
| `GET /api/v2/libs/search` | Search libraries by name     |
| `GET /api/v2/context`     | Fetch documentation snippets |

`c7 get` combines both in sequence: resolve → fetch → print.

## Testing

```bash
make test

# Or directly
go test -v ./...
```

## Cross-platform builds

```bash
make build-all
# Produces: dist/c7-linux-amd64, c7-darwin-arm64, c7-windows-amd64.exe, etc.
```

## Why Go vs MCPorter/npx?

|              | `npx mcporter`        | `c7` (Go)         |
| ------------ | --------------------- | ----------------- |
| Startup time | ~2–4s (Node spawn)    | ~20ms             |
| Dependencies | Node.js + npm         | None              |
| Binary size  | N/A (always fetches)  | ~6MB static       |
| Portable     | No                    | Yes (single file) |
| API overhead | MCP JSON-RPC protocol | Direct HTTP REST  |
