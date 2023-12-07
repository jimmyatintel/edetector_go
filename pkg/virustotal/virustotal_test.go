package virustotal

import (
	"testing"
)

func TestScanIP(t *testing.T) {
	test := []struct {
		ip     string
		apikey string
		want1  int
		want2  int
	}{
		{ip: "0.0.0.0", apikey: "", want1: 0, want2: 0},
	}
	for _, tt := range test {
		got1, got2, _ := ScanIP(tt.ip, tt.apikey)
		if got1 != tt.want1 || got2 != tt.want2{
			t.Errorf("ScanIP(%v) = (%v, %v), want = (%v, %v)", tt.ip, got1, got2, tt.want1, tt.want2)
		}
	}
}
