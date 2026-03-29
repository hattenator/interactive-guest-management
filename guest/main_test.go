package main

import (
	"testing"
)

func TestGetLastInputTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Windows API test in short mode")
	}
	t0, err := getLastInputTime()
	if err != nil {
		t.Fatalf("getLastInputTime returned error: %v", err)
	}
	if t0 == 0 {
		t.Fatalf("getLastInputTime returned zero time, expected >0")
	}
	t.Logf("Idle Start Time is %d", t0)
}

func TestGetIdleDuration(t *testing.T) {
	idleDuration, shouldReturn, failure := getIdleDuration()
	if failure {
		t.Fatalf("getIdleDuration failed")
	}
	if shouldReturn {
		t.Fatalf("getIdleDuration threw an exception")
	}
	if idleDuration <= 0 {
		t.Fatalf("getIdleDuration should have returned positive, not %v", idleDuration)
	}
	t.Logf("Idle Duration Time is %v", idleDuration.Seconds())
}
