package main

import (
	"fmt"
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
	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// Test for stdin as input
func TestStdin(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'hello world' | go run main.go -")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "hello world\n"
	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// Test concatentation of multiple files
func TestConcat(t *testing.T) {
	testFile1 := "./test.txt"
	testFile2 := "./test2.txt"
	cmd := exec.Command("go", "run", "main.go", testFile1, testFile2)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected1, err := os.ReadFile(testFile1)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	expected2, err := os.ReadFile(testFile2)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}
	expected := fmt.Sprintf("%s%s", expected1, expected2)

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// TODO: Add test for errors
