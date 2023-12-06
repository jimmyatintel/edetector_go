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
	if denominator == 0 {
		return 0, errors.New("denominator cannot be 0")
	}
	return int(math.Min((float64(numerator) / float64(denominator) * max), max)), nil
}

func getProgressByCount(numerator int, denominator int, base int, max float64) int {
	if base == 0 {
		return 0
	}
	newDenominator := float64(denominator / base)
	if newDenominator == 0 {
		return 0
	}
	return int(math.Min((float64(numerator) / float64(newDenominator) * max), max))
}
