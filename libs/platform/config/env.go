package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RuntimeMode string

const (
	RuntimeLocal  RuntimeMode = "local"
	RuntimeDocker RuntimeMode = "docker"
	RuntimeCloud  RuntimeMode = "cloud"
)

// RuntimeModeFromEnv resolves runtime mode from RUNTIME_MODE.
// If unset, it falls back to RuntimeLocal.
func RuntimeModeFromEnv() (RuntimeMode, error) {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("RUNTIME_MODE")))
	if mode == "" {
		return RuntimeLocal, nil
	}

	switch RuntimeMode(mode) {
	case RuntimeLocal, RuntimeDocker, RuntimeCloud:
		return RuntimeMode(mode), nil
	default:
		return "", fmt.Errorf("invalid RUNTIME_MODE %q (allowed: local|docker|cloud)", mode)
	}
}

// LoadModeEnvFile loads .env.<mode> from the given directory.
// Missing files are ignored.
// Existing process environment values are not overwritten.
func LoadModeEnvFile(dir string, mode RuntimeMode) error {
	filename := fmt.Sprintf(".env.%s", mode)
	path := filepath.Join(dir, filename)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("parse %s:%d: missing '='", path, lineNo)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("parse %s:%d: empty key", path, lineNo)
		}

		val = strings.TrimSpace(val)
		if len(val) >= 2 {
			if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
				(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
				val = val[1 : len(val)-1]
			}
		}

		// Respect values already injected by shell/docker/cloud runtime.
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, val); err != nil {
			return fmt.Errorf("set env %s from %s:%d: %w", key, path, lineNo, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %s: %w", path, err)
	}
	return nil
}
