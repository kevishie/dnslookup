package dnslookup

import (
	"context"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// QueryResult is the outcome of querying one server for one name and set of types.
type QueryResult struct {
	Server    Server
	ByType    map[string][]string
	Errs      []error
	Rcode     int
	Truncated bool
}

// QueryServer performs DNS queries for fqdn against srv for each record type.
func QueryServer(ctx context.Context, client *dns.Client, srv Server, fqdn string, types []RecordType) QueryResult {
	res := QueryResult{Server: srv, ByType: make(map[string][]string)}
	endpoint := ResolverAddr(srv.Addr)

	for _, rt := range types {
		qType := uint16(rt)
		typeName := RecordType(qType).String()
		typeCtx, cancel := context.WithTimeout(ctx, client.Timeout)
		m := new(dns.Msg)
		m.SetQuestion(fqdn, qType)
		m.RecursionDesired = true
		m.SetEdns0(4096, false)

		r, _, err := client.ExchangeContext(typeCtx, m, endpoint)
		cancel()
		if err != nil {
			res.Errs = append(res.Errs, fmt.Errorf("%s: %w", typeName, err))
			continue
		}
		if r.Truncated {
			res.Truncated = true
		}
		res.Rcode = r.Rcode
		if r.Rcode != dns.RcodeSuccess {
			rc := fmt.Sprintf("RCODE%d", r.Rcode)
			if s, ok := dns.RcodeToString[r.Rcode]; ok {
				rc = s
			}
			res.Errs = append(res.Errs, fmt.Errorf("%s: %s", typeName, rc))
			continue
		}
		res.ByType[typeName] = AnswersFromMsg(r, qType)
	}
	return res
}

// FlatAnswers renders answers in stable type order for display.
func FlatAnswers(types []RecordType, q QueryResult) []string {
	var out []string
	for _, rt := range types {
		key := rt.String()
		for _, a := range q.ByType[key] {
			if len(types) > 1 {
				out = append(out, fmt.Sprintf("%s: %s", key, a))
			} else {
				out = append(out, a)
			}
		}
	}
	return out
}

const defaultTimeout = 5 * time.Second

// NewClient returns a dns.Client with UDP transport and the given timeout.
func NewClient(timeout time.Duration) *dns.Client {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &dns.Client{Net: "udp", Timeout: timeout}
}
