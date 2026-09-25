package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/mascah/grove/internal/deps"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

// runDeps prints the dependency overview of unfinished work, or a preview of
// an explicit selection, bound to this checkout's records (G-161). Like
// versions, an inspection that could not read every source still prints and
// exits 1.
func runDeps(p *project.Project, a invocation, out, errOut io.Writer) int {
	v := deps.Overview(p.Records, false)
	if len(a.ids) != 0 {
		var err error
		if v, err = deps.Preview(p.Records, a.ids); err != nil {
			report(errOut, err)
			return 1
		}
	}
	code := 0
	checkout := map[string]any{"root": p.Root, "ref": nil, "head": nil}
	target := ""
	res, err := versions.Inspect(p.Root, "")
	if err != nil {
		v.Notes = append(v.Notes, "Other branches and checkouts could not be read ("+err.Error()+"), so nothing is said about other versions.")
		code = 1
	} else {
		target = res.Target
		if !res.Complete {
			code = 1
		}
		var here *versions.Source // this checkout, whose records the view holds
		for _, s := range res.Sources {
			if s.Kind == "live" && s.GitDir == res.GitDir {
				here = s
				checkout["head"] = s.Commit
				if s.Ref != "" {
					checkout["ref"] = s.Ref
				}
			}
		}
		v.Compare(res, here)
	}
	v.Deliver(target, deps.Ancestry(context.Background(), p.Root))
	if code != 0 {
		fmt.Fprintln(errOut, "grove: some sources could not be inspected; the result is incomplete")
	}

	var buffer bytes.Buffer
	if a.json {
		encoder := json.NewEncoder(&buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		var t any // null without a target
		if target != "" {
			t = target
		}
		encoder.Encode(struct {
			Checkout map[string]any `json:"checkout"`
			Target   any            `json:"target"`
			*deps.View
		}{checkout, t, v})
	} else {
		depsText(&buffer, v, checkout, target)
	}
	if writeResult(out, errOut, buffer.Bytes()) != 0 {
		return 1
	}
	return code
}

func depsText(w *bytes.Buffer, v *deps.View, checkout map[string]any, target string) {
	line := "Checkout: " + visible(checkout["root"].(string))
	if ref, ok := checkout["ref"].(string); ok {
		line += " on " + visible(ref)
	}
	if head, ok := checkout["head"].(string); ok {
		line += " at " + short(head)
	}
	if target != "" {
		line += "; target " + visible(target)
	}
	fmt.Fprintln(w, line)
	list := func(ids []string) string {
		if len(ids) == 0 {
			return "-"
		}
		return visible(strings.Join(ids, " "))
	}
	table := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	var rows, outside []deps.Item
	for _, it := range v.Items {
		if it.Outside {
			outside = append(outside, it)
		} else {
			rows = append(rows, it)
		}
	}
	switch {
	case v.Selected != nil:
		fmt.Fprintf(w, "Selected: %s\nOrder: %s\n", list(v.Selected), list(v.Order))
		fmt.Fprintln(table, "ORDER\tID\tSTATUS\tNEEDS\tDELIVERY\tTITLE")
		for i, it := range rows {
			fmt.Fprintf(table, "%d\t%s\t%s\t%s\t%s\t%s\n", i+1, it.ID, it.Status, list(it.Needs), visible(it.Delivery), visible(it.Title))
		}
	case rows == nil:
		fmt.Fprintln(w, "No unfinished work.")
	default:
		fmt.Fprintln(table, "GROUP\tLAYER\tID\tSTATUS\tNEEDS\tUNLOCKS\tDELIVERY\tTITLE")
		for _, it := range rows {
			fmt.Fprintf(table, "%d\t%d\t%s\t%s\t%s\t%s\t%s\t%s\n", it.Group, it.Layer, it.ID, it.Status, list(it.Needs), list(it.Unlocks), visible(it.Delivery), visible(it.Title))
		}
	}
	table.Flush()
	if outside != nil {
		if v.Selected != nil {
			fmt.Fprintln(w, "\nPrerequisites outside the selection, not added:")
		} else {
			fmt.Fprintln(w, "\nPrerequisites that are not unfinished:")
		}
		fmt.Fprintln(table, "ID\tSTATUS\tNEEDED BY\tDELIVERY\tTITLE")
		for _, it := range outside {
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\n", it.ID, it.Status, list(it.NeededBy), visible(it.Delivery), visible(it.Title))
		}
		table.Flush()
	}
	if len(v.Questions) != 0 {
		fmt.Fprintln(w, "\nOpen questions:")
		for _, q := range v.Questions {
			fmt.Fprintf(table, "%s\tblocks %s\t%s\n", q.ID, list(q.Blocks), visible(q.Title))
		}
		table.Flush()
	}
	if len(v.Notes) != 0 {
		fmt.Fprintln(w, "\nNotes:")
		for _, n := range v.Notes {
			fmt.Fprintln(w, "- "+visible(n))
		}
	}
	if v.Selected != nil {
		fmt.Fprintf(w, "\nThis preview adds no work, starts nothing, and authorizes nothing; grove context %s reads the selection in this order.\n", list(v.Order))
	} else if rows != nil {
		fmt.Fprintln(w, "\nOnly depends_on orders layers. Equal layers have no declared order, which is not evidence that they can proceed in parallel.")
	}
}
