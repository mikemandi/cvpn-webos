package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLunaSendNotInServicePath(t *testing.T) {
	temp := t.TempDir()
	fake := filepath.Join(temp, "luna-send")
	logs := filepath.Join(temp, "calls")
	script := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \"$CVPN_LUNA_TEST_LOG\"\nprintf '{\"returnValue\":true}\\n'\n"
	if err := os.WriteFile(fake, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	old := lunaSystemPaths
	lunaSystemPaths = []string{fake}
	defer func() { lunaSystemPaths = old }()
	t.Setenv("PATH", "/nonexistent")
	t.Setenv("CVPN_LUNA_TEST_LOG", logs)
	resolved, err := lunaSendExecutable()
	if err != nil || resolved != fake {
		t.Fatalf("executable=%q err=%v", resolved, err)
	}
	if err := showNativePairingAlert("Test native alert"); err != nil {
		t.Fatal(err)
	}
	if err := systemStatusToast("Test status toast"); err != nil {
		t.Fatal(err)
	}
	calls, err := os.ReadFile(logs)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(calls), "createAlert") || !strings.Contains(string(calls), "createToast") {
		t.Fatalf("missing Luna calls: %s", calls)
	}
}

func TestLunaSendAbsentInPATHAndDefaultLocations(t *testing.T) {
	old := lunaSystemPaths
	lunaSystemPaths = []string{filepath.Join(t.TempDir(), "missing")}
	defer func() { lunaSystemPaths = old }()
	t.Setenv("PATH", "/nonexistent")
	if binary, err := lunaSendExecutable(); err == nil {
		t.Fatalf("unexpected binary %q", binary)
	}
}
