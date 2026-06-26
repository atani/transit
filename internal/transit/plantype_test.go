package transit

import "testing"

func TestTypeLabel(t *testing.T) {
	tests := map[string]string{
		"departure": "出発",
		"arrival":   "到着",
		"first":     "始発",
		"last":      "終電",
		"":          "",
		"weird":     "weird",
	}
	for in, want := range tests {
		if got := TypeLabel(in); got != want {
			t.Errorf("TypeLabel(%q) = %q, want %q", in, got, want)
		}
	}
}
