package utils

import "strings"

func Sanatize(s string) string {
	return strings.TrimRight(s, "/")
}
