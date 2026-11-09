package Handlers

import "strings"

func UnMask(f string) string {
	f = strings.ReplaceAll(f, "-", "")
	f = strings.ReplaceAll(f, ".", "")
	f = strings.ReplaceAll(f, "/", "")
	f = strings.ReplaceAll(f, "_", "")
	f = strings.ReplaceAll(f, "(", "")
	f = strings.ReplaceAll(f, ")", "")
	f = strings.ReplaceAll(f, " ", "")
	return f
}