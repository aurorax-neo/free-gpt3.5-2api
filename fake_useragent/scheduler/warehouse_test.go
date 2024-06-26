package scheduler

import "testing"

func TestAppendUrl(t *testing.T) {
	AppendUrl("https://www.google.com/")
	AppendUrl("https://golang.org/")
}

func TestPopUrl(t *testing.T) {
	if url := PopUrl(); url == "" {
		t.Errorf("Expected value, but empty")
	}
}

func TestCountUrl(t *testing.T) {
	if count := CountUrl(); count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}
