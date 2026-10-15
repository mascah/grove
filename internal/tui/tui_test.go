package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestReadAnItem(t *testing.T) {
	m := &model{header: "main@1a2b3c4", items: []Item{
		{Label: "maintained spec", ID: "record-editing", Title: "Record editing", Sections: []string{"# Record editing\n\nWhen a user edits a record, the system shall save it."}},
		{Label: "approved change", ID: "261014-k3m9q", Title: "Allow editing\x1b]0;x\x07 older records", Sections: []string{"one", "two"}},
	}}
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	if v := m.View().Content; !strings.Contains(v, "> maintained spec  record-editing  Record editing") || !strings.Contains(v, `Allow editing\x1b]0;x\a older records`) {
		t.Errorf("list:\n%s", v)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if v := m.View().Content; !strings.Contains(v, "the system shall save it") {
		t.Errorf("open item:\n%s", v)
	}
	m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if v := m.View().Content; !strings.Contains(v, "one") || !strings.Contains(v, "────") || !strings.Contains(v, "two") {
		t.Errorf("sections:\n%s", v)
	}
}
