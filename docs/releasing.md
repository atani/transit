# Releasing

Releases are automated with [GoReleaser](https://goreleaser.com). Pushing a
`v*` tag triggers `.github/workflows/release.yml`, which:

1. Cross-compiles `transit` for macOS, Linux, and Windows (amd64 / arm64).
2. Creates a GitHub release with the archives and `checksums.txt`.
3. Updates the Homebrew cask in [`atani/homebrew-tap`](https://github.com/atani/homebrew-tap).
4. Opens a pull request to [`microsoft/winget-pkgs`](https://github.com/microsoft/winget-pkgs) from the `atani/winget-pkgs` fork.

## One-time setup

- Fork `microsoft/winget-pkgs` to `atani/winget-pkgs`.
- Create a personal access token with `repo` scope on `atani/homebrew-tap` and
  `atani/winget-pkgs`, and add it as the `TAP_GITHUB_TOKEN` repository secret.

## Cutting a release

```bash
git tag v0.1.0
git push origin v0.1.0
```

Verify the config before tagging:

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

`version` is injected at build time, so `transit version` prints the released
tag (it prints `dev` for local `go build` / `go run`).
