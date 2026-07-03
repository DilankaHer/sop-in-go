package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()

	oldValue, hadValue := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("set %s: %v", key, err)
	}
	t.Cleanup(func() {
		if hadValue {
			_ = os.Setenv(key, oldValue)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	oldValue, hadValue := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if hadValue {
			_ = os.Setenv(key, oldValue)
			return
		}
		_ = os.Unsetenv(key)
	})
}

func chdir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldDir)
	})
}

func resetViper(t *testing.T) {
	t.Helper()
	viper.Reset()
}

func writeConfigFile(t *testing.T, dir, name, contents string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestGetConfigLoadsLocalEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.local.yml", "server:\n  port: \"5000\"\n")
	chdir(t, dir)
	resetViper(t)
	setEnv(t, "ENV", "local")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Server.Port != "5000" {
		t.Fatalf("Port = %q, want %q", cfg.Server.Port, "5000")
	}
}

func TestGetConfigTreatsUnsetEnvAsLocal(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.local.yml", "server:\n  port: \"4000\"\n")
	chdir(t, dir)
	resetViper(t)
	unsetEnv(t, "ENV")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Server.Port != "4000" {
		t.Fatalf("Port = %q, want %q", cfg.Server.Port, "4000")
	}
}

func TestGetConfigLoadsProdEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.prod.yml", "server:\n  port: \"9000\"\n")
	chdir(t, dir)
	resetViper(t)
	setEnv(t, "ENV", "prod")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Server.Port != "9000" {
		t.Fatalf("Port = %q, want %q", cfg.Server.Port, "9000")
	}
}

func TestGetConfigSitEnvLoadsSitEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.sit.yml", "server:\n  port: \"7000\"\n")
	chdir(t, dir)
	resetViper(t)
	setEnv(t, "ENV", "sit")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Server.Port != "7000" {
		t.Fatalf("Port = %q, want %q", cfg.Server.Port, "7000")
	}
}

func TestGetConfigUnknownEnvLoadsLocalEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.local.yml", "server:\n  port: \"6000\"\n")
	chdir(t, dir)
	resetViper(t)
	setEnv(t, "ENV", "development")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Server.Port != "6000" {
		t.Fatalf("Port = %q, want %q", cfg.Server.Port, "6000")
	}
}

func TestGetConfigReturnsErrorWhenLocalEnvFileMissing(t *testing.T) {
	chdir(t, t.TempDir())
	resetViper(t)
	setEnv(t, "ENV", "local")

	_, err := GetConfig()
	if err == nil {
		t.Fatal("GetConfig() error = nil, want error")
	}
	if !strings.Contains(err.Error(), ".env.local") {
		t.Fatalf("error = %q, want it to mention .env.local", err.Error())
	}
}

func TestGetConfigReturnsErrorWhenProdEnvFileMissing(t *testing.T) {
	chdir(t, t.TempDir())
	resetViper(t)
	setEnv(t, "ENV", "prod")

	_, err := GetConfig()
	if err == nil {
		t.Fatal("GetConfig() error = nil, want error")
	}
	if !strings.Contains(err.Error(), ".env.prod") {
		t.Fatalf("error = %q, want it to mention .env.prod", err.Error())
	}
}

func TestGetConfigValidation(t *testing.T) {
	t.Run("missing port", func(t *testing.T) {
		dir := t.TempDir()
		writeConfigFile(t, dir, ".env.local.yml", "server:\n")
		chdir(t, dir)
		resetViper(t)
		unsetEnv(t, "ENV")

		_, err := GetConfig()
		if err == nil {
			t.Fatal("GetConfig() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "missing/invalid vars: Port") {
			t.Fatalf("error = %q, want missing Port", err.Error())
		}
	})

	t.Run("empty port", func(t *testing.T) {
		dir := t.TempDir()
		writeConfigFile(t, dir, ".env.local.yml", "server:\n  port: \"\"\n")
		chdir(t, dir)
		resetViper(t)
		unsetEnv(t, "ENV")

		_, err := GetConfig()
		if err == nil {
			t.Fatal("GetConfig() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "missing/invalid vars: Port") {
			t.Fatalf("error = %q, want empty Port", err.Error())
		}
	})
}
