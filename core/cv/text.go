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
	"мая":   "05",
	"июня":  "06",
	"июля":  "07",
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
	s = strings.ReplaceAll(s, "O", "0")
	s = strings.ReplaceAll(s, "D", "0")
	s = strings.ReplaceAll(s, "o", "0")
	s = strings.ReplaceAll(s, "О", "0")
	s = strings.ReplaceAll(s, "о", "0")
	s = strings.ReplaceAll(s, "б", "6")
	return s
}
func parseRecognizedTimestamp(s string) int64 {
	if s == "" {
		return 0
	}
	if idx := strings.IndexByte(s, '\n'); idx != -1 {
		s = s[:idx]
	}
	s = rusMonthToFormat(s)
	if s == "" {
		return 0
	}
	//28 сент. 2021
	layout := "02.01.2006"
	t, err := time.Parse(layout, s)
	if err != nil {
		return 0
	}
	return t.Unix()
}
func rusMonthToFormat(s string) string {
	parts := strings.Split(s, " ")
	if len(parts) != 4 {
		return ""
	}
	b := strings.Builder{}
	b.WriteString(replaceCharsWithSimilarNumbers(parts[0]))
	b.WriteRune('.')
	b.WriteString(months[parts[1]])
	b.WriteRune('.')
	b.WriteString(strings.TrimRight(replaceCharsWithSimilarNumbers(parts[2]), ","))

	return b.String()
}
