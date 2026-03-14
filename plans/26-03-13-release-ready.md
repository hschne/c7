# Release Ready

## Goal

Make `c7` feel like a proper, releasable CLI: solid CI, polished docs, automated semantic releases, multi-architecture binaries for Linux and macOS, installable release artifacts attached directly to GitHub releases, and a simple install script that fetches the right asset from GitHub Releases.

## Current State

`c7` is a small Go CLI with a clean command structure and tests already in place.

Relevant files:

- `main.go` embeds a `version` variable passed into `cmd.Execute(version)`
- `cmd/root.go` exposes Cobra root metadata and `--version`
- `cmd/get.go`, `cmd/search.go`, `cmd/docs.go`, `cmd/cache.go` implement the user-facing CLI
- `internal/api.go` talks to the Context7 API
- `internal/cache.go` stores cached library lookups under the user cache dir
- `README.md` is functional but still closer to a developer README than a release-grade project homepage
- `Makefile` supports local builds and a manual `build-all`, but not release automation

Repository status today:

- No GitHub Actions workflows are present
- No GitHub releases exist yet
- No changelog or automated release notes flow exists
- Versioning is manually set in `Makefile` (`VERSION := 0.1.0`) and injected into `main.version`
- Cross-platform builds exist only as a manual Make target and do not cover Linux ARM64 or archives/checksums
- No packaging/distribution setup exists for GitHub release artifacts or an install script that can fetch the right release asset automatically

## New Behavior

Releasing `c7` should be straightforward and repeatable:

1. Development changes land through PRs with CI running tests, `go vet`, linting, and build verification.
2. A maintainer creates and pushes a semantic version tag like `v0.2.0`.
3. GitHub Actions runs GoReleaser on the tag.
4. GoReleaser publishes a GitHub release with:
   - Linux/macOS binaries for `amd64` and `arm64`
   - archives and checksums
   - generated changelog/release notes
5. This repository ships an install script that detects the user platform, looks at GitHub releases, downloads the matching asset, and installs `c7` without requiring PPAs or other package repositories.
6. The README clearly explains what the tool is, how to install it from GitHub releases or via the install script, how versioning works, and how to use it.
7. Homebrew or other package repositories can be added later without changing the core release flow.

## Decisions

### How should releases be built and published?

Use GoReleaser, triggered by pushing semantic version tags.

**Why:** This is the standard, lowest-maintenance path for a Go CLI. It gives us multi-arch builds, archives, checksums, GitHub releases, Homebrew support, and Debian packages in one place.

### What platforms should the first release support?

Support Linux and macOS only, both `amd64` and `arm64`.

**Why:** This matches the desired initial audience while keeping the matrix small and maintainable. Windows can be added later without redesigning the release flow.

### What release/versioning policy should `c7` follow?

Adopt Semantic Versioning immediately, using `vMAJOR.MINOR.PATCH` tags.

**Why:** The user explicitly wants SemVer. Even if the CLI is still young, setting the contract now makes releases easier to understand and automations easier to reason about.

### How polished should CI be in this pass?

Medium polish: tests, `go vet`, linting, and release dry-run/build verification.

**Why:** This provides meaningful quality gates without turning the project into a CI science project.

### Which distribution channels are in scope now?

Keep the first distribution pass simple: publish release artifacts to GitHub Releases and provide an install script in this repository that automatically downloads the matching asset.

**Why:** This is the lowest-friction way to ship release-ready binaries quickly without taking on package repository maintenance. Package repositories like Homebrew can come later once the release pipeline is stable.

### How should changelogs be handled?

Generate changelogs automatically as part of releases.

**Why:** Manual changelogs are easy to neglect. GoReleaser can provide a consistent starting point and reduce release friction.

### What install path should users have on day one?

Users should be able to either download binaries directly from GitHub Releases or run an install script from this repository that automatically picks the right release asset.

**Why:** This keeps the initial rollout KISS while still providing a polished installation story. Users get a one-command install path without requiring us to stand up and maintain package repositories.

## Implementation Plan

1. **Define release metadata and version injection**
   - Update build/version handling so release artifacts embed at least:
     - version
     - commit
     - build date
   - Remove duplicated/manual version assumptions from `Makefile` where possible.
   - Ensure `c7 --version` prints useful release information.
   - Files likely involved:
     - `main.go`
     - `cmd/root.go`
     - `Makefile`

2. **Add CI with medium-quality checks**
   - Add GitHub Actions workflow(s) for:
     - `go test ./...`
     - `go vet ./...`
     - linting (likely `golangci-lint`)
     - build verification for Linux/macOS targets or a Go matrix plus normal build
     - GoReleaser dry-run on non-tag changes to catch release config breakage early
   - Decide Go version matrix (for example current stable + previous stable, or a pinned single version if simplicity wins).
   - Files to add:
     - `.github/workflows/ci.yml`
     - maybe `.golangci.yml`

3. **Add GoReleaser for tagged releases**
   - Create `.goreleaser.yml` with:
     - GitHub release publishing
     - Linux/macOS builds for `amd64` and `arm64`
     - tarballs/checksums
     - changelog generation
   - Configure tag-based releases only.
   - Add release archives with sensible naming conventions.
   - Ensure release assets are named consistently so an install script can select them predictably.
   - Ensure ldflags populate version metadata from GoReleaser.
   - Files to add/update:
     - `.goreleaser.yml`
     - `Makefile`

4. **Add an install script that installs from GitHub Releases**
   - Add a shell install script in the repository (for example `scripts/install.sh` or `install.sh`) that:
     - detects OS and architecture
     - resolves the latest release or a requested version
     - downloads the matching GitHub release asset
     - unpacks it
     - installs `c7` into a conventional destination
   - Keep the script dependency-light and easy to audit.
   - Document supported environment variables/flags for version selection and install destination.
   - Files to add/update:
     - `install.sh` or `scripts/install.sh`
     - `README.md`

5. **Rewrite the README in the project’s established tone**
   - Rewrite `README.md` to feel closer to the user’s pinned projects:
     - concise positioning
     - polished top section
     - badges with actual signal
     - strong install section
     - clear examples early
     - practical release/install options
   - Add sections for:
     - What’s this?
     - Getting Started / Install
     - Install script
     - Install from GitHub Releases
     - Usage
     - Authentication and rate limits
     - Caching behavior
     - Shell completions
     - Development
     - Releases / Versioning
   - Add concrete copy-paste install examples for the install script and for downloading binaries from release assets.
   - Keep the voice compact, confident, and a little playful without overdoing it.

6. **Add release operations documentation**
   - Add a small maintainer doc describing:
     - how to cut a release
     - how tags trigger publishing
     - how GitHub release artifacts are produced and expected to be consumed
     - that install script behavior depends on predictable asset naming in GitHub Releases
     - that the coding agent must never create tags, trigger publishing, or cut releases on the maintainer’s behalf
   - This can live in `README.md`, `CONTRIBUTING.md`, or `docs/releasing.md`.

7. **Do one end-to-end dry run before the first real tag**
   - Run local checks only:
     - tests
     - lint
     - GoReleaser dry-run
   - Verify artifact names, archive contents, checksums, install script behavior, and version output.
   - Fix anything awkward before the maintainer cuts the first public release.
   - Under no circumstances should the agent create or push the initial tag, trigger the release workflow, or publish the release.

## What We're NOT Doing

- Windows builds in this pass
- Ubuntu PPA setup
- Any package repository setup in this pass
- Homebrew tap setup in this pass
- Official Arch repository packaging
- Fully automated AUR publishing
- Advanced supply-chain work like signing, SBOMs, provenance attestations, or cosign integration
- Coverage reporting or Code Climate integration unless explicitly requested later
- A full docs site

## Risks & Edge Cases

- **Install script robustness:** OS/arch detection, download URL construction, and install destinations need to be straightforward and well-tested or users will hit confusing setup failures.
- **Asset naming stability matters:** if release asset names change, the install script and README examples can break.
- **SemVer expectations:** committing to SemVer means CLI flags, command names, output contracts, and environment variable behavior should be treated as public API.
- **Version injection drift:** today the version is partly handled via `Makefile`; if local builds and GoReleaser builds diverge, `--version` output may become inconsistent.
- **GoReleaser dry-run limitations:** some integration pieces are only fully validated against real repos/tags and actual GitHub releases.
- **README install instructions can rot:** install snippets and asset references should be checked when the first release is cut.

## Follow-Ups

- Add signed releases and/or checksums verification guidance
- Add SBOM/provenance if desired
- Add Windows artifacts if demand appears
- Consider a `c7 completion` section with install helpers per shell
- Add Homebrew distribution once the GitHub-release-based flow is stable
- Revisit whether `v1.0.0` needs a documented CLI compatibility policy beyond plain SemVer
