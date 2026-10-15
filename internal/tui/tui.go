// Package tui is the reader: a list of documents, and the selected one
// rendered as Markdown. It knows nothing about Git or Grove's documents; the
// caller hands it items.
package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Item is one readable entry.
type Item struct {
	Label, ID, Title string
	// Sections are Markdown documents shown one after another, each rendered
	// on its own so that one cannot restyle or hide the next.
	Sections []string
}

// Run shows the reader until the person quits.
func Run(header string, items []Item) error {
	_, err := tea.NewProgram(&model{header: header, items: items}).Run()
	return err
}

type model struct {
	header        string
	items         []Item
	cursor        int
	open          bool     // showing the selected item rather than the list
	lines         []string // the open item, rendered
	top           int      // first line shown of the open item
	width, height int
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.open {
			m.render()
		}
	case tea.KeyPressMsg:
		page := max(m.height-2, 1)
		switch k := msg.String(); {
		case k == "q" || k == "ctrl+c":
			return m, tea.Quit
		case m.open && (k == "esc" || k == "backspace" || k == "left"):
			m.open = false
		case m.open && (k == "down" || k == "j"):
			m.top = min(m.top+1, max(len(m.lines)-page, 0))
		case m.open && (k == "up" || k == "k"):
			m.top = max(m.top-1, 0)
		case m.open && (k == "pgdown" || k == "space"):
			m.top = min(m.top+page, max(len(m.lines)-page, 0))
		case m.open && k == "pgup":
			m.top = max(m.top-page, 0)
		case !m.open && (k == "down" || k == "j"):
			m.cursor = min(m.cursor+1, len(m.items)-1)
		case !m.open && (k == "up" || k == "k"):
			m.cursor = max(m.cursor-1, 0)
		case !m.open && (k == "enter" || k == "right") && len(m.items) > 0:
			m.open, m.top = true, 0
			m.render()
		}
	}
	return m, nil
}

func (m *model) render() {
	var parts []string
	for _, s := range m.items[m.cursor].Sections {
		parts = append(parts, Markdown(s, m.width-2))
	}
	m.lines = strings.Split(strings.Join(parts, "\n"+strings.Repeat("─", max(m.width, 1))+"\n"), "\n")
}

func (m *model) View() tea.View {
	var b strings.Builder
	if m.open {
		it := m.items[m.cursor]
		fmt.Fprintf(&b, "%s  %s  %s\n", it.Label, Sanitize(it.ID), Sanitize(it.Title))
		end := min(m.top+max(m.height-2, 1), len(m.lines))
		b.WriteString(strings.Join(m.lines[m.top:end], "\n"))
		b.WriteString("\n\x1b[2m↑↓ scroll · esc back · q quit\x1b[m")
	} else {
		fmt.Fprintf(&b, "%s\n\n", Sanitize(m.header))
		for i, it := range m.items {
			mark := "  "
			if i == m.cursor {
				mark = "> "
			}
			fmt.Fprintf(&b, "%s%-16s %s  %s\n", mark, it.Label, Sanitize(it.ID), Sanitize(it.Title))
		}
		if len(m.items) == 0 {
			b.WriteString("Nothing is agreed or approved yet.\n")
		}
		b.WriteString("\n\x1b[2m↑↓ select · enter read · q quit\x1b[m")
	}
	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
