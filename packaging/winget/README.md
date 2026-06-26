# winget packaging

Each release attaches Windows ZIP artifacts (x64 + arm64) and a generated
winget manifest bundle (`transit_<version>_winget_manifests.zip`).

## First submission (manual, once)

A brand-new package must be introduced to the community repository by hand:

1. Download `transit_<version>_winget_manifests.zip` from the release.
2. Unzip it into a local clone of `microsoft/winget-pkgs` under
   `manifests/a/atani/transit/<version>/`.
3. Validate and submit:

   ```powershell
   winget validate manifests\a\atani\transit\<version>\
   wingetcreate submit manifests\a\atani\transit\<version>\
   ```

Once Microsoft merges the manifest, users can install with:

```powershell
winget install atani.transit
```

## Subsequent versions (automated)

The `publish-winget` job in `.github/workflows/release-please.yml` opens the
update PR to `microsoft/winget-pkgs` automatically on each release, using
[`vedantmgoyal9/winget-releaser`](https://github.com/vedantmgoyal9/winget-releaser).
It is dormant until you enable it:

1. Fork `microsoft/winget-pkgs` to `atani/winget-pkgs`.
2. Create a classic PAT with `public_repo` scope that can push to that fork,
   and store it as the `WINGET_TOKEN` repository secret.
3. Set the repository variable `ENABLE_WINGET` to `true`.

After that, every release regenerates manifests from the Windows zips and
submits the winget-pkgs PR without manual steps.
