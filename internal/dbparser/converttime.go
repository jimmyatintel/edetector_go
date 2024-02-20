package dbparser

import (
	"edetector_go/pkg/logger"
	"fmt"
	"strings"
	"time"
)

func RFCToTimestamp(original string) string {
	original = strings.TrimSpace(original)
	if original == "0" || original == "-1" {
		return "0"
	}
	layout := "Mon, 02 Jan 2006 15:04:05 GMT"
	t, err := time.Parse(layout, original)
	if err != nil {
		logger.Error("Error parsing time: " + err.Error())
		return "0"
	}
	return fmt.Sprintf("%d", t.Unix())
}

func DigitToTimestamp(original string) string {
	original = strings.TrimSpace(original)
	if original == "0" || original == "-1" {
		return "0"
	}
	original = original + "000000"
	layout := "20060102150405"
	t, err := time.Parse(layout, original)
	if err != nil {
		logger.Error("Error parsing time: " + err.Error())
		return "0"
	}
	return fmt.Sprintf("%d", t.Unix())
}
