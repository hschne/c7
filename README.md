<div align="center">

# c7

### Context7 from the terminal. Single binary, instant docs.

[![CI](https://github.com/hschne/c7/actions/workflows/ci.yml/badge.svg)](https://github.com/hschne/c7/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/hschne/c7)](https://github.com/hschne/c7/releases/latest)
[![License](https://img.shields.io/github/license/hschne/c7)](https://github.com/hschne/c7/blob/main/LICENSE)

</div>

## What's this?

A lightweight CLI for [Context7](https://context7.com) — fetch up-to-date library documentation without leaving the terminal. No Node.js, no dependencies, just a single binary.

## Install

### Install script (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh
```

Detects your OS and architecture, downloads the right binary from GitHub Releases, and installs it to `/usr/local/bin` by default.

Want a user-local install instead?

```bash
curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh -s -- --local
```

You can also customize with environment variables:

```bash
# Install a specific version
C7_VERSION=v0.2.0 curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh

# Install to a different directory
C7_INSTALL=$HOME/.local/bin curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh
```

### Download from GitHub Releases

Grab the right archive for your platform from [Releases](https://github.com/hschne/c7/releases/latest), extract, and put the binary somewhere on your `PATH`:

```bash
tar -xzf c7_*_linux_amd64.tar.gz
install -m 755 c7 /usr/local/bin/c7
```

### Build from source

```bash
go install github.com/hschne/c7@latest
```

## Usage

### `c7 get` — one-shot lookup

Resolves the library name, fetches docs, done. Caches the library ID for faster repeat lookups.

```bash
c7 get rails "active record scopes"
c7 get hotwire "form submission with turbo frames"
c7 get kamal "deploy with secrets" --tokens 8000
```

### `c7 search` — find a library

```bash
c7 search rails
c7 search "ruby on rails" "active record"
```

### `c7 docs` — fetch docs by library ID

```bash
c7 docs /rails/rails "how to write custom validations"
c7 docs /vercel/next.js "middleware" --topic routing --page 2
```

### `c7 cache clear` — clear cached lookups

```bash
c7 cache clear
```

## Authentication

Works without an API key for basic usage. For higher rate limits, get a key at [context7.com/dashboard](https://context7.com/dashboard):

```bash
export CONTEXT7_API_KEY="your-key-here"
```

## Caching

`c7 get` caches resolved library IDs under your user cache directory (`$XDG_CACHE_HOME/c7/` or equivalent). Entries expire after 7 days. Clear manually with `c7 cache clear`.

## Shell completions

```bash
c7 completion bash   # also: zsh, fish, powershell

# Example: add to ~/.bashrc
echo 'eval "$(c7 completion bash)"' >> ~/.bashrc
```

## Development

```bash
make build     # build the binary
make test      # run tests
make lint      # run golangci-lint
make clean     # remove build artifacts
```

## Releases

`c7` follows [Semantic Versioning](https://semver.org/). Releases are automated via [GoReleaser](https://goreleaser.com/) — pushing a tag like `v0.2.0` triggers a GitHub Actions workflow that builds binaries for Linux and macOS (amd64 + arm64), creates archives and checksums, and publishes a GitHub Release.

See [docs/releasing.md](docs/releasing.md) for maintainer instructions.

## License

[MIT](LICENSE)
