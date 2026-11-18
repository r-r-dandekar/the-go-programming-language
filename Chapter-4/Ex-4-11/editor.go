package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetInputFromEditor() (string, error) {
	tmpFile, err := os.CreateTemp("", "cliedit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	// Close the file handle immediately so the external editor can open/lock it safely
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		// Fallbacks based on operating system
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			editor = "nano" // or "vi" / "vim"
		}
	}

	cmd := exec.Command(editor, tmpName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}
	contentBytes, err := os.ReadFile(tmpName)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file contents: %w", err)
	}
	return strings.TrimSpace(string(contentBytes)), nil
}
