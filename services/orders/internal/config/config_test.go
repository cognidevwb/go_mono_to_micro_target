package config

import "testing"

func TestLoadFallsBackToLocalDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	cfg := Load()
	if cfg.Port != "8080" || cfg.Service == "" || cfg.MaxConns <= 0 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadReadsTheEnvironment(t *testing.T) {
	t.Setenv("PORT", "9999")
	if got := Load().Port; got != "9999" {
		t.Fatalf("PORT = %q", got)
	}
}
