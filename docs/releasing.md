# Releasing

Releases are driven by [release-please](https://github.com/googleapis/release-please)
and built by [GoReleaser](https://goreleaser.com).

## Flow

1. Merge changes to `main` using [Conventional Commits](https://www.conventionalcommits.org)
   (`feat:`, `fix:`, `feat!:`, ...).
2. `release-please` opens (and keeps updating) a release pull request that bumps
   the version in `.release-please-manifest.json` and updates `CHANGELOG.md`.
3. Merging that release pull request creates the `vX.Y.Z` tag and GitHub release.
4. The `goreleaser` job then attaches the build artifacts to that release:
   cross-platform binaries (macOS / Linux / Windows on amd64 / arm64), the
   Homebrew cask in [`atani/homebrew-tap`](https://github.com/atani/homebrew-tap),
   and a winget manifest pull request to
   [`microsoft/winget-pkgs`](https://github.com/microsoft/winget-pkgs).

GoReleaser runs with `release.mode: keep-existing`, so it never creates a second
release; it only uploads to the one release-please made.

## One-time setup

- Fork `microsoft/winget-pkgs` to `atani/winget-pkgs`.
- Create a personal access token with `repo` scope on `atani/homebrew-tap` and
  `atani/winget-pkgs`, and add it as the `TAP_GITHUB_TOKEN` repository secret.

## Local validation

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

`version` is injected at build time, so `transit version` prints the released
tag (it prints `dev` for local `go build` / `go run`).
