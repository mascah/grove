package cli

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

func instructionsExplore(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Explore(ip))
}

func instructionsPrototype(stdout io.Writer, question string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Prototype(ip, question))
}

// instructionsPropose prints the propose instructions: for a new change in
// change schema schema, or grove.yaml's when it is empty, when arg is empty
// or a description, or for revising the drafted change whose ID arg is.
func instructionsPropose(stdout io.Writer, schema, arg string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	if !doc.ValidID(arg) {
		now := time.Now()
		return show(stdout)(instruct.Propose(ip, cmp.Or(schema, p.Schema), doc.NewID(now), now.UTC().Format(time.RFC3339), arg))
	}
	if schema != "" {
		return errors.New("a change keeps its schema; to change it, edit schema in the change's .grove.yaml")
	}
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout, arg)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Revise(ip, c, files))
}

// instructionsPlan prints the plan instructions: for a new cycle when arg
// is empty or a description, for a new cycle of one when it is an approved
// change's ID, or for revising the cycle in this working tree whose ID it is.
func instructionsPlan(stdout io.Writer, arg string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	now := time.Now()
	if !doc.ValidID(arg) {
		return show(stdout)(instruct.Plan(ip, doc.NewID(now), now.UTC().Format(time.RFC3339), arg, nil))
	}
	if c := ip.Agreed.Change(arg); c != nil {
		return show(stdout)(instruct.Plan(ip, doc.NewID(now), now.UTC().Format(time.RFC3339), "", c))
	}
	fsys := os.DirFS(filepath.Join(p.Top, p.Root))
	tree, err := doc.Load(fsys, ip.Layout)
	if err != nil {
		return err
	}
	cy := tree.Cycle(arg)
	if cy == nil {
		return fmt.Errorf("%s approves no change %s, and this working tree holds no cycle %s", p.Target, arg, arg)
	}
	b, err := fs.ReadFile(fsys, cy.Path)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Replan(ip, cy, string(b)))
}

// instructionsArchive prints the archive instructions for change id as the
// working tree holds it.
func instructionsArchive(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout, id)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Archive(ip, c, files))
}

// target returns the project, the target's documents and this working
// tree's.
func target() (project, instruct.Project, *doc.Collection, error) {
	p, err := openProject()
	if err != nil {
		return p, instruct.Project{}, nil, err
	}
	ip, err := p.agreed()
	if err != nil {
		return p, ip, nil, err
	}
	tree, err := doc.Load(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout)
	return p, ip, tree, err
}

// instructionsApply prints the apply instructions for the target's cycle
// ids[0], or, with ids[1], for that member of it. Apply takes nothing else.
func instructionsApply(stdout io.Writer, ids ...string) error {
	if len(ids) == 0 {
		return errors.New("apply takes a cycle ID; build a change with the owner with /grove-implement <change-id>")
	}
	p, ip, tree, err := target()
	if err != nil {
		return err
	}
	cy := ip.Agreed.Cycle(ids[0])
	switch {
	case cy == nil && ip.Agreed.Change(ids[0]) != nil:
		return fmt.Errorf("apply takes a cycle ID; build %s with the owner with /grove-implement %s", ids[0], ids[0])
	case cy == nil:
		return fmt.Errorf("%s holds no cycle %s", p.Target, ids[0])
	case len(ids) == 1:
		return instructionsCycle(stdout, p, ip, tree, cy, "apply")
	case !slices.Contains(cy.Meta.Changes, ids[1]):
		return fmt.Errorf("%s is not in cycle %s", ids[1], cy.ID)
	}
	return instructionsChange(stdout, p, ip, tree, ids[1], instruct.Apply, cy)
}

// instructionsImplement prints the implement instructions for the target's
// cycle or change id.
func instructionsImplement(stdout io.Writer, id string) error {
	p, ip, tree, err := target()
	if err != nil {
		return err
	}
	if cy := ip.Agreed.Cycle(id); cy != nil {
		return instructionsCycle(stdout, p, ip, tree, cy, "implement")
	}
	return instructionsChange(stdout, p, ip, tree, id, instruct.Implement, nil)
}

// instructionsChange prints the instructions skill builds, apply's or
// implement's, for change id as the target holds it, with cycle cy's plan
// when it is not nil, unless it waits on a change tree has not archived as
// finished (D-56). For a change in no cycle that is not done, they name
// each task's file that changed on the target since the change's base.
func instructionsChange(stdout io.Writer, p project, ip instruct.Project, tree *doc.Collection, id string,
	skill func(instruct.Project, *doc.Change, map[string]string) (instruct.Instructions, error), cy *doc.Cycle) error {
	c := ip.Agreed.Change(id)
	if c == nil {
		return fmt.Errorf("%s holds no approved change %s", p.Target, id)
	}
	if waits := slices.DeleteFunc(slices.Clone(c.Meta.DependsOn), tree.Finished); len(waits) > 0 {
		return fmt.Errorf("%s waits on %s, which this working tree has not archived as finished", id, strings.Join(waits, " and "))
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	c, files, err := change(fsys, ip.Layout, id)
	if err != nil {
		return err
	}
	in, err := skill(ip, c, files)
	if err == nil && cy == nil && !inOpenCycle(tree, id) && c.Meta.Base != "" {
		var tasks []doc.Task
		tasks, err = moved(p, ip.Commit, c)
		var named []string
		for _, t := range tasks {
			for _, f := range t.Files {
				named = append(named, fmt.Sprintf("%s's `%s`", t.ID, f))
			}
		}
		if len(named) > 0 {
			in.Requires += fmt.Sprintf("- Changed on %s since base %s: %s.\n", p.Target, c.Meta.Base[:min(7, len(c.Meta.Base))], strings.Join(named, ", "))
		}
	}
	if err == nil && cy != nil {
		var b []byte
		if b, err = fs.ReadFile(fsys, cy.Path); err == nil {
			// Every member in build order, archived or not, so a member's place does not move.
			outside := func(id string) bool { return tree.Finished(id) && !slices.Contains(cy.Meta.Changes, id) }
			ordered, _ := doc.Order(members(ip.Agreed, cy, func(string) bool { return false }), outside)
			in.InCycle(cy, string(b), ordered, id, tree.InArchive)
		}
	}
	return show(stdout)(in, err)
}

// instructionsCycle prints the instructions for building each member of
// cycle cy that tree has not archived with skill, apply's or implement's,
// in build order.
func instructionsCycle(stdout io.Writer, p project, ip instruct.Project, tree *doc.Collection, cy *doc.Cycle, skill string) error {
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	b, err := fs.ReadFile(fsys, cy.Path)
	if err != nil {
		return err
	}
	ordered, waiting := buildable(ip.Agreed, tree, cy)
	return show(stdout)(instruct.Cycle(ip, skill, cy, string(b), ordered, waiting))
}

// members returns the changes agreed approves that cycle cy lists and that
// done does not report.
func members(agreed *doc.Collection, cy *doc.Cycle, done func(id string) bool) []*doc.Change {
	var out []*doc.Change
	for _, ch := range agreed.Changes {
		if slices.Contains(cy.Meta.Changes, ch.ID) && !done(ch.ID) {
			out = append(out, ch)
		}
	}
	return out
}

// buildable returns the members of cycle cy that agreed approves and tree
// has not archived, in build order, then those that wait on a change
// neither being built nor finished in tree.
func buildable(agreed, tree *doc.Collection, cy *doc.Cycle) (ordered, waiting []*doc.Change) {
	ordered, waiting = doc.Order(members(agreed, cy, tree.InArchive), tree.Finished)
	for _, ch := range waiting { // what keeps it out: neither being built nor finished here
		ch.WaitsOn = slices.DeleteFunc(slices.Clone(ch.Meta.DependsOn), func(id string) bool {
			return tree.Finished(id) || slices.ContainsFunc(ordered, func(o *doc.Change) bool { return o.ID == id })
		})
	}
	return ordered, waiting
}

// show returns a function printing instructions unless building them failed.
func show(stdout io.Writer) func(instruct.Instructions, error) error {
	return func(in instruct.Instructions, err error) error {
		if err == nil {
			fmt.Fprint(stdout, in.Markdown())
		}
		return err
	}
}

// change loads change id from the document root fsys, with every file it
// holds.
func change(fsys fs.FS, layout doc.Layout, id string) (*doc.Change, map[string]string, error) {
	all, err := doc.Load(fsys, layout)
	if err != nil {
		return nil, nil, err
	}
	c := all.Change(id)
	if c == nil {
		return nil, nil, fmt.Errorf("no change %s in this working tree", id)
	}
	files := map[string]string{}
	for _, name := range c.Files {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, nil, err
		}
		files[name] = string(b)
	}
	return c, files, nil
}
