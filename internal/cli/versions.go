package cli

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/mascah/grove/internal/versions"
)

// runVersions prints every observation of the selected project's records
// across branch tips and worktrees. Sources and their diagnostics go to
// stderr; rows or JSON go to stdout even when the result is incomplete.
func runVersions(root string, a invocation, out, errOut io.Writer) int {
	res, err := versions.Inspect(root, a.id)
	if err != nil {
		report(errOut, err)
		return 1
	}
	fmt.Fprintf(errOut, "Repository: %s\n", visible(res.Repository))
	for _, s := range res.Sources {
		fmt.Fprintf(errOut, "Source: %s\n", visible(describeSource(s)))
		if s.Note != "" {
			fmt.Fprintf(errOut, "  note: %s\n", visible(s.Note))
		}
		for _, d := range s.Diagnostics {
			fmt.Fprintf(errOut, "  %s\n", visible(d))
		}
	}
	if res.Target != "" {
		fmt.Fprintf(errOut, "Target: %s\n", visible(res.Target))
	}
	for _, n := range res.Notes {
		fmt.Fprintf(errOut, "Target: %s\n", visible(n))
	}
	for _, g := range res.Groups {
		for _, n := range g.Notes {
			fmt.Fprintf(errOut, "%s: %s\n", g.ID, visible(n))
		}
	}
	code := 0
	if !res.Complete {
		fmt.Fprintln(errOut, "grove: some sources could not be inspected; the result is incomplete")
		code = 1
	}
	if a.id != "" && len(res.Groups) == 0 {
		fmt.Fprintf(errOut, "grove: record %s not found in any valid source\n", visible(a.id))
		code = 1
	}
	var content []byte
	if a.json {
		content = marshal(versionsJSON(res))
	} else {
		var buffer bytes.Buffer
		table := tabwriter.NewWriter(&buffer, 0, 4, 2, ' ', 0)
		fmt.Fprintln(table, "ID\tSTATUS\tSOURCE\tCHANGE\tCURRENT\tTARGET\tSELECTOR")
		for _, g := range res.Groups {
			for _, v := range g.Versions {
				status, change := "-", "-"
				if v.Record != nil {
					status = cmp.Or(v.Record.Status, "-") // a page has no status
				}
				if v.Change != "" {
					change = v.Change
				}
				current := "yes"
				if v.Older != "" {
					current = "older"
				}
				onTarget := "-"
				if res.Target != "" {
					onTarget = map[bool]string{true: "yes", false: "no"}[v.OnTarget]
				}
				fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", g.ID, visible(status), visible(sourceCell(v.Source)), change, current, onTarget, cmp.Or(v.Selector, "-"))
			}
		}
		table.Flush()
		content = buffer.Bytes()
	}
	if writeResult(out, errOut, content) != 0 {
		return 1
	}
	return code
}

// short abbreviates a commit for one-line output.
func short(commit string) string {
	if len(commit) > 12 {
		return commit[:12]
	}
	return commit
}

func sourceCell(s *versions.Source) string {
	ref := s.Ref
	if ref == "" {
		ref = "detached"
	}
	if s.Kind == "live" {
		return "live " + s.Locator + " " + ref
	}
	return "committed " + ref
}

func describeSource(s *versions.Source) string {
	text := sourceCell(s) + " " + short(s.Commit)
	if s.Kind == "live" {
		text += " " + s.Worktree
	}
	switch {
	case !s.Present && len(s.Diagnostics) == 0:
		text += " (no project)"
	case len(s.Diagnostics) != 0:
		text += " (invalid)"
	}
	return text
}

func versionsJSON(res *versions.Result) map[string]any {
	sources := make([]any, 0, len(res.Sources))
	for _, s := range res.Sources {
		sources = append(sources, sourceJSON(s, true))
	}
	records := make([]any, 0, len(res.Groups))
	for _, g := range res.Groups {
		vs := make([]any, 0, len(g.Versions))
		for _, v := range g.Versions {
			o := sourceJSON(v.Source, false)
			o["path"], o["current"] = v.Path, v.Older == ""
			if v.Older != "" {
				o["older"] = v.Older
			}
			o["on_target"] = nil
			if res.Target != "" {
				o["on_target"] = v.OnTarget
			}
			if v.Change != "" {
				o["change"] = v.Change
			}
			if v.HeadPath != "" {
				o["head_path"] = v.HeadPath
			}
			if v.Record != nil {
				o["selector"], o["revision"], o["config_revision"] = v.Selector, v.Revision, v.Source.ConfigRevision
				o["type"], o["title"], o["status"], o["source"] = v.Record.Type, v.Record.Title, v.Record.Status, string(v.Record.Source)
			}
			vs = append(vs, o)
		}
		notes := make([]string, 0, len(g.Notes))
		records = append(records, map[string]any{"id": g.ID, "versions": vs, "notes": append(notes, g.Notes...)})
	}
	var target any // null without a target
	if res.Target != "" {
		target = res.Target
	}
	return map[string]any{
		"project": res.Project, "repository": res.Repository, "prefix": res.Prefix,
		"complete": res.Complete, "sources": sources, "records": records,
		"target": target, "notes": append(make([]string, 0, len(res.Notes)), res.Notes...),
	}
}

func sourceJSON(s *versions.Source, full bool) map[string]any {
	o := map[string]any{"kind": s.Kind, "commit": s.Commit}
	o["ref"] = nil
	if s.Ref != "" {
		o["ref"] = s.Ref
	}
	if s.Kind == "live" {
		o["worktree"], o["detached"] = s.Worktree, s.Ref == ""
		if s.Locator != "" {
			o["locator"] = s.Locator
		}
	}
	if !full {
		return o
	}
	o["present"], o["valid"] = s.Present, s.Valid
	if s.Valid {
		o["config_revision"] = s.ConfigRevision
	}
	if s.Note != "" {
		o["note"] = s.Note
	}
	diagnostics := make([]string, 0, len(s.Diagnostics))
	o["diagnostics"] = append(diagnostics, s.Diagnostics...)
	return o
}
