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

	expected := "hello world"
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

// Test n flag functionality
func TestNFlag(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'line 1\nline 2\nline 3' | go run main.go -n -")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "1  line 1\n2  line 2\n3  line 3"
	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// Test numbering blank lines
func TestNumberBlankLines(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'line 1\n\nline 2\n\nline 3' | go run main.go -n -")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "1  line 1\n2  \n3  line 2\n4  \n5  line 3"
	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// Test numbering blank lines
func TestNumberSkipBlankLines(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo 'line 1\n\nline 2\n\nline 3' | go run main.go -b -")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "1  line 1\n\n2  line 2\n\n3  line 3"
	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// TODO: Add test for errors
