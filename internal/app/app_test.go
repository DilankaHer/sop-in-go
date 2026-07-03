package app

import "testing"

func TestNewAppCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	writeConfigFile(t, dir, ".env.local.yml", "server:\n  port: \"8080\"\n")
	chdir(t, dir)
	resetViper(t)
	unsetEnv(t, "ENV")

	application, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	if application.Config == nil {
		t.Fatal("Config is nil")
	}
	if application.Config.Server.Port != "8080" {
		t.Fatalf("Port = %q, want %q", application.Config.Server.Port, "8080")
	}
}

func TestNewAppReturnsConfigLoadError(t *testing.T) {
	chdir(t, t.TempDir())
	resetViper(t)
	unsetEnv(t, "ENV")

	_, err := NewApp()
	if err == nil {
		t.Fatal("NewApp() error = nil, want error")
	}
}
