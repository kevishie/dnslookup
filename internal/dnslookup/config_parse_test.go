package dnslookup

import (
	"strings"
	"testing"
)

func TestParseServerReader(t *testing.T) {
	const in = `
# comment
a 1.1.1.1
b 8.8.8.8
c 9.9.9.9:5353   # inline
`
	got, err := parseServerReader(strings.NewReader(in), "test")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("len=%d", len(got))
	}
	if got[0].Name != "a" || got[0].Addr != "1.1.1.1" {
		t.Fatalf("first: %+v", got[0])
	}
	if got[2].Addr != "9.9.9.9:5353" {
		t.Fatalf("port: %+v", got[2])
	}
}

func TestParseServerReaderBadLine(t *testing.T) {
	_, err := parseServerReader(strings.NewReader("only-one-field"), "x")
	if err == nil {
		t.Fatal("expected error")
	}
}
