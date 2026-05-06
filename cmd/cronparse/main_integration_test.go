package main

import (
	"os/exec"
	"strings"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	tmpBin := t.TempDir() + "/cronparse"
	cmd := exec.Command("go", "build", "-o", tmpBin, ".")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return tmpBin
}

func TestCLI_ValidExpression(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "* * * * *").CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	outStr := string(out)
	if !strings.Contains(outStr, "Expression") {
		t.Errorf("expected 'Expression' in output, got:\n%s", outStr)
	}
	if !strings.Contains(outStr, "Description") {
		t.Errorf("expected 'Description' in output, got:\n%s", outStr)
	}
	if !strings.Contains(outStr, "Next 5 runs") {
		t.Errorf("expected 'Next 5 runs' in output, got:\n%s", outStr)
	}
}

func TestCLI_InvalidExpression(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "invalid expression")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for invalid expression, got output:\n%s", out)
	}
	if !strings.Contains(string(out), "Error") {
		t.Errorf("expected 'Error' in stderr output, got:\n%s", out)
	}
}

func TestCLI_CustomN(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "-n", "3", "0 12 * * *").CombinedOutput()
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	outStr := string(out)
	if !strings.Contains(outStr, "Next 3 runs") {
		t.Errorf("expected 'Next 3 runs' in output, got:\n%s", outStr)
	}
}

func TestCLI_NoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit when no args given, got output:\n%s", out)
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage info in output, got:\n%s", out)
	}
}
