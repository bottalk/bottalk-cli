package main

import "testing"

func TestCli(t *testing.T) {
	testCommands := getTestCommands()
	if len(testCommands) != 1 {
		t.Errorf("Test commands count expected to be 1, but was %v", len(testCommands))
	}
}
