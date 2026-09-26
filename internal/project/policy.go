package project

import (
	"path"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Policy is grove.yaml's optional policy: mapping, the owner's standing
// delegation (G-182): what grove sweep may do to a candidate in review with
// no per-candidate human act. A nil Policy, the default, means nothing is
// automatic; a missing section means that act waits for the owner.
type Policy struct {
	BudgetUSD        string   // policy.budget: every automatic spend in one sweep
	Resolve          bool     // policy.resolve: start a resolution attempt for a conflict
	ResolveBudgetUSD string   // policy.resolve.budget; "" uses run: budget
	Approve          bool     // policy.approve: approve a candidate that meets the conditions
	Verify           []string // policy.approve.verify: commands the merged result must pass
	MaxLines         int      // policy.approve.max_lines; 0 is no bound
	Never            []string // policy.approve.never: paths whose change always waits
	Integrate        bool     // policy.integrate: merge and write done after a delegated approval
}

// Matches returns the never pattern a project-relative path matches, if any:
// path.Match's syntax, or DIR/** for everything under DIR.
func (p *Policy) Matches(file string) (string, bool) {
	for _, pattern := range p.Never {
		if dir, ok := strings.CutSuffix(pattern, "/**"); ok {
			if strings.HasPrefix(file, dir+"/") {
				return pattern, true
			}
		} else if ok, _ := path.Match(pattern, file); ok {
			return pattern, true
		}
	}
	return "", false
}

// policyField reads the policy: mapping. Problems name the key as
// policy.KEY or policy.SECTION.KEY, and any problem is a diagnostic, so a
// malformed policy stops every command rather than half applying.
func (m *metadata) policyField() *Policy {
	top := m.section("policy", "a mapping of budget, resolve, approve and integrate")
	if top == nil {
		return nil
	}
	defer top.report(m, "policy")
	var p Policy
	for key := range top.fields {
		if key != "budget" && key != "resolve" && key != "approve" && key != "integrate" {
			top.problem(key, "unknown key; policy: takes budget, resolve, approve and integrate")
		}
	}
	p.BudgetUSD = top.budgetField("budget")
	if resolve := top.section("resolve", "a mapping, {} for the run: defaults"); resolve != nil {
		p.Resolve = true
		for key := range resolve.fields {
			if key != "budget" {
				resolve.problem(key, "unknown key; policy.resolve takes budget")
			}
		}
		p.ResolveBudgetUSD = resolve.budgetField("budget")
		if p.BudgetUSD == "" {
			top.problem("budget", "required with resolve: the aggregate every automatic attempt in one sweep spends")
		}
		resolve.report(top, "resolve")
	}
	if approve := top.section("approve", "a mapping of verify, max_lines and never"); approve != nil {
		p.Approve = true
		for key := range approve.fields {
			if key != "verify" && key != "max_lines" && key != "never" {
				approve.problem(key, "unknown key; policy.approve takes verify, max_lines and never")
			}
		}
		p.Verify = approve.strings("verify", "a command")
		if len(p.Verify) == 0 {
			approve.problem("verify", "required: the commands the merged result must pass, one or more")
		}
		if n, ok := approve.integerField("max_lines", false); ok && n <= 0 {
			approve.problem("max_lines", "expected a positive number of lines")
		} else {
			p.MaxLines = n
		}
		p.Never = approve.strings("never", "a pattern")
		for _, pattern := range p.Never {
			rest := strings.TrimSuffix(pattern, "/**")
			if _, err := path.Match(rest, ""); err != nil || strings.Contains(rest, "**") || !dedicated(rest) {
				approve.problem("never", "bad pattern "+pattern+": a project-relative path.Match pattern, or DIR/** for a directory's tree")
			}
		}
		approve.report(top, "approve")
	}
	if n, ok := top.fields["integrate"]; ok {
		if n.Kind != yaml.ScalarNode || n.Tag != "!!bool" || n.Decode(&p.Integrate) != nil {
			top.problem("integrate", "expected true or false")
		} else if p.Integrate && !p.Approve {
			top.problem("integrate", "needs approve: integration follows only a delegated approval")
		}
	}
	return &p
}

// section is the mapping under key, or nil when it is absent or, reported,
// not a mapping.
func (m *metadata) section(key, what string) *metadata {
	n, ok := m.fields[key]
	if !ok {
		return nil
	}
	if n.Kind != yaml.MappingNode || n.Tag != "!!map" {
		m.problem(key, "expected "+what)
		return nil
	}
	sub := &metadata{path: m.path, offset: m.offset, fields: map[string]*yaml.Node{}}
	sub.addFields(n.Content)
	return sub
}

// report moves a section's problems to its parent, named under key.
func (m *metadata) report(parent *metadata, key string) {
	for _, d := range m.errors {
		d.Field = strings.TrimSuffix(key+"."+d.Field, ".")
		parent.errors = append(parent.errors, d)
	}
}

func (m *metadata) budgetField(key string) string {
	v, ok := m.fields[key]
	if !ok {
		return ""
	}
	if v.Kind == yaml.ScalarNode && (v.Tag == "!!int" || v.Tag == "!!float" || v.Tag == "!!str") && ValidBudget(v.Value) {
		return v.Value
	}
	m.problem(key, "expected a positive decimal dollar amount")
	return ""
}

// strings reads a list of nonempty strings.
func (m *metadata) strings(key, what string) []string {
	n, ok := m.fields[key]
	if !ok {
		return nil
	}
	if n.Kind != yaml.SequenceNode || n.Tag != "!!seq" {
		m.problem(key, "expected a list, each "+what)
		return nil
	}
	var out []string
	for _, item := range n.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" || strings.TrimSpace(item.Value) == "" {
			m.problem(key, "each item must be "+what+", a nonempty string")
			continue
		}
		out = append(out, item.Value)
	}
	return out
}
