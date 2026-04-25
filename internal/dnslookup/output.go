package dnslookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const jsonSchemaVersion = 1

// ExitStatus returns 0 on success, 2 if no server succeeded or strict saw any failure.
func ExitStatus(results []QueryResult, strict bool) int {
	var anySuccess, anyFail bool
	for _, r := range results {
		if len(r.Errs) > 0 {
			anyFail = true
		}
		if len(r.ByType) > 0 {
			anySuccess = true
		}
	}
	if strict && anyFail {
		return 2
	}
	if !anySuccess {
		return 2
	}
	return 0
}

// WriteJSON writes machine-readable output to w.
func WriteJSON(w io.Writer, name string, types []RecordType, results []QueryResult) error {
	type row struct {
		Server    string              `json:"server"`
		Address   string              `json:"address"`
		OK        bool                `json:"ok"`
		ByType    map[string][]string `json:"by_type,omitempty"`
		Errors    []string            `json:"errors,omitempty"`
		Rcode     int                 `json:"rcode"`
		Truncated bool                `json:"truncated"`
	}
	typeNames := make([]string, 0, len(types))
	for _, t := range types {
		typeNames = append(typeNames, t.String())
	}
	rows := make([]row, 0, len(results))
	for _, r := range results {
		errStrs := make([]string, 0, len(r.Errs))
		for _, e := range r.Errs {
			errStrs = append(errStrs, e.Error())
		}
		rows = append(rows, row{
			Server:    r.Server.Name,
			Address:   r.Server.Addr,
			OK:        len(r.Errs) == 0,
			ByType:    r.ByType,
			Errors:    errStrs,
			Rcode:     r.Rcode,
			Truncated: r.Truncated,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]any{
		"version": jsonSchemaVersion,
		"name":    name,
		"types":   typeNames,
		"results": rows,
	})
}

// WriteTable writes aligned rows for humans. If useColor, ANSI codes are used when w is a terminal.
func WriteTable(w io.Writer, types []RecordType, results []QueryResult, useColor bool) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintf(tw, "SERVER\tADDRESS\tRESULT\n"); err != nil {
		return err
	}
	const (
		reset  = "\033[0m"
		red    = "\033[31m"
		yellow = "\033[33m"
	)
	for _, r := range results {
		ans := strings.Join(FlatAnswers(types, r), ", ")
		errJoined := errors.Join(r.Errs...)
		line := ""
		switch {
		case errJoined != nil:
			msg := errJoined.Error()
			if useColor {
				line = fmt.Sprintf("%s%s\t%s\t%s%s\n", red, r.Server.Name, r.Server.Addr, msg, reset)
			} else {
				line = fmt.Sprintf("%s\t%s\t%s\n", r.Server.Name, r.Server.Addr, msg)
			}
		case r.Truncated && useColor:
			line = fmt.Sprintf("%s%s\t%s\t%s (truncated)%s\n", yellow, r.Server.Name, r.Server.Addr, ans, reset)
		default:
			line = fmt.Sprintf("%s\t%s\t%s\n", r.Server.Name, r.Server.Addr, ans)
		}
		if _, err := io.WriteString(tw, line); err != nil {
			return err
		}
	}
	return tw.Flush()
}
