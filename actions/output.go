package actions

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

type Output interface {
	ToPrettyOutput() string
}

type ListOutput struct {
	Labels []string
	Rows   []map[string]string
}

func (o *ListOutput) AddRow(m map[string]string) {
	o.Rows = append(o.Rows, m)
}

func (o *ListOutput) ToPrettyOutput() string {
	b := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(b, 0, 8, 2, '\t', 0)

	// Print Heading Row
	for _, l := range o.Labels {
		fmt.Fprintf(w, "%v\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	// Print Rows
	for _, r := range o.Rows {
		row := make([]string, len(o.Labels))
		for i, l := range o.Labels {
			row[i] = r[l]
		}

		fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
	}

	w.Flush()
	return b.String()
}

type PlainOutput struct {
	Output string
}

func (o PlainOutput) ToPrettyOutput() string {
	return o.Output
}

type DetailOutput struct {
	Details map[string]string
}

func (o DetailOutput) ToPrettyOutput() string {
	b := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(b, 0, 8, 2, '\t', 0)

	var keys []string
	for k := range o.Details {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%v\n", k, o.Details[k])
	}

	w.Flush()
	return b.String()
}
