package main

import (
	"bytes"
	"strings"
	"testing"
)

func captureStdout(f func() int) (string, int) {
	old := stdout
	buf := new(bytes.Buffer)
	stdout = buf
	exitCode := f()
	stdout = old
	return buf.String(), exitCode
}

func captureStderr(f func() int) (string, int) {
	old := stderr
	buf := new(bytes.Buffer)
	stderr = buf
	exitCode := f()
	stderr = old
	return buf.String(), exitCode
}

func TestRunParse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"valid semver", []string{"1.2.3"}, "1.2.3\n", 0},
		{"no args", []string{}, "", 2},
		{"too many args", []string{"1.2.3", "extra"}, "", 2},
		{"empty version", []string{""}, "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code := captureStdout(func() int {
				return runParse(tt.args)
			})
			if code != tt.wantCode {
				t.Errorf("runParse() code = %d, want %d", code, tt.wantCode)
			}
			if out != tt.wantOut {
				t.Errorf("runParse() output = %q, want %q", out, tt.wantOut)
			}
		})
	}
}

func TestRunCompare(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"less than", []string{"1.2.3", "1.2.4"}, "1.2.3 < 1.2.4\n", 0},
		{"greater than", []string{"1.2.4", "1.2.3"}, "1.2.4 > 1.2.3\n", 0},
		{"equal", []string{"1.2.3", "1.2.3"}, "1.2.3 == 1.2.3\n", 0},
		{"no args", []string{}, "", 2},
		{"one arg", []string{"1.2.3"}, "", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code := captureStdout(func() int {
				return runCompare(tt.args)
			})
			if code != tt.wantCode {
				t.Errorf("runCompare() code = %d, want %d", code, tt.wantCode)
			}
			if out != tt.wantOut {
				t.Errorf("runCompare() output = %q, want %q", out, tt.wantOut)
			}
		})
	}
}

func TestRunCheck(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"satisfies", []string{"1.2.3", ">=1.2.0"}, "1.2.3 satisfies >=1.2.0\n", 0},
		{"satisfies caret", []string{"1.2.3", "^1.2.0"}, "1.2.3 satisfies ^1.2.0\n", 0},
		{"does not satisfy", []string{"1.1.0", ">=1.2.0"}, "1.1.0 does not satisfy >=1.2.0\n", 1},
		{"no args", []string{}, "", 2},
		{"one arg", []string{"1.2.3"}, "", 2},
		{"invalid version", []string{"", ">=1.0.0"}, "", 1},
		{"invalid constraint", []string{"1.2.3", ""}, "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code := captureStdout(func() int {
				return runCheck(tt.args)
			})
			if code != tt.wantCode {
				t.Errorf("runCheck() code = %d, want %d", code, tt.wantCode)
			}
			if out != tt.wantOut {
				t.Errorf("runCheck() output = %q, want %q", out, tt.wantOut)
			}
		})
	}
}

func TestRunDetect(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantCode int
	}{
		{"semver", []string{"1.2.3"}, "semver\n", 0},
		{"no args", []string{}, "", 2},
		{"undetectable", []string{"not a version"}, "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code := captureStdout(func() int {
				return runDetect(tt.args)
			})
			if code != tt.wantCode {
				t.Errorf("runDetect() code = %d, want %d", code, tt.wantCode)
			}
			if out != tt.wantOut {
				t.Errorf("runDetect() output = %q, want %q", out, tt.wantOut)
			}
		})
	}
}

func TestRunUnknownCommand(t *testing.T) {
	errOut, code := captureStderr(func() int {
		return run([]string{"bogus"})
	})
	if code != 2 {
		t.Errorf("run() code = %d, want 2", code)
	}
	if !strings.Contains(errOut, "unknown command") {
		t.Errorf("expected 'unknown command' in stderr, got %q", errOut)
	}
}

func TestRunNoArgs(t *testing.T) {
	out, code := captureStdout(func() int {
		return run([]string{})
	})
	if code != 2 {
		t.Errorf("run() code = %d, want 2", code)
	}
	if !strings.Contains(out, "Usage:") {
		t.Errorf("expected 'Usage:' in output, got %q", out)
	}
}

func TestRunHelp(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"-h"}, {"--help"}} {
		out, code := captureStdout(func() int {
			return run(args)
		})
		if code != 0 {
			t.Errorf("run(%v) code = %d, want 0", args, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("expected 'Usage:' in output, got %q", out)
		}
	}
}
