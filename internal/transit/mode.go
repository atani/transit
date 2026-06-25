package transit

// ModeLabel returns a Japanese label for a transit mode. Unknown modes are
// returned unchanged so new feed values still display something useful.
func ModeLabel(mode string) string {
	switch mode {
	case "rail", "train":
		return "電車"
	case "subway", "metro":
		return "地下鉄"
	case "bus":
		return "バス"
	case "tram", "streetcar":
		return "路面電車"
	case "ferry", "boat":
		return "フェリー"
	case "monorail":
		return "モノレール"
	case "funicular", "cablecar":
		return "ケーブルカー"
	case "gondola":
		return "ロープウェイ"
	case "walk":
		return "徒歩"
	default:
		return mode
	}
}
