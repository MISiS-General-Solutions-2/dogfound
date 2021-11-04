package cv

import (
	"github.com/otiai10/gosseract"
)

func tessParseCamID(img []byte) string {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(img)
	client.Languages = []string{"eng"}
	text, err := client.Text()
	if err != nil {
		panic(err)
	}
	return parseRecognizedCamID(text)
}

func parseTimestamp(img []byte) int64 {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(img)
	client.Languages = []string{"rus"}
	text, err := client.Text()
	if err != nil {
		panic(err)
	}
	return parseRecognizedTimestamp(text)
}
