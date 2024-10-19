package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestSimpleEmptyExpression(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", `""`, testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	f, err := os.ReadFile(testFile)
	expected := append(f, '\n')
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestComplexEmptyExpression(t *testing.T) {
	cmd := exec.Command("sh", "-c", `go run main.go "" test.txt | diff test.txt -`)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Received output does not match expected")
	}
}

func TestOneLineExpressionMatchAll(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", "l", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	f, err := os.ReadFile(testFile)
	expected := append(f, '\n')
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestOneLineExpressionMatchOne(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", "2", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "line 2\n"
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestMatchExitCode(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", "2", testFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	if cmd.ProcessState.ExitCode() != 0 {
		t.Errorf("Expect exit code: 0\n Received %d", cmd.ProcessState.ExitCode())
	}
}

func TestNoMatch(t *testing.T) {
	testFile := "simple-test.txt"
	cmd := exec.Command("go", "run", "main.go", "NoMatch", testFile)
	output, _ := cmd.Output()

	expected := ""

	if string(output) != string(expected) {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
	if cmd.ProcessState.ExitCode() != 1 {
		t.Errorf("Expect exit code: 1\n Received %d", cmd.ProcessState.ExitCode())
	}
}

func TestSimpleRecurse(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-r", `Two`, "./test-simple-subdir")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "test-simple-subdir/dir1/file2.txt:TwoThree\ntest-simple-subdir/file1.txt:OneTwo\ntest-simple-subdir/file1.txt:TwoThree\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

// TODO: Complex recurse test

func TestSimpleInvert(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-v", `2`, "./simple-test.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "line 1\nline 3\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestDigitMatch(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", `"\d"`, "./symbols-digit-word.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "12345\n67890\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestWordCharMatch(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", `"\w"`, "./symbols-digit-word.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "12345\nabcde\n67890\nfghij\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestBeginningMatch(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", `"^AB"`, "./match-beg-end.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "ABA\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}

func TestEndMatch(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", `"AB$"`, "./match-beg-end.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %v", string(output))
		t.Fatalf("Failed to run command: %v", err)
	}

	expected := "BAB\n"

	if string(output) != expected {
		t.Errorf("Expected: %q\nReceived: %q", expected, string(output))
	}
}
