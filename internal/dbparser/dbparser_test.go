package dbparser

import (
	"database/sql"
	"edetector_go/pkg/file"
	"testing"
)

func init() {
	for i := 0; i < 2; i++ {
		file.MoveToParentDir()
	}
}

func TestRFCToTimestamp(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"Sat, 25 Nov 2023 06:30:03 GMT", "1700893803"},
		{"Wed, 28 Nov 2023 23:32:04 MST", "0"},
		{"Wed, 28 2023 23:32:04", "0"},
		{"abc", "0"},
	}
	for _, tt := range tests {
		data := RFCToTimestamp(tt.value)
		if data != tt.want {
			t.Errorf("Failed: RFCToTimestamp(%v) = %v, want %v", tt.value, data, tt.want)
		}
	}
}

func TestDigitToTimestamp(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"20211227", "1640563200"},
		{"20230506", "1683331200"},
		{"2023", "0"},
		{"abc", "0"},
	}
	for _, tt := range tests {
		data := DigitToTimestamp(tt.value)
		if data != tt.want {
			t.Errorf("Failed: DigitToTimestamp(%v) = %v, want %v", tt.value, data, tt.want)
		}
	}
}

func TestGetTableNames(t *testing.T) {
	tests := []struct {
		db   string //db file path
		want []string
		err  error
	}{
		{"test/test.db", []string{"StartRun"}, nil},
		{"test/bad.db", nil, sql.ErrNoRows},
	}
	for _, tt := range tests {
		db, err := sql.Open("sqlite3", tt.db)
		if err != nil {
			t.Errorf("Error opening database file: " + err.Error())
		}
		data, err := getTableNames(db)
		if err != nil && tt.err == nil {
			t.Errorf("Unexpected error: " + err.Error())
			continue
		}
		if len(data) != len(tt.want) {
			t.Errorf("Failed: GetTableNames(%v) = %v, want %v", tt.db, data, tt.want)
			continue
		}
		for i := 0; i < len(data); i++ {
			if data[i] != tt.want[i] {
				t.Errorf("Failed: GetTableNames(%v) = %v, want %v", tt.db, data, tt.want)
			}
		}
	}
}