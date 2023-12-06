package ip2location

import (
	"edetector_go/config"
	"edetector_go/pkg/file"
	"testing"
)

func init() {
	for i := 0; i < 2; i++ {
		file.MoveToParentDir()
	}
	vp, _ := config.LoadConfig()
	if vp == nil {
		panic("Error loading config file")
	}
}

func TestToCountry(t *testing.T) {
	tests := []struct {
		ip   string
		want string
	}{
		{"13.107.5.91", "US"},
		{"0.0.0.0", "-"},
		{"123", "invalid IP"},
		{"", "invalid IP"},
	}

	for _, tt := range tests {
		data, _ := ToCountry(tt.ip)
		if data != tt.want {
			t.Errorf("Failed: ToCountry(%v) = %v, want %v", tt.ip, data, tt.want)
		}
	}
}

func TestToLatitudeLongtitude(t *testing.T) {
	tests := []struct {
		ip     string
		wantLo int
		wantLa int
	}{
		{"13.107.5.91", -122, 47},
		{"0.0.0.0", 0, 0},
		{"123", 0, 0},
		{"", 0, 0},
	}

	for _, tt := range tests {
		lo, la, _ := ToLatitudeLongtitude(tt.ip)
		if lo != tt.wantLo || la != tt.wantLa {
			t.Errorf("Failed: ToLatitudeLongtitude(%v) = %v, %v, want %v, %v", tt.ip, lo, la, tt.wantLo, tt.wantLa)
		}
	}
}
