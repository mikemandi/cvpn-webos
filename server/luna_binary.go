package main

import (
	"fmt"
	"os"
	"os/exec"
)

// The background panel may be started by a Homebrew boot hook with a minimal
// PATH, unlike an interactive root SSH shell. Resolve the system tool by a
// trusted absolute location instead of relying on the inherited PATH.
// These locations are read-only system directories on ordinary webOS images.
var lunaSystemPaths = []string{
	"/usr/bin/luna-send",
	"/usr/sbin/luna-send",
	"/bin/luna-send",
	"/sbin/luna-send",
	"/usr/local/bin/luna-send",
	"/usr/local/sbin/luna-send",
}

func lunaSendExecutable() (string, error) {
	for _, candidate := range lunaSystemPaths {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return candidate, nil
		}
	}
	// Compatibility for nonstandard webOS installs. LookPath may return an
	// absolute path or a local relative path if PATH was explicitly configured.
	// The standard paths above always take precedence.
	if candidate, err := exec.LookPath("luna-send"); err == nil {
		return candidate, nil
	}
	return "", fmt.Errorf("luna-send not found in system directories or panel PATH")
}
