package work

import (
	"edetector_go/config"
	"edetector_go/pkg/file"
	"edetector_go/pkg/mariadb"
	"edetector_go/pkg/redis"
	"testing"
)

func init() {
	for i := 0; i < 2; i++ {
		file.MoveToParentDir()
	}
	vp, err := config.LoadConfig()
	if vp == nil {
		panic(err)
	}
	_, err = mariadb.Connect_init()
	if err != nil {
		panic(err)
	}
	if db := redis.Redis_init(); db == nil {
		panic(err)
	}
}

func TestParseNetowrk(t *testing.T) {
	// networkSet := make(map[string]struct{})
	// ip := "0.0.0.0"
	// test := []struct {
	// 	line       string
	// 	networkSet *map[string]struct{}
	// 	ip         string
	// 	want       []string
	// }{
	// 	{
	// 		line:       "3648|20.90.153.243:443|1701872474|1700057349|0|53751",
	// 		networkSet: &networkSet,
	// 		ip:         ip,
	// 		want:       []string{"3648", "1700057349", "1701872474", "20.90.153.243", "443", "0.0.0.0", "53751", "unknown", "in", "detect", "53751", "-", "0", "0", "20.90.153.243", "443", "GB", "0", "51", "0", "88"},
	// 	},
	// }
	// for _, tt := range test {
	// 	data := parseNetowrk(tt.line, tt.networkSet, tt.ip)
	// 	if data == nil && tt.want == nil {
	// 		continue
	// 	} else if (data == nil && tt.want != nil) || (data != nil && tt.want == nil) || len(data) != len(tt.want) {
	// 		t.Errorf("Failed: parseNetowrk(%v, %v, %v) = %v want %v", tt.line, tt.networkSet, tt.ip, data, tt.want)
	// 		continue
	// 	}
	// 	for i := range data {
	// 		if data[i] != tt.want[i] {
	// 			t.Errorf("Failed: detectNetworkElastic(%v) = %v want %v", tt.line, data, tt.want)
	// 			break
	// 		}
	// 	}
	// }
}

func TestGetProgressByMsg(t *testing.T) {
	test := []struct {
		msg  string
		max  float64
		want int
	}{
		{"1/2", 50, 25},
		{"100/1", 50, 50},
		{"1/0", 50, 0},
		{"", 0, 0},
		{"0/b", 0, 0},
	}
	for _, tt := range test {
		data, _ := getProgressByMsg(tt.msg, tt.max)
		if data != tt.want {
			t.Errorf("Failed: getProgressByMsg(%v, %v) = %v want %v", tt.msg, tt.max, data, tt.want)
		}
	}
}

func TestGetProgressByCount(t *testing.T) {
	test := []struct {
		numerator   int
		denominator int
		base        int
		max         float64
		want        int
	}{
		{1, 200, 100, 50, 25},
		{100, 200, 100, 50, 50},
		{1, 1, 0, 100, 0},
		{1, 0, 1, 100, 0},
	}
	for _, tt := range test {
		data := getProgressByCount(tt.numerator, tt.denominator, tt.base, tt.max)
		if data != tt.want {
			t.Errorf("Failed: getProgressByCount(%v, %v, %v, %v) = %v want %v", tt.numerator, tt.denominator, tt.base, tt.max, data, tt.want)
		}
	}
}

func TestGetriskscore(t *testing.T) {
	tests := []struct {
		name      string
		memory    *Memory
		wantLevel string
		wantScore string
	}{
		{
			name: "Test1",
			memory: &Memory{
				ProcessName:       "",
				ProcessCreateTime: 0,
				DynamicCommand:    "",
				ProcessMD5:        "",
				ProcessPath:       "",
				ParentProcessId:   0,
				ParentProcessName: "",
				ParentProcessPath: "",
				DigitalSign:       "",
				ProcessId:         0,
				InjectActive:      "0", // error format
				ProcessBeInjected: 0,
				Boot:              "0", // error format
				Hide:              "0", // error format
				ImportOtherDLL:    "null",
				Hook:              "",
				ProcessConnectIP:  "false",
				RiskLevel:         0,
				RiskScore:         0,
				Mode:              "",
				ProcessKey:        "",
				UUID:              "",
				Agent:             "",
				AgentIP:           "",
				AgentName:         "",
				ItemMain:          "",
				DateMain:          0,
				TypeMain:          "",
				EtcMain:           "",
			},
			wantLevel: "0",
			wantScore: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLevel, gotScore, _, _, _ := Getriskscore(*tt.memory, 0)
			if gotLevel != tt.wantLevel {
				t.Errorf("Getriskscore() gotLevel = %v, want %v", gotLevel, tt.wantLevel)
			}
			if gotScore != tt.wantScore {
				t.Errorf("Getriskscore() gotScore = %v, want %v", gotScore, tt.wantScore)
			}
		})
	}
}
