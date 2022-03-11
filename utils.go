package flagx

import "strings"

// splitTrimSpace slices `s` into all substrings separated by `sep`
// and returns a slice of the substrings between those separators,
// with all leading and trailing white space removed, as defined by Unicode.
//
// If `s` does not contain `sep` and `sep` is not empty, splitTrimSpace returns a
// slice of length 1 whose only element is `s` after TrimSpace.
func splitTrimSpace(s string, sep string) []string {
	res := []string{}
	for _, item := range strings.Split(s, sep) {
		ss := strings.TrimSpace(item)
		if ss != "" {
			res = append(res, ss)
		}
	}
	return res
}

// contains tells whether the array of string `a` contains the string `x`.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
