package transit

import "testing"

func TestFormatServiceSeconds(t *testing.T) {
	tests := map[string]struct {
		in   int
		want string
	}{
		"midnight":       {0, "00:00"},
		"morning":        {9*3600 + 5*60, "09:05"},
		"after midnight": {25*3600 + 30*60, "01:30(+1d)"},
		"previous day":   {-30 * 60, "23:30(-1d)"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := FormatServiceSeconds(tt.in); got != tt.want {
				t.Fatalf("FormatServiceSeconds(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
