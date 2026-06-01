package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func writeEnvFile(t *testing.T, dir, name, contents string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func setValidPort(t *testing.T) {
	t.Helper()
	setEnv(t, "PORT", "8080")
}

func TestGetConfigLoadsLocalEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.local", "PORT=5000\n")
	chdir(t, dir)
	setEnv(t, "ENV", "local")
	unsetEnv(t, "PORT")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Port != "5000" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "5000")
	}
}

func TestGetConfigTreatsUnsetEnvAsLocal(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.local", "PORT=4000\n")
	chdir(t, dir)
	unsetEnv(t, "ENV")
	unsetEnv(t, "PORT")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Port != "4000" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "4000")
	}
}

func TestGetConfigLoadsProdEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.prod", "PORT=9000\n")
	chdir(t, dir)
	setEnv(t, "ENV", "prod")
	unsetEnv(t, "PORT")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Port != "9000" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "9000")
	}
}

func TestGetConfigSitEnvLoadsLocalEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.sit", "PORT=7000\n")
	chdir(t, dir)
	setEnv(t, "ENV", "sit")
	unsetEnv(t, "PORT")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Port != "7000" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "7000")
	}
}

func TestGetConfigUnknownEnvLoadsLocalEnvFile(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.local", "PORT=6000\n")
	chdir(t, dir)
	setEnv(t, "ENV", "development")
	unsetEnv(t, "PORT")

	cfg, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg.Port != "6000" {
		t.Fatalf("Port = %q, want %q", cfg.Port, "6000")
	}
}

func TestGetConfigReturnsErrorWhenLocalEnvFileMissing(t *testing.T) {
	chdir(t, t.TempDir())
	setEnv(t, "ENV", "local")
	unsetEnv(t, "PORT")

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
	setEnv(t, "ENV", "prod")
	unsetEnv(t, "PORT")

	_, err := GetConfig()
	if err == nil {
		t.Fatal("GetConfig() error = nil, want error")
	}
	if !strings.Contains(err.Error(), ".env.prod") {
		t.Fatalf("error = %q, want it to mention .env.prod", err.Error())
	}
}

func TestValidateVars(t *testing.T) {
	t.Run("valid port in range", func(t *testing.T) {
		setValidPort(t)

		if err := validateVars(); err != nil {
			t.Fatalf("validateVars() error = %v, want nil", err)
		}
	})

	t.Run("accepts upper bound port", func(t *testing.T) {
		setEnv(t, "PORT", "65535")

		if err := validateVars(); err != nil {
			t.Fatalf("validateVars() error = %v, want nil", err)
		}
	})

	t.Run("unset port", func(t *testing.T) {
		unsetEnv(t, "PORT")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "unset variables => PORT") {
			t.Fatalf("error = %q, want unset PORT", err.Error())
		}
	})

	t.Run("empty port", func(t *testing.T) {
		setEnv(t, "PORT", "")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "empty variables => PORT") {
			t.Fatalf("error = %q, want empty PORT", err.Error())
		}
	})

	t.Run("non numeric port", func(t *testing.T) {
		setEnv(t, "PORT", "abc")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "invalid variables => PORT: abc is not a valid port") {
			t.Fatalf("error = %q, want invalid PORT", err.Error())
		}
	})

	t.Run("port below range", func(t *testing.T) {
		setEnv(t, "PORT", "2999")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "invalid variables => PORT: 2999 is not a valid port") {
			t.Fatalf("error = %q, want invalid PORT", err.Error())
		}
	})

	t.Run("port above max", func(t *testing.T) {
		setEnv(t, "PORT", "65536")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		if !strings.Contains(err.Error(), "invalid variables => PORT: 65536 is not a valid port") {
			t.Fatalf("error = %q, want invalid PORT", err.Error())
		}
	})

	t.Run("aggregates unset and invalid port errors", func(t *testing.T) {
		unsetEnv(t, "PORT")

		err := validateVars()
		if err == nil {
			t.Fatal("validateVars() error = nil, want error")
		}
		errMsg := err.Error()
		if !strings.Contains(errMsg, "unset variables => PORT") {
			t.Fatalf("error = %q, want unset PORT", errMsg)
		}
		if !strings.Contains(errMsg, "invalid variables => PORT:") {
			t.Fatalf("error = %q, want invalid PORT", errMsg)
		}
	})
}
