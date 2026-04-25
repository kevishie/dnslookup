# dnslookup

CLI to query many DNS resolvers at once and compare answers—useful for spotting propagation differences across public DNS and your own server list.

## Prerequisites

Install a recent Go toolchain (this project targets Go 1.22+). On macOS with Homebrew:

```bash
brew install go
```

From the repository root, download module dependencies and refresh `go.sum`:

```bash
go mod tidy
```

## Install

From a clone of this repository (after `go mod tidy` above):

```bash
go install ./cmd/dnslookup
```

The installed binary is named `dnslookup`. Go places it in **`$(go env GOBIN)`** when that is non-empty; otherwise **`$(go env GOPATH)/bin`** (often **`~/go/bin`**). That directory is usually **not** on your `PATH`, which leads to `command not found`.

Add the install directory to `PATH` for the current shell (zsh):

```bash
d="$(go env GOBIN)"
[ -n "$d" ] || d="$(go env GOPATH)/bin"
export PATH="$d:$PATH"
```

To make that permanent, add the same three lines to `~/.zshrc`, then run `source ~/.zshrc` or open a new terminal.

Check the binary exists:

```bash
d="$(go env GOBIN)"; [ -n "$d" ] || d="$(go env GOPATH)/bin"; ls -la "$d/dnslookup"
```

Without installing, you can run from the repo with:

```bash
go run ./cmd/dnslookup -- draftkings.com
```

(`--` separates `go run` flags from the program’s arguments.)

## Usage

```bash
dnslookup example.com
dnslookup -t A -t AAAA example.com
dnslookup -json example.com
dnslookup -strict -t A example.com
```

### Flags

| Flag | Default | Meaning |
|------|---------|---------|
| `-t` | `A` | Record type (repeatable): `A`, `AAAA`, `CNAME`, `MX`, `NS`, `TXT`. Comma-separated values in one `-t` are allowed (e.g. `-t A,AAAA`). |
| `-timeout` | `5s` | Timeout for each DNS exchange (per type, per server). |
| `-c` | `32` | Maximum concurrent server queries. |
| `-json` | off | Print JSON (schema version `1`) to stdout. |
| `-strict` | off | Exit `2` if any resolver fails (partial failures count). |
| `-no-color` | off | Disable ANSI colors in the table. |

### Exit codes

- `0` — At least one resolver returned data (not all failed).
- `1` — Usage error or invalid flags/config.
- `2` — Every resolver failed, or `-strict` and any resolver failed.

## Configuration

Resolver lists use two whitespace-separated columns: `name` and `address`. Lines starting with `#` and blank lines are ignored. Optional `:port` is supported (default `53`).

**Preferred:** `$XDG_CONFIG_HOME/dnslookup/servers` (when `XDG_CONFIG_HOME` is unset, `~/.config/dnslookup/servers`).

**Legacy:** `~/.resolv.conf` in your home directory (same format as above).

If neither file exists or both are empty, built-in public resolvers are used (Google, Cloudflare, Quad9, OpenDNS).

Example `servers` file:

```text
google-public-dns-a.google.com 8.8.8.8
google-public-dns-b.google.com 8.8.4.4
office-resolver 10.0.0.53
```

When a config file is used, the path is printed on stderr before results (same idea as the original tool).

## Example output

```text
Reading servers from /home/you/.config/dnslookup/servers

SERVER                        ADDRESS         RESULT
cloudflare-dns                1.1.1.1         93.184.216.34
google-public-dns-a           8.8.8.8         93.184.216.34
...
```
