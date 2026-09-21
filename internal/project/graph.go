package project

import (
	"slices"
	"strings"
)

func relationships(r *Record) []struct {
	field   string
	targets []string
} {
	return []struct {
		field   string
		targets []string
	}{
		{"depends_on", r.DependsOn}, {"members", r.Members},
		{"blocks", r.Blocks}, {"relates_to", r.RelatesTo}, {"work", r.Work},
	}
}

// Validate reports identity and relationship problems across a complete
// record set. Callers substituting a candidate record must pass the whole set.
func Validate(records []*Record) []Diagnostic {
	var ds []Diagnostic
	index := map[string][]*Record{}
	for _, r := range records {
		// An independent metadata error must not hide a duplicate identity or
		// let another record's references silently select this definition.
		if r.ID != "" {
			index[r.ID] = append(index[r.ID], r)
		}
	}
	unique := map[string]*Record{}
	for _, r := range records {
		matches := index[r.ID]
		if len(matches) > 1 {
			paths := make([]string, len(matches))
			for i, match := range matches {
				paths[i] = match.Path
			}
			ds = append(ds, Diagnostic{Path: r.Path, Field: "id", Message: "duplicate " + r.ID + " in " + strings.Join(paths, ", ")})
		} else if len(matches) == 1 {
			unique[r.ID] = r
		}
		for _, rel := range relationships(r) {
			for _, target := range rel.targets {
				message := ""
				switch matches := index[target]; {
				case target == r.ID:
					message = "self-reference is not allowed"
				case len(matches) == 0:
					message = "unresolved target " + target
				case len(matches) > 1:
					message = "ambiguous target " + target
				case rel.field != "relates_to" && matches[0].Type != "work":
					message = "target " + target + " must be work"
				}
				if message != "" {
					ds = append(ds, Diagnostic{Path: r.Path, Field: rel.field, Message: message})
				}
			}
		}
	}
	// A term's title is the term, so two records for one term are a conflict.
	terms := map[string]*Record{}
	for _, r := range records {
		if r.Type != "term" || r.Title == "" {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(r.Title))
		if first := terms[name]; first != nil {
			ds = append(ds, Diagnostic{Path: r.Path, Field: "title", Message: "term already defined by " + first.ID + " in " + first.Path})
		} else {
			terms[name] = r
		}
	}
	// Membership is decomposition, not an execution-order edge. Mixing these
	// graphs would falsely reject valid parent/member dependency relationships.
	for _, field := range []string{"depends_on", "members"} {
		state := map[string]int{}
		var stack []string
		var visit func(*Record)
		visit = func(r *Record) {
			state[r.ID] = 1
			stack = append(stack, r.ID)
			targets := r.DependsOn
			if field == "members" {
				targets = r.Members
			}
			for _, id := range targets {
				target := unique[id]
				if target == nil || target.Type != "work" || id == r.ID {
					continue
				}
				switch state[id] {
				case 0:
					visit(target)
				case 1:
					start := slices.Index(stack, id)
					cycle := append(append([]string{}, stack[start:]...), id)
					ds = append(ds, Diagnostic{Path: r.Path, Field: field, Message: "cycle " + strings.Join(cycle, " -> ")})
				}
			}
			stack = stack[:len(stack)-1]
			state[r.ID] = 2
		}
		for _, r := range records {
			if unique[r.ID] == r && r.Type == "work" && state[r.ID] == 0 {
				visit(r)
			}
		}
	}
	return ds
}
