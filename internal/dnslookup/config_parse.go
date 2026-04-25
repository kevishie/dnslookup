package dnslookup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func parseServerReader(r io.Reader, label string) ([]Server, error) {
	var out []Server
	sc := bufio.NewScanner(r)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("%s:%d: expected two fields (name address)", label, lineNum)
		}
		out = append(out, Server{Name: fields[0], Addr: fields[1]})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func parseServerFile(path string) ([]Server, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parseServerReader(f, path)
}
