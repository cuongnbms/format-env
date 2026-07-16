package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func defaultFunc(value any, fallback string) string {
	if stringValue, ok := value.(string); ok && stringValue != "" {
		return stringValue
	}
	return fallback
}

// readConfigFile reads the key-value pairs from a .env file
func readConfigFile(filePath string) (map[string]string, error) {
	config := make(map[string]string)

	file, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// formatEnv generates a formatted .env file using a template
func formatEnv(envDir string, stage string) error {
	log.Printf("using template %s/_template.env format %s/%s.env", envDir, envDir, stage)

	configPath := filepath.Join(envDir, stage+".env")
	templatePath := filepath.Join(envDir, "_template.env")

	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	config, err := readConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	tmpl, err := template.New("env").Option("missingkey=zero").Funcs(template.FuncMap{
		"df": defaultFunc,
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(configPath, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Println("Usage: fenv <env_dir> <stages>")
		return
	}

	envDir := os.Args[1]
	stageStr := os.Args[2]
	result := strings.Split(stageStr, ",")
	for _, v := range result {
		err := formatEnv(envDir, v)
		if err != nil {
			log.Fatalf("Error formatting env: %v", err)
		}
	}
}
