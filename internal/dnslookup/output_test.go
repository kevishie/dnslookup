package dnslookup

import (
	"errors"
	"testing"
)

func TestExitStatus(t *testing.T) {
	ok := []QueryResult{{ByType: map[string][]string{"A": {"1.1.1.1"}}}}
	if n := ExitStatus(ok, false); n != 0 {
		t.Fatalf("got %d", n)
	}
	if n := ExitStatus(ok, true); n != 0 {
		t.Fatalf("strict ok got %d", n)
	}

	fail := []QueryResult{{Errs: []error{errors.New("x")}}}
	if n := ExitStatus(fail, false); n != 2 {
		t.Fatalf("all fail got %d", n)
	}

	mixed := []QueryResult{
		{ByType: map[string][]string{"A": {"1.1.1.1"}}},
		{Errs: []error{errors.New("timeout")}},
	}
	if n := ExitStatus(mixed, false); n != 0 {
		t.Fatalf("partial success got %d", n)
	}
	if n := ExitStatus(mixed, true); n != 2 {
		t.Fatalf("strict partial got %d", n)
	}
}
