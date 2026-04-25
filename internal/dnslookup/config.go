package dnslookup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Server is a display name and resolver address (IP or host; optional :port).
type Server struct {
	Name string
	Addr string
}

// DefaultServers returns built-in public resolvers (sorted by name).
func DefaultServers() []Server {
	m := map[string]string{
		"cloudflare-dns":           "1.1.1.1",
		"cloudflare-dns-secondary": "1.0.0.1",
		"google-public-dns-a":      "8.8.8.8",
		"google-public-dns-b":      "8.8.4.4",
		"opendns":                  "208.67.222.222",
		"opendns-secondary":        "208.67.220.220",
		"quad9":                    "9.9.9.9",
	}
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	out := make([]Server, 0, len(names))
	for _, name := range names {
		out = append(out, Server{Name: name, Addr: m[name]})
	}
	return out
}

// LoadServers resolves the server list: XDG config file, then ~/.resolv.conf, then defaults.
// configPath is set when servers were loaded from a non-empty config file.
func LoadServers(homeDir string, getenv func(string) string) (servers []Server, configPath string, err error) {
	if homeDir == "" {
		return nil, "", fmt.Errorf("home directory unknown")
	}
	xdg := getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		xdg = filepath.Join(homeDir, ".config")
	}
	xdgFile := filepath.Join(xdg, "dnslookup", "servers")
	if st, e := os.Stat(xdgFile); e == nil && !st.IsDir() {
		data, pe := parseServerFile(xdgFile)
		if pe != nil {
			return nil, "", pe
		}
		if len(data) > 0 {
			return data, xdgFile, nil
		}
	}
	legacy := filepath.Join(homeDir, ".resolv.conf")
	if st, e := os.Stat(legacy); e == nil && !st.IsDir() {
		data, pe := parseServerFile(legacy)
		if pe != nil {
			return nil, "", pe
		}
		if len(data) > 0 {
			return data, legacy, nil
		}
	}
	return DefaultServers(), "", nil
}
