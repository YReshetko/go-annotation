package utils

import "strings"

func TrimQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func LastDir(s string) string {
	return s[strings.LastIndex(s, "/")+1:]
}

func Root(s string) string {
	return s[:strings.LastIndex(s, "/")]
}
