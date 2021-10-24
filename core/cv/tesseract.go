package cv

import (
	"strings"
	"time"

	"github.com/otiai10/gosseract"
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

func parseTimestamp(img []byte) (int64, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(img)
	client.Languages = []string{"rus"}
	text, err := client.Text()
	if err != nil {
		return 0, err
	}
	return parseRecognizedTimestamp(text), nil
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
