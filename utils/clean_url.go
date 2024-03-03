package utils

import "strings"

func CleanURL(url string) string {
	return strings.TrimRight(url, "/")
}
