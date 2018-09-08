package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func expand(s string) string {
	if strings.HasPrefix(s, "~/") {
		return filepath.Join(os.Getenv("HOME"), s[1:])
	}

	return s
}
