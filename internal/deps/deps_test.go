package deps

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/versions"
)

func work(id, status string, needs ...string) *project.Record {
	return &project.Record{ID: id, Type: "work", Title: "Title " + id, Status: status, DependsOn: needs, Source: []byte(id + status)}
}

// backlog is G-165's synthetic backlog: a chain, a shared prerequisite, a
// convergence, unrelated work, a candidate in review, an abandoned
// prerequisite, and a blocking question.
func backlog() []*project.Record {
	rs := []*project.Record{
		work("S-01", "done"), work("S-02", "proposed", "S-01"), work("S-03", "active", "S-02"),
		work("S-04", "proposed", "S-02"), work("S-05", "proposed", "S-03", "S-04"), work("S-06", "proposed", "S-05"),
		work("S-07", "review", "S-03"), work("S-08", "proposed", "S-04", "S-10"), work("S-09", "proposed"),
		work("S-10", "abandoned"), work("S-11", "active"), work("S-13", "proposed", "S-05"),
		{ID: "S-12", Type: "question", Title: "Round how?", Status: "open", Blocks: []string{"S-11"}},
		{ID: "S-14", Type: "question", Title: "Answered", Status: "resolved", Blocks: []string{"S-05"}},
	}
	rs[0].Candidate, rs[6].Candidate = "1111111aaaa", "7777777bbbb"
	return rs
}

func TestOverviewLayersGroupsAndUnlocks(t *testing.T) {
	v := Overview(backlog(), false)
	type row struct {
		id            string
		group, layer  int
		needs, unlock string
	}
	var got []row
	var outside []string
	for _, it := range v.Items {
		if it.Outside {
			outside = append(outside, it.ID+"<"+strings.Join(it.NeededBy, ","))
			continue
		}
		got = append(got, row{it.ID, it.Group, it.Layer, strings.Join(it.Needs, ","), strings.Join(it.Unlocks, ",")})
	}
	want := []row{
		{"S-02", 1, 0, "S-01", "S-03,S-04"}, {"S-03", 1, 1, "S-02", "S-05,S-07"}, {"S-04", 1, 1, "S-02", "S-05,S-08"},
		{"S-05", 1, 2, "S-03,S-04", "S-06,S-13"}, {"S-07", 1, 2, "S-03", ""}, {"S-08", 1, 2, "S-04,S-10", ""},
		{"S-06", 1, 3, "S-05", ""}, {"S-13", 1, 3, "S-05", ""},
		{"S-09", 2, 0, "", ""}, {"S-11", 3, 0, "", ""},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("rows\n got %v\nwant %v", got, want)
	}
	if !reflect.DeepEqual(outside, []string{"S-01<S-02", "S-10<S-08"}) {
		t.Errorf("outside %v", outside)
	}
	if !reflect.DeepEqual(v.Questions, []Question{{"S-12", "Round how?", []string{"S-11"}}}) {
		t.Errorf("questions %v", v.Questions)
	}
	if len(v.Notes) != 1 || !strings.Contains(v.Notes[0], "S-10 is abandoned") || !strings.Contains(v.Notes[0], "S-08 still needs it") {
		t.Errorf("notes %q", v.Notes)
	}
	if v.Selected != nil {
		t.Errorf("an overview selects nothing: %v", v.Selected)
	}
}

func TestPreviewKeepsTheSelectionAndAddsNothing(t *testing.T) {
	v, err := Preview(backlog(), []string{"S-05", "S-03", "S-09"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v.Selected, []string{"S-05", "S-03", "S-09"}) || !reflect.DeepEqual(v.Order, []string{"S-03", "S-05", "S-09"}) {
		t.Fatalf("selected %v order %v", v.Selected, v.Order)
	}
	var outside []string
	for _, it := range v.Items {
		if it.Outside {
			outside = append(outside, it.ID+"<"+strings.Join(it.NeededBy, ","))
		}
	}
	// Transitive prerequisites through unselected work are listed, never added.
	if !reflect.DeepEqual(outside, []string{"S-01<S-03,S-05", "S-02<S-03,S-05", "S-04<S-05"}) {
		t.Errorf("outside %v", outside)
	}
	notes := strings.Join(v.Notes, "\n")
	for _, want := range []string{"S-02, S-04: unfinished and not selected", "No declared order between S-03 and S-09; S-05 and S-09.", "not evidence"} {
		if !strings.Contains(notes, want) {
			t.Errorf("notes lack %q:\n%s", want, notes)
		}
	}
	if len(v.Questions) != 0 {
		t.Errorf("a resolved question and one blocking other work are not listed: %v", v.Questions)
	}
	if _, err := Preview(backlog(), []string{"S-12"}); err == nil {
		t.Error("a question was selectable")
	}
	if _, err := Preview(backlog(), []string{"S-99"}); err == nil {
		t.Error("an unknown ID was selectable")
	}
}

// Membership and priority never order; only depends_on does, including
// through unselected work.
func TestOnlyDependenciesOrder(t *testing.T) {
	one, five := 1, 5
	a, b, c := work("A", "proposed", "C"), work("B", "proposed"), work("C", "proposed", "B")
	a.Priority, b.Priority = &five, &one
	b.Members = []string{"A"}
	v, err := Preview([]*project.Record{a, b, c}, []string{"A", "B"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(v.Order, []string{"B", "A"}) {
		t.Errorf("order %v", v.Order)
	}
	v, _ = Preview([]*project.Record{work("A", "proposed"), b}, []string{"B", "A"})
	if !reflect.DeepEqual(v.Order, []string{"B", "A"}) {
		t.Errorf("membership or priority reordered: %v", v.Order)
	}
}

func TestDeliverExplainsEachStatus(t *testing.T) {
	rs := append(backlog(), work("S-20", "done"), work("S-21", "done"), work("S-30", "proposed", "S-20", "S-21", "S-07"))
	rs[len(rs)-3].Candidate = "2020202" // in HEAD, not on main
	rs[len(rs)-2].Candidate = "2121212" // unreadable
	v, err := Preview(rs, []string{"S-30", "S-08", "S-02"})
	if err != nil {
		t.Fatal(err)
	}
	v.Deliver("main", func(commit, ref string) (bool, error) {
		switch {
		case commit == "2121212":
			return false, errors.New("bad object")
		case commit == "7777777bbbb":
			return false, nil
		case commit == "2020202":
			return ref == "HEAD", nil
		}
		return true, nil
	}, nil)
	got := map[string]string{}
	for _, it := range v.Items {
		got[it.ID] = it.Delivery
	}
	want := map[string]string{
		"S-01": "candidate 1111111 in HEAD, on main",
		"S-02": "awaiting implementation",
		"S-03": "awaiting implementation",
		"S-04": "awaiting implementation",
		"S-07": "awaiting review; candidate 7777777 not in HEAD, not on main",
		"S-08": "awaiting implementation",
		"S-10": "abandoned: will not be delivered",
		"S-20": "candidate 2020202 in HEAD, not on main",
		"S-21": "candidate 2121212 cannot be read here",
		"S-30": "awaiting implementation",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("delivery\n got %v\nwant %v", got, want)
	}
	v, _ = Preview([]*project.Record{work("H", "done"), work("W", "proposed", "H")}, []string{"W"})
	v.Deliver("", func(string, string) (bool, error) { t.Fatal("asked Git without a candidate"); return false, nil }, nil)
	if d := v.Items[1].Delivery; !strings.HasPrefix(d, "done without a candidate: delivery unrecorded") {
		t.Errorf("historical done: %q", d)
	}
}

// Compare describes other versions and never merges their edges.
func TestCompareDescribesOtherVersions(t *testing.T) {
	here := &versions.Source{Kind: "live", Ref: "refs/heads/main", Locator: ".", GitDir: "/g"}
	branch := &versions.Source{Kind: "committed", Ref: "refs/heads/feature"}
	other := &versions.Source{Kind: "live", Ref: "refs/heads/other", Locator: "other", GitDir: "/g/w"}
	a, b, c, d := work("A", "proposed", "B"), work("B", "proposed"), work("C", "proposed"), work("D", "proposed")
	e, eReordered := work("E", "proposed", "B", "C"), work("E", "active", "C", "B")
	aElsewhere := work("A", "active", "B", "C")
	aElsewhere.Source = []byte("A elsewhere")
	rev := func(r *project.Record) string { return project.Revision(r.Source) }
	stale := work("D", "proposed")
	stale.Source = []byte("changed since")
	res := &versions.Result{GitDir: "/g", Complete: false, Sources: []*versions.Source{branch, here, other}, Groups: []versions.Group{
		{ID: "A", Versions: []versions.Version{
			{Source: here, Record: a, Revision: rev(a), Change: "unchanged", Older: "branch feature changed it"},
			{Source: branch, Record: aElsewhere, Revision: rev(aElsewhere)},
			{Source: other, Record: aElsewhere, Revision: rev(aElsewhere), Change: "unchanged"},
		}},
		{ID: "B", Versions: []versions.Version{{Source: here, Record: b, Revision: rev(b), Change: "modified"}}, Notes: []string{"x and y could not be ordered"}},
		{ID: "C", Versions: []versions.Version{
			{Source: here, Record: c, Revision: rev(c), Change: "unchanged"},
			{Source: branch, Revision: ""}, // deleted there, and current: a divergence
		}},
		{ID: "D", Versions: []versions.Version{{Source: here, Record: stale, Revision: rev(stale), Change: "unchanged"}}},
		{ID: "E", Versions: []versions.Version{
			{Source: here, Record: e, Revision: rev(e), Change: "unchanged", Older: "branch feature changed it"},
			{Source: branch, Record: eReordered, Revision: rev(eReordered)},
		}},
	}}
	v, err := Preview([]*project.Record{a, b, c, d, e}, []string{"A", "C", "D", "E"})
	if err != nil {
		t.Fatal(err)
	}
	v.Notes = nil
	v.Compare(res, here)
	want := []string{
		"Some branches or checkouts could not be read (grove versions lists them), so what is said here about other versions may be incomplete.",
		"A here is older than its current version: active on branch feature, checkout other (other).",
		`A depends on other work elsewhere: branch feature ["B" "C"]; this uses this checkout's ["B"].`,
		"C diverges; other current versions: deleted on branch feature.",
		"D changed while it was being read; rerun.",
		"E here is older than its current version: active on branch feature.",
		"B has uncommitted changes (modified) in this checkout, which is what is read here.",
		"B: x and y could not be ordered",
	}
	if !reflect.DeepEqual(v.Notes, want) {
		t.Errorf("notes\n got %q\nwant %q", v.Notes, want)
	}
	if !reflect.DeepEqual(v.Items[0].Needs, []string{"B"}) {
		t.Errorf("edges were merged: %v", v.Items[0].Needs)
	}
}

// A board's current view can hold work whose prerequisite's current state
// deletes it, and every work includes done and abandoned rows.
func TestOverviewEveryWorkAndMissingPrerequisites(t *testing.T) {
	v := Overview(append(backlog(), work("S-20", "proposed", "S-99")), false)
	last := v.Items[len(v.Items)-1]
	if last.ID != "S-99" || !last.Outside || last.Status != "" || !reflect.DeepEqual(last.NeededBy, []string{"S-20"}) {
		t.Fatalf("missing prerequisite %+v", last)
	}
	v.Deliver("", func(string, string) (bool, error) { return false, nil }, nil)
	if last = v.Items[len(v.Items)-1]; !strings.Contains(last.Delivery, "not among the records read") {
		t.Errorf("delivery %q", last.Delivery)
	}
	every := Overview(backlog(), true)
	var rows []string
	for _, it := range every.Items {
		if it.Outside {
			t.Errorf("every work lists %s outside", it.ID)
		}
		rows = append(rows, it.ID)
	}
	if len(rows) != 12 || every.Items[0].ID != "S-01" || every.Items[0].Layer != 0 {
		t.Errorf("every work %v", rows)
	}
}

// A selection's candidates are merged in its order, skipping one already on
// the target, and what follows the first conflict is named as not tried.
func TestDeliverMergesInTheSelectionsOrder(t *testing.T) {
	rs := []*project.Record{work("A", "review"), work("B", "review"), work("C", "review"), work("D", "review"), work("E", "proposed")}
	for _, r := range rs[:4] {
		r.Candidate = strings.Repeat(strings.ToLower(r.ID), 7)
	}
	v, err := Preview(rs, []string{"D", "C", "B", "A", "E"})
	if err != nil {
		t.Fatal(err)
	}
	var asked [][]string
	v.Deliver("main", func(string, string) (bool, error) { return false, nil }, func(commits []string) ([]versions.Merge, error) {
		asked = append(asked, commits)
		var out []versions.Merge
		for i, c := range commits {
			m := versions.Merge{Target: "ttttttt", Commit: c, Outcome: "clean", Conflicts: []string{}}
			switch {
			case c == "ccccccc":
				m.Outcome = "integrated"
			case len(commits) > 1 && i == 1:
				m.Outcome, m.Conflicts = "conflict", []string{"x.go", "y.go"}
			}
			out = append(out, m)
			if m.Outcome == "conflict" {
				break
			}
		}
		return out, nil
	})
	if want := [][]string{{"ddddddd"}, {"ccccccc"}, {"bbbbbbb"}, {"aaaaaaa"}, {"ddddddd", "bbbbbbb", "aaaaaaa"}}; !reflect.DeepEqual(asked, want) {
		t.Fatalf("asked %v", asked)
	}
	if d := v.Items[1].Delivery; d != "awaiting review; candidate ccccccc not in HEAD, not on main; integrated: main at ttttttt holds it" {
		t.Errorf("C: %q", d)
	}
	if len(v.MergeOrder) != 2 || v.MergeOrder[1].ID != "B" {
		t.Errorf("merge order %+v", v.MergeOrder)
	}
	if want := "Merged into main at ttttttt in this order, each onto the ones before, in objects only: D merges cleanly, B conflicts in x.go, y.go; the first conflict is B's, and A after it was not tried. Grove chose no order, and a clean order is not evidence that the changes work together."; !slices.Contains(v.Notes, want) {
		t.Errorf("notes %q", v.Notes)
	}
}
