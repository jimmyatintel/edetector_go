package parsedb

var dbMap map[string]func(agent string, lines []string) error

func init() {
	dbMap = map[string]func(agent string, lines []string) error {
		"ARPCache": ARPCacheTable,
	}
}
