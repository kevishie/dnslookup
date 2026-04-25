package dnslookup

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestParseRecordTypesDefault(t *testing.T) {
	types, err := ParseRecordTypes(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 1 || types[0] != RecordType(dns.TypeA) {
		t.Fatalf("%v", types)
	}
}

func TestParseRecordTypesDedupe(t *testing.T) {
	types, err := ParseRecordTypes([]string{"A", "AAAA", "A"})
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 2 {
		t.Fatalf("%v", types)
	}
}

func TestParseRecordTypesComma(t *testing.T) {
	types, err := ParseRecordTypes([]string{"A, AAAA", "MX"})
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 3 {
		t.Fatalf("%v", types)
	}
}

func TestAnswersFromMsgA(t *testing.T) {
	m := new(dns.Msg)
	m.Answer = append(m.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: "x.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1},
		A:   net.IPv4(1, 2, 3, 4),
	})
	got := AnswersFromMsg(m, dns.TypeA)
	if len(got) != 1 || got[0] != "1.2.3.4" {
		t.Fatalf("%v", got)
	}
}

func TestResolverAddr(t *testing.T) {
	if ResolverAddr("8.8.8.8") != "8.8.8.8:53" {
		t.Fatal(ResolverAddr("8.8.8.8"))
	}
	if ResolverAddr("127.0.0.1:5353") != "127.0.0.1:5353" {
		t.Fatal(ResolverAddr("127.0.0.1:5353"))
	}
}
