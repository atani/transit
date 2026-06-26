# Releasing

Versioning is tracked by [release-please](https://github.com/googleapis/release-please).
Binaries and the winget manifest are built by [GoReleaser](https://goreleaser.com).
Homebrew is a source-build formula in [`atani/homebrew-tap`](https://github.com/atani/homebrew-tap).

## Homebrew (source formula)

`Formula/transit.rb` in the tap builds from the release source tarball with
`go build` and injects the version, matching the rest of the tap. To cut a
release and update the formula:

```bash
# 1. tag and create the GitHub release (source tarball is auto-generated)
gh release create vX.Y.Z --repo atani/transit --target main --title vX.Y.Z --notes "..."

# 2. compute the tarball checksum
curl -sL -o /tmp/transit.tar.gz \
  "https://github.com/atani/transit/archive/refs/tags/vX.Y.Z.tar.gz"
shasum -a 256 /tmp/transit.tar.gz

# 3. update url + sha256 in atani/homebrew-tap Formula/transit.rb, then:
brew upgrade atani/tap/transit
```

This needs no secrets because Homebrew builds from source on the user's machine.

## winget (GoReleaser)

`release-please` opens a release pull request from Conventional Commits.
Merging it creates the tag and release, then the `goreleaser` job builds the
cross-platform binaries and opens a winget manifest pull request to
[`microsoft/winget-pkgs`](https://github.com/microsoft/winget-pkgs).

This path needs `TAP_GITHUB_TOKEN`: a personal access token with `repo` scope
on the `atani/winget-pkgs` fork. GoReleaser runs with `release.mode:
keep-existing`, so it only attaches artifacts to the release-please release.

## Local validation

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

`version` is injected at build time, so `transit version` prints the released
tag (it prints `dev` for local `go build` / `go run`).
