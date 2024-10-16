package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestEmptyExpression(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", "\"\"", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}
