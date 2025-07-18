// Package stringid is used to generate strictly increasing, unique string-based IDs
// starting from a given numeric value. These IDs are sortable using plain string sorting,
// either in ascending or descending order.
//
// This allows consistent, predictable ordering of keys (such as treasures) within a swamp.
package stringid

// New generates a unique string ID suitable for lexicographic (ascending) ordering,
// based on the given numeric input. The resulting string looks like:
// 0 -> aaaaaaa, 1 -> aaaaaab, etc.
//
// Swamp returns the number of treasures within a swamp and the next ID to be used.
func New(num int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 7)
	for i := 6; i >= 0; i-- {
		result[i] = chars[num%26]
		num /= 26
	}
	return string(result)
}
