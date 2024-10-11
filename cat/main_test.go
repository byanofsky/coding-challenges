package main

import (
	"os"
	"os/exec"
	"testing"
)

// Test for basic my-cat functionality
func TestBasic(t *testing.T) {
	testFile := "./test.txt"
	cmd := exec.Command("go", "run", "main.go", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	if string(output) != "Hi" {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// TODO: Add test for errors
