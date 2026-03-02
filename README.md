# c7 — Context7 CLI

A lightweight CLI for [Context7](https://context7.com) — fetch up-to-date library documentation from the terminal. Single binary, no Node.js required.

## Install

```bash
go install github.com/hschne/c7@latest

# Or build from source
make install
```

## Usage

```bash
# Optional: set API key for higher rate limits
export CONTEXT7_API_KEY="your-key-here"  # get at context7.com/dashboard
```

### `c7 get` — one-shot lookup (easiest)

Resolves the library name, then fetches docs in one command. Caches the resolved library ID for faster repeat lookups.

```bash
c7 get rails "active record scopes"
c7 get hotwire "form submission with turbo frames"
c7 get kamal "deploy with secrets" --tokens 8000
```

### `c7 search` — find a library ID

```bash
c7 search rails
c7 search "ruby on rails" "active record"
```

### `c7 docs` — fetch docs by known library ID

```bash
c7 docs /rails/rails "how to write custom validations"
c7 docs /vercel/next.js "middleware" --topic routing --page 2
```

### `c7 cache clear` — clear cached lookups

```bash
c7 cache clear
```

### Shell completions

```bash
c7 completion bash  # also: zsh, fish, powershell
```

## Testing

```bash
make test
```

## Project structure

```
main.go          # entrypoint, calls cmd.Execute()
cmd/             # cobra command definitions
  root.go
  get.go
  search.go
  docs.go
  cache.go
internal/        # business logic (API client, cache, formatting)
  api.go
  cache.go
  format.go
test/            # tests
  api_test.go
  cache_test.go
  format_test.go
```

## Cross-platform builds

```bash
make build-all
```
