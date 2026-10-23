package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

const defaultTimeout = 2 * time.Hour

// stopGrace is how long a stopped session's processes have to exit before
// they are killed.
const stopGrace = 5 * time.Second

// attempt is the record of one run, kept as <ID>.json in .git/grove/.
type attempt struct {
	ID       string    `json:"id"`
	Started  time.Time `json:"started"`
	Ended    time.Time `json:"ended,omitzero"`
	Outcome  string    `json:"outcome,omitempty"` // exit status N, timed out, cancelled or failed to start
	Cycle    string    `json:"cycle"`
	Changes  []string  `json:"changes"` // the cycle's, in build order
	Commit   string    `json:"commit"`
	Checkout string    `json:"checkout"`
	Branch   string    `json:"branch"`
	Output   string    `json:"output"`
}

// runAttempt starts one attempt: it updates the target from origin, if the
// checkout has one, and,
// unless the latest attempt's changes are still open or retry is set,
// launches one Claude Code session on the apply skill for the oldest
// eligible cycle, on a new branch, within timeout.
func runAttempt(stdout io.Writer, retry bool, timeout time.Duration) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	gitDir, err := gitx.Run(p.Top, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return err
	}
	dir := filepath.Join(strings.TrimSpace(gitDir), "grove")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	lock, err := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer lock.Close() // the lock goes with the file, so a killed run leaves none
	if syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
		return errors.New("a run is in progress in this checkout")
	}
	if out, err := gitx.Run(p.Top, "status", "--porcelain"); err != nil {
		return err
	} else if out != "" {
		return errors.New("the checkout has uncommitted changes; commit, stash or discard them first")
	}
	local, err := gitx.ResolveBranch(p.Top, p.Target)
	if err != nil {
		return err
	}
	remotes, err := gitx.Run(p.Top, "remote")
	if err != nil {
		return err
	}
	commit := local // with no origin, the local target is the fetched one
	if slices.Contains(strings.Fields(remotes), "origin") {
		if _, err := gitx.Run(p.Top, "fetch", "-q", "origin", p.Target); err != nil {
			return err
		}
		if commit, err = gitx.Run(p.Top, "rev-parse", "--verify", "FETCH_HEAD^{commit}"); err != nil {
			return err
		}
		commit = strings.TrimSpace(commit)
		if _, err := gitx.Run(p.Top, "merge-base", "--is-ancestor", local, commit); err != nil {
			return fmt.Errorf("the local %s has commits origin's does not, so it cannot be fast-forwarded", p.Target)
		}
	}
	config, err := gitx.Run(p.Top, "show", commit+":grove.yaml")
	if err != nil {
		return err
	}
	if p, err = parseProject(p.Top, []byte(config)); err != nil {
		return err
	}
	skill := map[string]string{"skills": "/grove-apply", "commands": "/grove:apply"}[p.Agents["claude"]]
	if skill == "" {
		return errors.New("run needs Claude Code listed under agents in grove.yaml, as skills or commands")
	}
	ip, err := p.at(commit)
	if err != nil {
		return err
	}
	last, err := latest(dir)
	if err != nil {
		return err
	}
	for _, id := range last.Changes {
		if !retry && !ip.Agreed.InArchive(id) {
			return fmt.Errorf("the project is paused after attempt %s, whose output is %s; once you have dealt with it, run grove run --retry", last.ID, last.Output)
		}
	}
	cy, ordered := eligible(p, ip)
	if cy == nil {
		fmt.Fprintf(stdout, "Nothing to apply on %s at %s.\n", p.Target, commit[:7])
		return nil
	}

	now := time.Now()
	a := attempt{ID: doc.NewID(now), Started: now, Cycle: cy.ID, Commit: commit, Checkout: p.Top}
	a.Branch, a.Output = "grove/run-"+a.ID, filepath.Join(dir, a.ID+".log")
	for _, ch := range ordered {
		a.Changes = append(a.Changes, ch.ID)
	}
	if _, err := gitx.Run(p.Top, "switch", "-q", "-c", a.Branch, commit); err != nil {
		return err
	}
	// The apply instructions read the local target, so it moves too.
	if _, err := gitx.Run(p.Top, "update-ref", "refs/heads/"+p.Target, commit, local); err != nil {
		return err
	}
	if err := a.save(dir); err != nil {
		return err
	}
	out, err := os.Create(a.Output)
	if err != nil {
		return err
	}
	defer out.Close()
	a.Outcome = session(stdout, out, p.Top, skill+" "+cy.ID, timeout)
	a.Ended = time.Now()
	if err := a.save(dir); err != nil {
		return err
	}
	ended := fmt.Sprintf("%s ended: %s; its output is %s", a.ID, a.Outcome, a.Output)
	if a.Outcome != "exit status 0" {
		return errors.New("attempt " + ended)
	}
	fmt.Fprintf(stdout, "Attempt %s.\n", ended)
	return nil
}

// eligible returns the oldest cycle the target ip holds that is marked
// unattended, is not done, names no member in needs_human, has no member
// waiting on a change neither in it nor archived as finished, and whose
// members are as its plan read them; and its members in build order.
func eligible(p project, ip instruct.Project) (*doc.Cycle, []*doc.Change) {
	for _, cy := range ip.Agreed.Cycles {
		ordered, waiting := buildable(ip.Agreed, ip.Agreed, cy)
		if !cy.Meta.Unattended || len(ordered) == 0 || len(waiting) > 0 ||
			slices.ContainsFunc(cy.Meta.Changes, func(id string) bool { _, ok := cy.Meta.NeedsHuman[id]; return ok }) {
			continue
		}
		// A planned that is no commit is check's to report; the cycle is skipped.
		if changed, err := changedSincePlan(p, ip.Commit, ip.Agreed, cy); err == nil && len(changed) == 0 {
			return cy, ordered
		}
	}
	return nil, nil
}

// session runs skill in Claude Code in dir, writing every event it emits
// to out and printing a line to stdout for each message, tool call and
// result, and returns how it ended. Whichever way it ends, its process
// group is stopped.
func session(stdout io.Writer, out *os.File, dir, skill string, timeout time.Duration) string {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(out, err)
		return "failed to start"
	}
	defer r.Close()
	cmd := exec.Command("claude", "-p", skill, "--output-format", "stream-json", "--verbose")
	// A file, not a writer, so Wait returns when claude exits, not when the last child closes it.
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, w, out
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // its own group, so it can be stopped whole
	err = cmd.Start()
	w.Close()
	if err != nil {
		fmt.Fprintln(out, err)
		return "failed to start"
	}
	shown := make(chan struct{})
	go func() { relay(stdout, out, r); close(shown) }()
	done := make(chan struct{})
	go func() { cmd.Wait(); close(done) }()
	outcome := ""
	select {
	case <-done:
		outcome = cmd.ProcessState.String()
	case <-time.After(timeout):
		outcome = "timed out"
	case <-signals:
		outcome = "cancelled"
	}
	stop(cmd.Process.Pid)
	<-done
	r.SetReadDeadline(time.Now().Add(stopGrace)) // a process that left the group may hold the pipe open
	<-shown
	return outcome
}

// relay copies each line of r, one stream-json event, to out, and prints to
// stdout the first line of each assistant message, each tool call with its
// input, and each tool result.
func relay(stdout, out io.Writer, r io.Reader) {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')
		out.Write(line)
		if err != nil {
			return
		}
		var e struct {
			Type    string
			Message struct {
				Content []struct {
					Type, Text, Name string
					Input            map[string]any
					Content          json.RawMessage
				}
			}
		}
		if json.Unmarshal(line, &e) != nil {
			continue
		}
		for _, c := range e.Message.Content {
			switch {
			case e.Type == "assistant" && c.Type == "text":
				fmt.Fprintln(stdout, firstLine(c.Text))
			case e.Type == "assistant" && c.Type == "tool_use":
				input := ""
				// ponytail: the usual tools' main argument; add a key when a tool shows only its name.
				for _, k := range []string{"command", "file_path", "pattern", "url", "description"} {
					if v, ok := c.Input[k].(string); ok {
						input = v
						break
					}
				}
				fmt.Fprintln(stdout, c.Name+": "+firstLine(input))
			case e.Type == "user" && c.Type == "tool_result":
				var text string
				if json.Unmarshal(c.Content, &text) != nil {
					var parts []struct{ Text string }
					json.Unmarshal(c.Content, &parts)
					for _, p := range parts {
						text += p.Text
					}
				}
				fmt.Fprintln(stdout, "  "+firstLine(text))
			}
		}
	}
}

func firstLine(s string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(s), "\n")
	return line
}

// stop asks every process in group pgid to exit, kills those left after
// stopGrace, and returns once none is left.
func stop(pgid int) {
	syscall.Kill(-pgid, syscall.SIGTERM)
	kill := time.Now().Add(stopGrace)
	for syscall.Kill(-pgid, 0) == nil {
		if time.Now().After(kill) {
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (a attempt) save(dir string) error {
	b, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, a.ID+".json.tmp")
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, a.ID+".json")) // whole, even if the run is killed while saving
}

// latest returns the most recently started attempt in dir, or an empty one.
func latest(dir string) (attempt, error) {
	var last attempt
	names, err := filepath.Glob(filepath.Join(dir, "*.json"))
	for _, name := range names {
		b, err := os.ReadFile(name)
		if err != nil {
			return last, err
		}
		var a attempt
		if err := json.Unmarshal(b, &a); err != nil {
			return last, fmt.Errorf("%s: %w", name, err)
		}
		if a.Started.After(last.Started) {
			last = a
		}
	}
	return last, err
}
