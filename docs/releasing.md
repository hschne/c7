# Releasing c7

## How it works

1. Development happens on `main` through pull requests. CI runs tests, vet, lint, and a GoReleaser config check on every push/PR.
2. When ready to release, a maintainer creates and pushes a semantic version tag.
3. The `release.yml` workflow triggers on the tag, runs GoReleaser, and publishes a GitHub Release with:
   - Linux and macOS binaries (`amd64` + `arm64`)
   - `.tar.gz` archives
   - `checksums.txt`
   - Auto-generated changelog

## Cutting a release

```bash
# Make sure main is up to date
git checkout main
git pull

# Tag the release
git tag v0.2.0
git push origin v0.2.0
```

Pick the version number according to [Semantic Versioning](https://semver.org/):

- **patch** (`v0.1.1`): bug fixes, doc tweaks
- **minor** (`v0.2.0`): new features, non-breaking changes
- **major** (`v1.0.0`): breaking changes to CLI flags, commands, or output format

## Verifying a release

After pushing the tag:

1. Check the [Actions tab](https://github.com/hschne/c7/actions/workflows/release.yml) — the release workflow should complete in a few minutes.
2. Verify the [Releases page](https://github.com/hschne/c7/releases) shows the new version with all expected assets.
3. Test the install script: `curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh`
4. Confirm `c7 --version` prints the correct version, commit, and build date.

## Asset naming

GoReleaser produces archives named `c7_{version}_{os}_{arch}.tar.gz` (e.g., `c7_0.2.0_linux_amd64.tar.gz`). The install script depends on this naming convention — if it changes, the install script must be updated to match.

## Local dry run

To test the release pipeline locally without publishing:

```bash
goreleaser build --snapshot --clean   # build only
goreleaser release --snapshot --clean # full dry run (no publish)
```

## Important

Tags, releases, and publishing are **maintainer-only** actions. Automated agents must never create tags, trigger the release workflow, or publish releases.
