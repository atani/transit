package transit

import "testing"

func TestFeedLabel(t *testing.T) {
	tests := map[string]string{
		"埼京線":                      "埼京線", // already a clean line name
		"湘南新宿ライン":                  "湘南新宿ライン",
		"odpt.Operator:TokyoMetro": "東京メトロ", // known operator
		"odpt.Operator:Toei":       "都営",
		"odpt.Operator:JR-East":    "JR東日本",
		// The upstream API shards one operator across several feeds; the
		// "分割NN" marker is internal, so it is dropped.
		"odpt.Operator:TokyoMetro 分割01": "東京メトロ",
		"odpt.Operator:Toei 分割01":       "都営",
		"odpt.Operator:TokyuBus 分割05":   "東急バス",
		// A meaningful suffix (a line name) is kept.
		"odpt.Operator:TokyoMetro 千代田線": "東京メトロ 千代田線",
		// Unknown operators keep their name, with prefix and shard stripped.
		"odpt.Operator:UnknownCorp":      "UnknownCorp",
		"odpt.Operator:UnknownCorp 分割02": "UnknownCorp",
		"odpt.Operator:UnknownCorp 特急線":  "UnknownCorp 特急線",
		"":                               "",
	}
	for in, want := range tests {
		if got := FeedLabel(in); got != want {
			t.Errorf("FeedLabel(%q) = %q, want %q", in, got, want)
		}
	}
}
