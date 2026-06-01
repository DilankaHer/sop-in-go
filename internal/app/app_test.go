package app

import "testing"

func TestNewAppCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	writeEnvFile(t, dir, ".env.local", "PORT=8080\n")
	chdir(t, dir)
	unsetEnv(t, "ENV")
	unsetEnv(t, "PORT")

	application, err := NewApp()
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	if application.Config == nil {
		t.Fatal("Config is nil")
	}
	if application.Config.Port != "8080" {
		t.Fatalf("Port = %q, want %q", application.Config.Port, "8080")
	}
}

func TestNewAppReturnsConfigLoadError(t *testing.T) {
	chdir(t, t.TempDir())
	unsetEnv(t, "ENV")
	unsetEnv(t, "PORT")

	_, err := NewApp()
	if err == nil {
		t.Fatal("NewApp() error = nil, want error")
	}
}
