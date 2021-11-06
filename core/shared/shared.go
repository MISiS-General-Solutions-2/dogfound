package shared

import "strings"

func GetExtension(s string) string {
	idx := strings.LastIndexByte(s, '.')
	if idx == -1 {
		return ""
	}
	return s[idx:]
}
func ChangeExtension(s, newExt string) string {
	idx := strings.LastIndexByte(s, '.')
	if idx == -1 {
		return s + newExt
	}
	return s[:idx] + newExt
}
