package channelmap

import (
	"errors"
	"testing"
)

func TestAssignDiskChannel(t *testing.T) {
	channel := make(chan string)
	tests := []struct {
		key       string
		task_chan chan string
	}{
		{
			key:       "test",
			task_chan: channel,
		},
	}
	for ind, tt := range tests {
		AssignDiskChannel(tt.key, &tt.task_chan)
		if *UserDiskChannel[tt.key] != tt.task_chan {
			t.Errorf("Failed TestCase %v: UserDiskChannel", ind)
		}
	}
}

func TestGetDiskChannel(t *testing.T) {
	channel := make(chan string)
	tests := []struct {
		key       string
		task_chan chan string
		err       error
	}{
		{key: "no", task_chan: channel, err: errors.New("invalid key")},
	}
	for ind, tt := range tests {
		channel, err := GetDiskChannel(tt.key)
		if err != nil {
			if tt.err == nil {
				t.Errorf("Failed TestCase %v: GetDiskChannel", ind)
			}
		} else {
			if channel != tt.task_chan {
				t.Errorf("Failed TestCase %v: GetDiskChannel", ind)
			}
		}
	}
}
