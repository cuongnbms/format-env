package main

import (
	"os"
	"strings"
	"testing"
)

func TestPreCommitHookMetadata(t *testing.T) {
	content, err := os.ReadFile(".pre-commit-hooks.yaml")
	if err != nil {
		t.Fatalf("read hook metadata: %v", err)
	}

	metadata := string(content)
	for _, expected := range []string{
		"- id: format-env",
		"entry: format-env",
		"language: golang",
		"pass_filenames: false",
		`files: '\.env$'`,
	} {
		if !strings.Contains(metadata, expected) {
			t.Errorf("hook metadata does not contain %q", expected)
		}
	}
}
