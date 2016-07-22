package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"

	"github.com/miekg/dns"
)

func main() {

	if len(os.Args) > 1 {
		target := os.Args[1]
		dnsServers := CustomServerList()

		if len(dnsServers) == 0 {
			dnsServers = DefaultServerList()
		}

		for alias, address := range dnsServers {
			result := Lookup(target, address)
			fmt.Printf("%v ==> %v\n", alias, result)
		}

	} else {
		Help()
	}
}

// Help prints info on how to use this program.
func Help() {
	fmt.Println("Please provide a target.")
}

// Lookup host (target) per DNS server.
func Lookup(target string, server string) []net.IP {
	result := []net.IP{}
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)

	r, _, err := c.Exchange(&m, server+":53")
	if err == nil {
		for _, ans := range r.Answer {
			record, isType := ans.(*dns.A)
			if isType == true {
				result = append(result, record.A)
			}
		}
	}

	return result
}

// DefaultServerList returns map of Google public DNS.
func DefaultServerList() map[string]string {
	servers := make(map[string]string)

	servers["google-public-dns-a.google.com"] = "8.8.8.8"
	servers["google-public-dns-b.google.com"] = "8.8.4.4"

	return servers
}

// CustomServerList looks for .resolv.conf file in user home directory
func CustomServerList() map[string]string {

	result := make(map[string]string)
	usr, err := user.Current()
	checkError(err)

	filename := usr.HomeDir + "/.resolv.conf"

	if filename != "" {
		if file, err := os.Open(filename); err == nil {

			fmt.Println("Reading servers from", filename)

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				s := strings.Split(scanner.Text(), " ")
				if len(s) == 2 {
					alias, addr := s[0], s[1]
					result[alias] = addr
				}
			}
		}
	}

	return result
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
