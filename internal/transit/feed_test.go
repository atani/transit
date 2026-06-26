package transit

import "testing"

func TestFeedLabel(t *testing.T) {
	tests := map[string]string{
		"埼京線":                       "埼京線", // already a clean line name
		"湘南新宿ライン":                   "湘南新宿ライン",
		"odpt.Operator:TokyoMetro":  "東京メトロ", // known operator
		"odpt.Operator:Toei":        "都営",
		"odpt.Operator:JR-East":     "JR東日本",
		"odpt.Operator:UnknownCorp": "UnknownCorp", // unknown: strip prefix only
		"":                          "",
	}
	for in, want := range tests {
		if got := FeedLabel(in); got != want {
			t.Errorf("FeedLabel(%q) = %q, want %q", in, got, want)
		}
	}
}
