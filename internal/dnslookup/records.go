package dnslookup

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

// RecordType is a supported DNS query type.
type RecordType uint16

// ParseRecordTypes parses comma-separated or repeated type names into miekg types.
func ParseRecordTypes(names []string) ([]RecordType, error) {
	if len(names) == 0 {
		return []RecordType{RecordType(dns.TypeA)}, nil
	}
	var out []RecordType
	for _, raw := range names {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			t, err := parseOneType(part)
			if err != nil {
				return nil, err
			}
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return []RecordType{RecordType(dns.TypeA)}, nil
	}
	return dedupeTypes(out), nil
}

func parseOneType(s string) (RecordType, error) {
	switch strings.ToUpper(s) {
	case "A":
		return RecordType(dns.TypeA), nil
	case "AAAA":
		return RecordType(dns.TypeAAAA), nil
	case "CNAME":
		return RecordType(dns.TypeCNAME), nil
	case "MX":
		return RecordType(dns.TypeMX), nil
	case "NS":
		return RecordType(dns.TypeNS), nil
	case "TXT":
		return RecordType(dns.TypeTXT), nil
	default:
		return 0, fmt.Errorf("unknown record type %q", s)
	}
}

func dedupeTypes(in []RecordType) []RecordType {
	seen := make(map[uint16]struct{})
	var out []RecordType
	for _, t := range in {
		u := uint16(t)
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, t)
	}
	return out
}

func (t RecordType) String() string {
	switch uint16(t) {
	case dns.TypeA:
		return "A"
	case dns.TypeAAAA:
		return "AAAA"
	case dns.TypeCNAME:
		return "CNAME"
	case dns.TypeMX:
		return "MX"
	case dns.TypeNS:
		return "NS"
	case dns.TypeTXT:
		return "TXT"
	default:
		return fmt.Sprintf("TYPE%d", uint16(t))
	}
}

// AnswersFromMsg extracts string answers for the requested type from a DNS response.
func AnswersFromMsg(r *dns.Msg, qType uint16) []string {
	if r == nil {
		return nil
	}
	var out []string
	for _, ans := range r.Answer {
		switch rr := ans.(type) {
		case *dns.A:
			if qType == dns.TypeA {
				out = append(out, rr.A.String())
			}
		case *dns.AAAA:
			if qType == dns.TypeAAAA {
				out = append(out, rr.AAAA.String())
			}
		case *dns.CNAME:
			if qType == dns.TypeCNAME {
				out = append(out, strings.TrimSuffix(rr.Target, "."))
			}
		case *dns.MX:
			if qType == dns.TypeMX {
				out = append(out, fmt.Sprintf("%d %s", rr.Preference, strings.TrimSuffix(rr.Mx, ".")))
			}
		case *dns.NS:
			if qType == dns.TypeNS {
				out = append(out, strings.TrimSuffix(rr.Ns, "."))
			}
		case *dns.TXT:
			if qType == dns.TypeTXT {
				out = append(out, strings.Join(rr.Txt, " "))
			}
		}
	}
	return out
}

// ResolverAddr returns host:port for miekg Exchange.
func ResolverAddr(addr string) string {
	if addr == "" {
		return net.JoinHostPort("", "53")
	}
	if host, port, err := net.SplitHostPort(addr); err == nil {
		if port == "" {
			return net.JoinHostPort(host, "53")
		}
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort(addr, "53")
}
