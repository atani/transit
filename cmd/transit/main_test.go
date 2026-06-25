package main

import (
	"reflect"
	"testing"
)

func TestReorderFlags(t *testing.T) {
	boolFlags := map[string]bool{"json": true}
	tests := map[string]struct {
		in   []string
		want []string
	}{
		"flags after positional": {
			in:   []string{"渋谷", "新宿", "--time", "09:00"},
			want: []string{"--time", "09:00", "渋谷", "新宿"},
		},
		"value flag with equals": {
			in:   []string{"渋谷", "--limit=5"},
			want: []string{"--limit=5", "渋谷"},
		},
		// Regression: a single-dash bool flag must not swallow the next
		// positional argument as if it were the flag's value.
		"single dash bool keeps positional": {
			in:   []string{"-json", "渋谷"},
			want: []string{"-json", "渋谷"},
		},
		"double dash bool keeps positional": {
			in:   []string{"渋谷", "--json"},
			want: []string{"--json", "渋谷"},
		},
		"value flag at end without value": {
			in:   []string{"渋谷", "--time"},
			want: []string{"--time", "渋谷"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := reorderFlags(tt.in, boolFlags); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("reorderFlags(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
