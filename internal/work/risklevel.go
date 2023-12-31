package work

import (
	"edetector_go/pkg/logger"
	"edetector_go/pkg/mariadb/query"
	"strconv"
	"strings"
)

func Getriskscore(info Memory) (string, string, error) {
	score := 0
	realPath := strings.Replace(info.ProcessPath, "\\\\", "\\", -1)
	// white list
	whiteList, err := query.Load_white_list()
	if err != nil {
		logger.Error("Error loading white list" + err.Error())
	} else {
		for _, white := range whiteList {
			if (white[0] == "" || info.ProcessName == white[0]) && (white[1] == "" || info.ProcessMD5 == white[1]) && strings.Contains(info.DigitalSign, white[2]) && strings.Contains(realPath, white[3]) {
				logger.Debug("Hit white list")
				return "0", "0", nil
			}
		}
	}
	// black list
	blackList, err := query.Load_black_list()
	if err != nil {
		logger.Error("Error loading black list" + err.Error())
	} else {
		for _, black := range blackList {
			if (black[0] == "" || info.ProcessName == black[0]) && (black[1] == "" || info.ProcessMD5 == black[1]) && strings.Contains(info.DigitalSign, black[2]) && strings.Contains(realPath, black[3]) {
				logger.Debug("Hit black list")
				return "3", "150", nil
			}
		}
	}
	// hack list
	hackList, err := query.Load_hack_list()
	if err != nil {
		logger.Error("Error loading hack list" + err.Error())
	} else {
		for _, hack := range hackList {
			if (hack[0] == "" || info.ProcessName == hack[0]) && strings.Contains(info.DynamicCommand, hack[1]) && strings.Contains(realPath, hack[2]) {
				point, err := strconv.Atoi(hack[3])
				if err != nil {
					logger.Error("Error converting adding_point to integer" + err.Error())
					continue
				}
				logger.Debug("Hit hack list")
				score += point
			}
		}
	}

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
	return level, strconv.Itoa(score), nil
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
