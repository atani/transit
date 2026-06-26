# transit

`transit` は日本の経路検索と発車案内を行う小さな Go 製 CLI です。公開 Transit API（<https://api.transit.ls8h.com/>）を利用します。MCP サーバーも計画中です。

## インストール

### Homebrew（macOS / Linux）

```bash
brew install atani/tap/transit
```

### winget（Windows）

```powershell
winget install atani.transit
```

### ソースから

```bash
go install github.com/atani/transit/cmd/transit@latest
```

## 使い方

### 駅名サジェスト

```bash
transit suggest 渋谷 --limit 3
```

`駅名 / よみ / 路線（事業者）` をタブ区切りで表示します。同名駅は路線ごとに複数並びます。

### 経路検索

```bash
transit plan 渋谷 新宿 --time 09:00 --num 1
```

オプションは次のとおりです。

- `--date YYYYMMDD` 日付
- `--time HH:MM` 基準時刻
- `--type departure|arrival|first|last` 出発基準・到着基準・始発・終電
- `--num N` 候補数

### 発車案内

```bash
transit departures 渋谷 --limit 5
```

`発車時刻 / 種別 / 路線（方面）-> 行先` を時刻順に表示します。

### 共通オプション

`--json` を付けると構造化 JSON を出力します。スクリプトやエージェント連携に使えます。`transit version` でバージョンを表示します。

## 駅の指定方法

駅名はそのまま渡せば `suggest` で解決されます。上級者は解決を飛ばして直接指定もできます。

- `geo:35.681,139.767`（緯度経度）
- フィード修飾の駅 ID（例: `scrape-fukuoka-subway:福岡市-空港線-博多`）

## API について

時刻はサービス日の午前0時からの秒数です。終電後は `01:30(+1d)`、前日扱いは `23:30(-1d)` のように日跨ぎを表示します。

経路はフィード（事業者）単位で計算されます。API は事業者をまたぐ乗換を計算しません。`plan` に駅名を渡すと、共通フィードを持つ駅の組に解決して経路を探します。出発・到着が別事業者にしか存在しない場合（例: 西鉄の駅と新幹線の駅）は直通経路が無く、結果は徒歩のみになります。その場合は CLI が注記を表示します。

環境変数 `TRANSIT_API_BASE_URL` で API のベース URL を差し替えられます。
