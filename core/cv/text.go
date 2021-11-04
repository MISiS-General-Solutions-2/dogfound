package cv

import (
	"strings"
	"time"
)

var months = map[string]string{
	"янв.":  "01",
	"февр.": "02",
	"март":  "03",
	"апр.":  "04",
	"май":   "05",
	"июнь":  "06",
	"июль":  "07",
	"авг.":  "08",
	"сент.": "09",
	"окт.":  "10",
	"нояб.": "11",
	"дек.":  "12",
}

func parseRecognizedCamID(s string) string {
	idx := strings.IndexByte(s, ' ')
	if idx != -1 {
		s = s[idx+1:]
	}
	return fixCamID(s)
}

func fixCamID(s string) string {
	s = fixLetters(s)
	s = fixNumbers(s)
	return s
}
func fixLetters(s string) string {
	s = strings.Replace(s, "PVYN", "PVN", 1)
	s = strings.Replace(s, "TSAD", "TSAO", 1)
	s = strings.Replace(s, "UAD", "UAO", 1)
	s = strings.Replace(s, "SAD", "SAO", 1)
	s = strings.Replace(s, "UAOD", "UAO", 1)
	return s
}
func fixNumbers(s string) string {
	ss := strings.Split(s, "_")
	for i := range ss {
		if strings.ContainsAny(ss[i], "0123456789") {
			for j := i; j < len(ss); j++ {
				ss[j] = replaceCharsWithSimilarNumbers(ss[j])
			}
		}
	}
	return strings.Join(ss, "_")
}
func replaceCharsWithSimilarNumbers(s string) string {
	s = strings.ReplaceAll(s, "S", "5")
	return s
}
func parseRecognizedTimestamp(s string) int64 {
	if s == "" {
		return 0
	}
	s = fixO(s)
	s = rusMonthToFormat(s)
	if s == "" {
		return 0
	}
	layout := "02.01.2006:15:04:05"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0
	}
	return t.Unix()
}
func rusMonthToFormat(s string) string {
	//28 сент. 2021, 06:50:42
	parts := strings.Split(s, " ")
	if len(parts) != 4 {
		return ""
	}
	b := strings.Builder{}
	b.WriteString(parts[0])
	b.WriteRune('.')
	b.WriteString(months[parts[1]])
	b.WriteRune('.')
	b.WriteString(strings.TrimRight(parts[2], ","))
	b.WriteRune(':')
	b.WriteString(strings.Trim(parts[3], " "))

	return b.String()
}
func fixO(s string) string {
	return strings.ReplaceAll(s, "O", "0")
}
