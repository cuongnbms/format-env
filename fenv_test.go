package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultFunc(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		fallback string
		want     string
	}{
		{name: "missing value", value: nil, fallback: "default", want: "default"},
		{name: "empty value", value: "", fallback: "default", want: "default"},
		{name: "configured value", value: "custom", fallback: "default", want: "custom"},
		{name: "non-string value", value: 123, fallback: "default", want: "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultFunc(tt.value, tt.fallback); got != tt.want {
				t.Fatalf("defaultFunc(%v, %q) = %q, want %q", tt.value, tt.fallback, got, tt.want)
			}
		})
	}
}

func TestReadConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dev.env")
	content := "\n# comment\n KEY = value \nTOKEN=abc=123\nINVALID\nEMPTY=\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readConfigFile(path)
	if err != nil {
		t.Fatalf("readConfigFile() error = %v", err)
	}
	want := map[string]string{
		"KEY":   "value",
		"TOKEN": "abc=123",
		"EMPTY": "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readConfigFile() = %#v, want %#v", got, want)
	}
}

func TestReadConfigFileCreatesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dev.env")

	got, err := readConfigFile(path)
	if err != nil {
		t.Fatalf("readConfigFile() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("readConfigFile() = %#v, want empty config", got)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("missing config file was not created: %v", err)
	}
}

func TestFormatEnv(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "_template.env"), "STAGE={{ .STAGE }}\nPORT={{ df .PORT \"8000\" }}\n")
	writeTestFile(t, filepath.Join(dir, "dev.env"), "STAGE=dev\n")

	if err := formatEnv(dir, "dev"); err != nil {
		t.Fatalf("formatEnv() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "dev.env"))
	if err != nil {
		t.Fatal(err)
	}
	want := "STAGE=dev\nPORT=8000\n"
	if string(got) != want {
		t.Fatalf("formatted env = %q, want %q", got, want)
	}
}

func TestFormatEnvCreatesMissingStage(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "_template.env"), "PORT={{ df .PORT \"8000\" }}\n")

	if err := formatEnv(dir, "dev"); err != nil {
		t.Fatalf("formatEnv() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "dev.env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "PORT=8000\n" {
		t.Fatalf("formatted env = %q, want %q", got, "PORT=8000\n")
	}
}

func TestFormatEnvInvalidTemplateReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "_template.env"), "VALUE={{ .VALUE\n")
	configPath := filepath.Join(dir, "dev.env")
	writeTestFile(t, configPath, "ORIGINAL=value\n")

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Errorf("formatEnv() panicked instead of returning an error: %v", recovered)
		}
	}()

	if err := formatEnv(dir, "dev"); err == nil {
		t.Fatal("formatEnv() error = nil, want template parse error")
	}

	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "ORIGINAL=value\n" {
		t.Fatalf("config after parse error = %q, want original content", got)
	}
}

func TestFormatEnvRenderErrorDoesNotOverwriteExistingFile(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "_template.env")
	configPath := filepath.Join(dir, "dev.env")
	writeTestFile(t, templatePath, "VALUE={{ df .MISSING 123 }}\n")
	writeTestFile(t, configPath, "ORIGINAL=value\n")

	if err := formatEnv(dir, "dev"); err == nil {
		t.Fatal("formatEnv() error = nil, want template execution error")
	}

	got, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "ORIGINAL=value\n" {
		t.Fatalf("config after render error = %q, want original content", got)
	}
}

func TestFormatEnvPreservesExistingPermissions(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "dev.env")
	writeTestFile(t, filepath.Join(dir, "_template.env"), "STAGE={{ .STAGE }}\n")
	writeTestFile(t, configPath, "STAGE=dev\n")
	if err := os.Chmod(configPath, 0600); err != nil {
		t.Fatal(err)
	}

	if err := formatEnv(dir, "dev"); err != nil {
		t.Fatalf("formatEnv() error = %v", err)
	}

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("config permissions = %o, want 600", got)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
