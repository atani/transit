# リリース

バージョン管理は [release-please](https://github.com/googleapis/release-please) が行います。
配布物のビルドと公開は `.github/workflows/release-please.yml` の各ジョブが担います。
構成は他の Go ツール（ctxpack など）と揃えています。

## 流れ

1. [Conventional Commits](https://www.conventionalcommits.org)（`feat:` / `fix:` など）で `main` にマージします。
2. release-please がリリース PR を作成・更新します（`.release-please-manifest.json` のバージョンと `CHANGELOG.md` を更新）。
3. そのリリース PR をマージすると `vX.Y.Z` のタグと GitHub Release が作られます。
4. 続けて以下のジョブが走ります。
   - `build-release`: 6 プラットフォームの ZIP と winget マニフェスト束を作り、Release に添付。
   - `publish-homebrew`: `atani/homebrew-tap` の `Formula/transit.rb` をソースビルド方式で更新。
   - `publish-winget`: `microsoft/winget-pkgs` へ更新 PR を送信（`ENABLE_WINGET` 有効時のみ）。

## 必要な secret / variable

- `HOMEBREW_TAP_GITHUB_TOKEN`: `atani/homebrew-tap` に push できる PAT。release-please の PR 作成と formula 更新に使います。
- `WINGET_TOKEN`: `atani/winget-pkgs` フォークに push できる PAT。winget 送信に使います。
- `ENABLE_WINGET`: リポジトリ変数。`true` で winget 自動送信を有効化します。

winget の初回登録だけは手動です。手順は `packaging/winget/README.md` を参照してください。

## ローカル確認

```bash
go build ./... && go vet ./... && go test ./...
```
