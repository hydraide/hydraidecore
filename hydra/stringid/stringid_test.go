package stringid

import (
	"sort"
	"strconv"
	"testing"
)

// TestNew tests the New function with various input values.
func TestNew(t *testing.T) {
	// Input numbers and their expected string ID results to be tested.
	tests := []struct {
		num  int
		want string
	}{
		{0, "aaaaaaa"},
		{1, "aaaaaab"},
		{25, "aaaaaaz"},
		{26, "aaaaaba"},
		{27, "aaaaabb"},
		{676, "aaaabaa"}, // 26^2
		{703, "aaaabbb"}, // 26^2 + 26
		{1000000, "aacexho"},
	}

	// Check if the generated IDs match the expected results.
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.num), func(t *testing.T) {
			got := New(tt.num)
			if got != tt.want {
				t.Errorf("New(%d) = %v, want %v", tt.num, got, tt.want)
			}
		})
	}

	// Verify that the generated IDs are in the correct order.
	ids := make([]string, 0, len(tests))
	for _, tt := range tests {
		ids = append(ids, New(tt.num))
	}
	if !sort.StringsAreSorted(ids) {
		t.Error("Generated IDs are not in sorted order")
	}
}

// goos: windows
// goarch: amd64
// pkg: github.com/hydraide/hydraidecore/hydra/stringid
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkNew
// BenchmarkNew-32         128173084                9.328 ns/op
// PASS
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(i)
	}
}
