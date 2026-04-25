package dnslookup

import (
	"context"
	"sync"

	"github.com/miekg/dns"
)

// Run queries all servers concurrently (bounded by concurrency) for name and types.
// Results are returned in the same order as servers.
func Run(ctx context.Context, client *dns.Client, servers []Server, name string, types []RecordType, concurrency int) []QueryResult {
	if concurrency < 1 {
		concurrency = 32
	}
	fqdn := dns.Fqdn(name)
	results := make([]QueryResult, len(servers))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, srv := range servers {
		wg.Add(1)
		go func(i int, srv Server) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = QueryServer(ctx, client, srv, fqdn, types)
		}(i, srv)
	}
	wg.Wait()
	return results
}
