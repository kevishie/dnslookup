package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"

	"dnslookup/internal/dnslookup"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "per-query timeout")
	concurrency := flag.Int("c", 32, "max concurrent server queries")
	jsonOut := flag.Bool("json", false, "print JSON to stdout")
	strict := flag.Bool("strict", false, "exit 2 if any server fails")
	noColor := flag.Bool("no-color", false, "disable ANSI colors")

	var typeArgs []string
	flag.Func("t", "record type (repeatable): A, AAAA, NS, MX, TXT, CNAME", func(s string) error {
		typeArgs = append(typeArgs, s)
		return nil
	})

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: dnslookup [flags] <name>\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	name := flag.Arg(0)

	types, err := dnslookup.ParseRecordTypes(typeArgs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	servers, cfgPath, err := dnslookup.LoadServers(home, os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if cfgPath != "" {
		fmt.Fprintf(os.Stderr, "Reading servers from %s\n\n", cfgPath)
	}

	client := dnslookup.NewClient(*timeout)
	results := dnslookup.Run(context.Background(), client, servers, name, types, *concurrency)

	if *jsonOut {
		if err := dnslookup.WriteJSON(os.Stdout, name, types, results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		useColor := !*noColor && term.IsTerminal(int(os.Stdout.Fd()))
		if err := dnslookup.WriteTable(os.Stdout, types, results, useColor); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	os.Exit(dnslookup.ExitStatus(results, *strict))
}
