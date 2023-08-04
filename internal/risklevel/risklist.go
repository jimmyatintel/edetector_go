package risklevel

var HighRiskMap map[string]bool

func init() {
	HighRiskMap = map[string]bool{
		"csrss.exe":    true,
		"wininit.exe":  true,
		"winlogon.exe": true,
		"explorer.exe": true,
		"smss.exe":     true,
		"services.exe": true,
		"svchost.exe":  true,
		"taskhost.exe": true,
		"lsass.exe":    true,
		"lsm.exe":      true,
		"iexplore.exe": true,
	}
}
