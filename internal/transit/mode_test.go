package transit

import "testing"

func TestModeLabel(t *testing.T) {
	tests := map[string]string{
		"rail":    "電車",
		"train":   "電車",
		"subway":  "地下鉄",
		"bus":     "バス",
		"walk":    "徒歩",
		"ferry":   "フェリー",
		"":        "",
		"skyhook": "skyhook", // unknown modes pass through unchanged
	}
	for in, want := range tests {
		if got := ModeLabel(in); got != want {
			t.Errorf("ModeLabel(%q) = %q, want %q", in, got, want)
		}
	}
}
