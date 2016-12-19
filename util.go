package main

import (
	"fmt"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("\n--\nfunction %s took %v\n", name, elapsed)
}

// StringLeft function
func StringLeft(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if L := len([]rune(s)); L < n {
		n = L
	}
	return s[:n]
}

// StringPad function
func StringPad(s string, n int, pad string) string {
	if n <= 0 {
		return ""
	}
	L := len([]rune(s))
	if L >= n {
		return s[:n]
	}
	return s + strings.Repeat(pad, n-L)
}
