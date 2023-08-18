package work

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func getProgressByMsg(msg string, max float64) (int, error) {
	parts := strings.Split(msg, "/")
	if len(parts) != 2 {
		return 0, errors.New("invalid progress format")
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return int(math.Min((float64(numerator) / float64(denominator) * max), max)), nil
}

func getProgressByCount(numerator int, denominator int, max float64) int {
	return int((100 - max) + float64(numerator)/(float64(denominator/65436))*max)
}
