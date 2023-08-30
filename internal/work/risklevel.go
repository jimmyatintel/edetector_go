package work

import (
	"strings"
)

func Getriskscore(info Memory) (string, error) {
	score := 0
	if info.ProcessBeInjected == 2 {
		if _, ok := HighRiskMap[info.ProcessName]; ok {
			score += 150
		} else {
			score += 90
		}
	}
	if info.ProcessBeInjected == 1 {
		score += 30
	}
	if info.InjectActive[0] == '1' && info.DigitalSign == "null" {
		score += 60
	}
	if info.InjectActive[2] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if info.Boot[0] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if info.Boot[2] == '1' && info.DigitalSign == "null" {
		score += 30
	}
	if info.ProcessConnectIP != "false" {
		score += 30
	}
	if info.ImportOtherDLL != "null" {
		score += 60
	}
	if info.Hide[0] == '1' {
		score += 150
	}
	if info.Hide[2] == '1' {
		score += 60
	}
	if strings.Contains(info.Hook, "NtQuerySystemInformation") || strings.Contains(info.Hook, "RtlGetNativeSystemInformation") || strings.Contains(info.Hook, "ZwQuerySystemInformation") {
		score += 150
	}
	if info.ParentProcessPath == "null" {
		if score > 60 {
			score -= 60
		} else {
			score = 0
		}
	}
	if info.DigitalSign == "null" && info.ProcessBeInjected == 0 && info.Hook == "null" {
		score = 0
	}
	level := scoretoLevel(score)
	return level, nil
}

func scoretoLevel(score int) string {
	if score >= 150 {
		return "3"
	} else if score > 90 {
		return "2"
	} else if score > 30 {
		return "1"
	} else {
		return "0"
	}
}
