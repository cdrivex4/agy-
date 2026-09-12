package agy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Installation represents a discovered Antigravity CLI installation.
type Installation struct {
	BinaryPath   string
	Version      string
	Capabilities *Capabilities
}

// Discover locates the installed agy binary and inspects its version and capabilities.
func Discover() (*Installation, error) {
	binPath := os.Getenv("AGY_BINARY_PATH")
	if binPath == "" {
		var err error
		binPath, err = exec.LookPath(binaryName())
		if err != nil {
			binPath = probeStandardLocations()
		}
	}

	if binPath == "" {
		return nil, errors.New("antigravity CLI (agy) executable not found. Ensure agy is installed and in your PATH")
	}

	// Probe version
	versionCmd := exec.Command(binPath, "--version")
	var verOut bytes.Buffer
	versionCmd.Stdout = &verOut
	versionCmd.Stderr = &verOut
	if err := versionCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to query version from %s: %w", binPath, err)
	}
	versionStr := strings.TrimSpace(verOut.String())

	// Probe capabilities from help
	helpCmd := exec.Command(binPath, "--help")
	var helpOut bytes.Buffer
	helpCmd.Stdout = &helpOut
	helpCmd.Stderr = &helpOut
	_ = helpCmd.Run() // Help output can exit with 0 or non-zero depending on flag parsing

	caps := ParseCapabilities(helpOut.String())

	return &Installation{
		BinaryPath:   binPath,
		Version:      versionStr,
		Capabilities: caps,
	}, nil
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "agy.exe"
	}
	return "agy"
}

func probeStandardLocations() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	var candidates []string
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			candidates = append(candidates, filepath.Join(localAppData, "agy", "bin", "agy.exe"))
		}
		candidates = append(candidates,
			filepath.Join(home, ".gemini", "antigravity-cli", "bin", "agy.exe"),
			filepath.Join(home, ".local", "bin", "agy.exe"),
		)
	} else {
		candidates = append(candidates,
			"/usr/local/bin/agy",
			filepath.Join(home, ".local", "bin", "agy"),
			filepath.Join(home, ".gemini", "antigravity-cli", "bin", "agy"),
		)
	}

	for _, path := range candidates {
		if fi, err := os.Stat(path); err == nil && !fi.IsDir() {
			return path
		}
	}

	return ""
}
