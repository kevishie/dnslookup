package dnslookup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadServersXDG(t *testing.T) {
	home := t.TempDir()
	xdg := filepath.Join(home, ".config", "dnslookup")
	if err := os.MkdirAll(xdg, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(xdg, "servers")
	if err := os.WriteFile(path, []byte("one 1.1.1.1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, cfg, err := LoadServers(home, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg != path {
		t.Fatalf("cfg %s", cfg)
	}
	if len(servers) != 1 || servers[0].Name != "one" {
		t.Fatalf("%+v", servers)
	}
}

func TestLoadServersLegacyFallback(t *testing.T) {
	home := t.TempDir()
	legacy := filepath.Join(home, ".resolv.conf")
	if err := os.WriteFile(legacy, []byte("legacy 8.8.8.8\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	servers, cfg, err := LoadServers(home, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg != legacy {
		t.Fatalf("cfg %s", cfg)
	}
	if len(servers) != 1 || servers[0].Name != "legacy" {
		t.Fatalf("%+v", servers)
	}
}

func TestLoadServersDefault(t *testing.T) {
	home := t.TempDir()
	servers, cfg, err := LoadServers(home, func(string) string { return "" })
	if err != nil {
		t.Fatal(err)
	}
	if cfg != "" {
		t.Fatalf("expected no cfg path")
	}
	if len(servers) < 2 {
		t.Fatalf("defaults short: %d", len(servers))
	}
}
