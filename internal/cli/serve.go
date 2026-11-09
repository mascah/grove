package cli

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
)

//go:embed page.html
var page []byte

// server answers the page and the project value it reads.
type server struct {
	p     project
	mu    sync.Mutex                // one value is built at a time, so scans is safe
	scans map[string][]branchChange // a branch tip's changes against a target commit, which never change
}

// serve answers the page on addr, a loopback address, until the process is
// stopped.
func serve(stdout io.Writer, addr string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	if host, _, err := net.SplitHostPort(addr); err != nil || !loopback(host) {
		return errors.New("--addr must be a loopback address and port, such as 127.0.0.1:7070")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Fprintf(stdout, "http://%s/\n", ln.Addr())
	s := &server{p: p, scans: map[string][]branchChange{}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(page)
	})
	mux.HandleFunc("GET /api/project", s.project)
	mux.HandleFunc("GET /api/doc", s.doc)
	return http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A site the browser visits can point a name of its own at loopback;
		// answering only loopback names keeps it from reading the project.
		if host, _, _ := net.SplitHostPort(r.Host); !loopback(host) {
			http.Error(w, "the page answers only at a loopback address", http.StatusForbidden)
			return
		}
		mux.ServeHTTP(w, r)
	}))
}

func loopback(host string) bool {
	ip := net.ParseIP(host)
	return host == "localhost" || ip != nil && ip.IsLoopback()
}

// value is what /api/project answers and grove changes --json prints.
type value struct {
	Root     string          `json:"root"`
	Target   string          `json:"target"`
	Commit   string          `json:"commit"`
	Read     time.Time       `json:"read"`
	Specs    []*doc.Document `json:"specs"`
	Changes  []pageChange    `json:"changes"`
	Cycles   []pageCycle     `json:"cycles"`
	Branches []branchRow     `json:"branches"`
	Runs     []pageRun       `json:"runs"`
	History  []event         `json:"history"`
	Problems []string        `json:"problems"`
}

// pageChange is a change as grove changes --json gives it.
type pageChange struct {
	*doc.Change
	Cycle string   `json:"cycle,omitempty"` // the cycle not done that lists it
	Ready bool     `json:"ready"`           // approved and waiting on nothing
	Files []string `json:"files"`           // its documents, from the root
}

type pageCycle struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Goal       string            `json:"goal,omitempty"`
	Created    string            `json:"created"`
	Path       string            `json:"path"`
	Planned    string            `json:"planned"`
	Changes    []string          `json:"changes"`  // its open members in build order, then those waiting outside it
	Archived   []string          `json:"archived"` // its members archived
	DependsOn  []string          `json:"depends_on,omitempty"`
	NeedsHuman map[string]string `json:"needs_human,omitempty"`
	Stale      []string          `json:"stale,omitempty"` // members changed on the target since the plan
	Unattended bool              `json:"unattended"`
	Done       bool              `json:"done"`
	Runs       bool              `json:"runs"`    // grove claim would claim it
	Verdict    string            `json:"verdict"` // what grove claim would do with it, or why not
}

// pageRun is an attempt with its branch's latest commit time, local or on origin.
type pageRun struct {
	attempt
	Moved string `json:"moved,omitempty"`
}

// branchRow is a branch, or the working tree, holding work the target lacks.
type branchRow struct {
	Name       string         `json:"name"`
	Tip        string         `json:"tip"` // worktree for the working tree
	Date       string         `json:"date"`
	Subject    string         `json:"subject"`
	OriginOnly bool           `json:"origin_only"` // no local branch has it
	Pushed     bool           `json:"pushed"`      // origin holds the same tip
	Changes    []branchChange `json:"changes"`
}

type branchChange struct {
	ID    string     `json:"id"`
	Title string     `json:"title"`
	Dir   string     `json:"dir"`
	Draft bool       `json:"draft"` // the target does not hold it; otherwise it approves it
	Tasks []doc.Task `json:"tasks"`
	New   []string   `json:"new"` // its files, from the root, whose blob the target's history never held
}

// event is a commit on the target that did something to the record.
type event struct {
	Commit  string        `json:"commit"`
	Date    string        `json:"date"`
	Subject string        `json:"subject"`
	Marks   []historyMark `json:"marks"`
}

// historyMark is approved, revised or archived for a change, planned or
// re-planned for a cycle, or spec for a maintained spec.
type historyMark struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

func (s *server) project(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	v, err := s.value()
	s.mu.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *server) value() (*value, error) {
	p := s.p
	commit, err := gitx.ResolveBranch(p.Top, p.Target)
	if err != nil {
		return nil, err
	}
	fsys, err := gitx.TreeFS(p.Top, commit, p.Root)
	if err != nil {
		return nil, err
	}
	ip := p.instruct()
	ip.Commit = commit
	if ip.Agreed, err = doc.Load(fsys, ip.Layout); err != nil {
		return nil, err
	}
	c := ip.Agreed
	v := &value{Root: p.Root, Target: p.Target, Commit: commit, Read: time.Now(), Specs: c.Specs, Changes: []pageChange{}} // [] for no changes, as scripts iterate it
	for _, ch := range slices.Concat(c.Changes, c.Archived) {
		row := *ch // runnable below rewrites WaitsOn of a member waiting outside its cycle
		pc := pageChange{Change: &row, Ready: ch.State == "approved" && len(ch.WaitsOn) == 0, Files: ch.Files}
		for _, cy := range c.Cycles {
			if !cy.Done && slices.Contains(cy.Meta.Changes, ch.ID) {
				pc.Cycle = cy.ID
			}
		}
		v.Changes = append(v.Changes, pc)
	}
	tip := tipOf(p.Top, "refs/remotes/origin/grove/runs")
	if tip == "" {
		tip = tipOf(p.Top, runs)
	}
	recs, err := records(p.Top, tip)
	if err != nil {
		return nil, err
	}
	for _, a := range recs {
		moved, err := gitx.Run(p.Top, "for-each-ref", "--sort=-committerdate", "--count=1", "--format=%(committerdate:iso-strict)", "refs/heads/"+a.Branch, "refs/remotes/origin/"+a.Branch)
		if err != nil {
			return nil, err
		}
		v.Runs = append(v.Runs, pageRun{a, strings.TrimSpace(moved)})
	}
	last := pausing(recs, c)
	for _, cy := range c.Cycles {
		pc := pageCycle{ID: cy.ID, Title: cy.Meta.Title, Goal: cy.Meta.Goal, Created: cy.Meta.Created, Path: cy.Path, Planned: cy.Meta.Planned,
			DependsOn: cy.Meta.DependsOn, NeedsHuman: cy.Meta.NeedsHuman, Unattended: cy.Meta.Unattended, Done: cy.Done}
		ordered, waiting := buildable(c, c, cy)
		for _, ch := range slices.Concat(ordered, waiting) {
			pc.Changes = append(pc.Changes, ch.ID)
		}
		pc.Archived = slices.DeleteFunc(slices.Clone(cy.Meta.Changes), func(id string) bool { return !c.InArchive(id) })
		if !cy.Done {
			pc.Stale, _ = changedSincePlan(p, commit, c, cy) // an error is runnable's to give as the verdict
		}
		if _, err := runnable(p, ip, cy); err != nil {
			pc.Verdict = err.Error()
		} else if last != nil && last.Cycle == cy.ID {
			pc.Verdict = fmt.Sprintf("claimed by attempt %s, which %s started at %s on branch %s", last.ID, last.Host, last.Started.Format(time.RFC3339), last.Branch)
		} else {
			cmd, why := "grove claim", ""
			if !cy.Meta.Unattended {
				cmd, why = cmd+" "+cy.ID, "; it is not marked unattended"
			}
			if last != nil {
				cmd, why = cmd+" --retry", why+"; the project is paused after attempt "+last.ID+" of cycle "+last.Cycle
			}
			pc.Runs, pc.Verdict = last == nil, cmd+" would claim it"+why
		}
		v.Cycles = append(v.Cycles, pc)
	}
	if v.Branches, err = s.branches(commit, c); err != nil {
		return nil, err
	}
	if v.History, err = history(p, commit); err != nil {
		return nil, err
	}
	v.Problems, err = diagnose(p, fsys, c)
	return v, err
}

// branches returns each local and origin branch the target at commit does
// not hold merged, and the working tree, with the changes whose files there
// hold a blob the target's history never held, newest first. Telling by
// blob rather than by content leaves out a branch squash-merged and then
// overtaken, whose files differ from the target's but are old.
func (s *server) branches(commit string, c *doc.Collection) ([]branchRow, error) {
	p := s.p
	objects, err := gitx.Run(p.Top, "rev-list", "--objects", commit, "--", p.Root)
	if err != nil {
		return nil, err
	}
	held := map[string]bool{}
	for line := range strings.Lines(objects) {
		sha, _, _ := strings.Cut(strings.TrimSpace(line), " ")
		held[sha] = true
	}
	// Sorted by refname on equal dates, a local branch comes before its origin twin.
	refs, err := gitx.Run(p.Top, "for-each-ref", "--no-merged="+commit, "--sort=refname", "--sort=-committerdate",
		"--format=%(refname)%00%(refname:short)%00%(objectname)%00%(committerdate:iso-strict)%00%(subject)", "refs/heads", "refs/remotes/origin")
	if err != nil {
		return nil, err
	}
	var rows []branchRow
	local := map[string]int{} // a local branch's name to its row
	for line := range strings.Lines(refs) {
		f := strings.Split(strings.TrimSuffix(line, "\n"), "\x00")
		name, remote := strings.CutPrefix(f[0], "refs/remotes/origin/")
		// grove/runs holds records, sharing no history with the target.
		if name == runs || remote && (name == "HEAD" || name == p.Target || name == "grove/runs") {
			continue
		}
		if i, ok := local[name]; remote && ok && rows[i].Tip == f[2] {
			rows[i].Pushed = true
			continue
		}
		if !remote {
			local[strings.TrimPrefix(f[0], "refs/heads/")] = len(rows)
		}
		changes, err := s.scan(f[2], commit, held, c)
		if err != nil {
			return nil, err
		}
		rows = append(rows, branchRow{Name: f[1], Tip: f[2], Date: f[3], Subject: f[4], OriginOnly: remote, Changes: changes})
	}
	status, err := gitx.Run(p.Top, "status", "--porcelain", "-z", "--untracked-files=all", "--no-renames", "--", p.Root)
	if err != nil {
		return nil, err
	}
	var edited []string
	for _, rec := range strings.Split(status, "\x00") {
		if len(rec) > 3 && !strings.Contains(rec[:2], "D") {
			edited = append(edited, rec[3:])
		}
	}
	if len(edited) > 0 {
		sums, err := gitx.Run(p.Top, append([]string{"hash-object", "--"}, edited...)...)
		if err != nil {
			return nil, err
		}
		var fresh []string
		for i, sum := range strings.Fields(sums) {
			if !held[sum] {
				fresh = append(fresh, edited[i])
			}
		}
		changes, err := s.classify(fresh, os.DirFS(filepath.Join(p.Top, p.Root)), c)
		if err != nil {
			return nil, err
		}
		rows = slices.Insert(rows, 0, branchRow{Name: "working tree", Tip: "worktree", Date: time.Now().Format(time.RFC3339), Subject: "uncommitted", Changes: changes})
	}
	return slices.DeleteFunc(rows, func(b branchRow) bool { return len(b.Changes) == 0 }), nil
}

// scan returns the changes of tip with a file whose blob held lacks.
func (s *server) scan(tip, commit string, held map[string]bool, c *doc.Collection) ([]branchChange, error) {
	key := tip + " " + commit
	if changes, ok := s.scans[key]; ok {
		return changes, nil
	}
	raw, err := gitx.Run(s.p.Top, "diff", "--raw", "--no-abbrev", "--no-renames", commit+"..."+tip, "--", s.p.Root)
	if err != nil {
		return nil, err
	}
	var fresh []string
	for line := range strings.Lines(raw) {
		meta, file, _ := strings.Cut(strings.TrimSuffix(line, "\n"), "\t")
		if m := strings.Fields(meta); m[4] != "D" && !held[m[3]] {
			fresh = append(fresh, file)
		}
	}
	var changes []branchChange
	if len(fresh) > 0 {
		fsys, err := gitx.TreeFS(s.p.Top, tip, s.p.Root)
		if err != nil {
			return nil, err
		}
		if changes, err = s.classify(fresh, fsys, c); err != nil {
			return nil, err
		}
	}
	s.scans[key] = changes
	return changes, nil
}

// classify returns the changes in fsys, a branch's root, that a file of
// fresh is in, leaving out those c, the target's, holds archived.
func (s *server) classify(fresh []string, fsys fs.FS, c *doc.Collection) ([]branchChange, error) {
	at, err := doc.Load(fsys, s.p.instruct().Layout)
	if err != nil {
		return nil, err
	}
	var out []branchChange
	for _, ch := range slices.Concat(at.Changes, at.Archived) {
		var files []string
		for _, f := range fresh {
			if rel := strings.TrimPrefix(f, s.p.Root+"/"); strings.HasPrefix(rel, ch.Dir+"/") {
				files = append(files, rel)
			}
		}
		if len(files) > 0 && !c.InArchive(ch.ID) {
			out = append(out, branchChange{ID: ch.ID, Title: ch.Meta.Title, Dir: ch.Dir, Draft: c.Change(ch.ID) == nil, Tasks: ch.Tasks, New: files})
		}
	}
	return out, nil
}

// history returns each commit on the first-parent history of commit,
// newest first, that approved, revised or archived a change, planned or
// re-planned a cycle, or edited a maintained spec, with what it did.
func history(p project, commit string) ([]event, error) {
	log, err := gitx.Run(p.Top, "log", "--first-parent", "--diff-merges=first-parent", "--no-renames", "--name-status", "--format=%x00%H%x01%cI%x01%s", commit)
	if err != nil {
		return nil, err
	}
	var out []event
	for _, block := range strings.Split(log, "\x00")[1:] {
		header, files, _ := strings.Cut(block, "\n")
		h := strings.SplitN(header, "\x01", 3)
		e := event{Commit: h[0], Date: h[1], Subject: h[2]}
		add := func(kind, name string) {
			if m := (historyMark{kind, name[:min(12, len(name))]}); !slices.Contains(e.Marks, m) {
				e.Marks = append(e.Marks, m)
			}
		}
		for line := range strings.Lines(files) {
			status, file, _ := strings.Cut(strings.TrimSpace(line), "\t")
			rel, under := strings.CutPrefix(file, p.Root+"/")
			parts := strings.Split(rel, "/")
			switch {
			case !under:
			case len(parts) == 4 && parts[0] == "changes" && parts[1] == "archive":
				if parts[3] == ".grove.yaml" && status == "A" {
					add("archived", parts[2])
				}
			case len(parts) >= 3 && parts[0] == "changes" && parts[1] != "archive":
				if parts[2] == ".grove.yaml" && status == "A" {
					add("approved", parts[1])
				} else if status != "D" {
					add("revised", parts[1])
				}
			case len(parts) == 2 && parts[0] == "cycles":
				if kind := map[string]string{"A": "planned", "M": "re-planned"}[status]; kind != "" {
					add(kind, parts[1])
				}
			case len(parts) == 3 && parts[0] == "specs" && parts[2] == "spec.md":
				e.Marks = append(e.Marks, historyMark{"spec", parts[1]})
			}
		}
		e.Marks = slices.DeleteFunc(e.Marks, func(m historyMark) bool {
			return m.Kind == "revised" && slices.Contains(e.Marks, historyMark{"approved", m.ID})
		})
		if len(e.Marks) > 0 {
			out = append(out, e)
		}
	}
	return out, nil
}

var commitRef = regexp.MustCompile(`^[0-9a-f]{40}$`)

// markdown renders as Glamour's goldmark does, with GFM tables. Without
// html.WithUnsafe it drops raw HTML, so a teammate's file cannot script
// the page.
var markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

// doc answers one Markdown or YAML file under the root, at the target
// (ref empty), at a commit or in the working tree, Markdown as HTML.
func (s *server) doc(w http.ResponseWriter, r *http.Request) {
	ref, name := r.URL.Query().Get("ref"), r.URL.Query().Get("path")
	if ext := path.Ext(name); !fs.ValidPath(name) || ext != ".md" && ext != ".yaml" || ref != "" && ref != "worktree" && !commitRef.MatchString(ref) {
		http.Error(w, "the page reads a Markdown or YAML file under the root, at the target, a commit or the working tree", http.StatusBadRequest)
		return
	}
	var b []byte
	var err error
	if ref == "worktree" {
		var root *os.Root
		if root, err = os.OpenRoot(filepath.Join(s.p.Top, s.p.Root)); err == nil {
			defer root.Close()
			b, err = fs.ReadFile(root.FS(), name)
		}
	} else {
		if ref == "" {
			ref, err = gitx.ResolveBranch(s.p.Top, s.p.Target)
		}
		var out string
		if err == nil {
			out, err = gitx.Run(s.p.Top, "show", ref+":"+s.p.Root+"/"+name)
		}
		b = []byte(out)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if path.Ext(name) == ".yaml" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(b)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	markdown.Convert(b, w)
}
