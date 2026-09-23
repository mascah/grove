package attempt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

// ActivityWindow is how much of the end of events.jsonl ReadActivity reads:
// enough that a final result line of MaxLine always fits whole.
const ActivityWindow = MaxLine + 64<<10

// maxActivity bounds the lines kept, and the bytes of each.
const maxActivity, maxActivityLine = 200, 300

// Activity is what the end of an attempt's events shows: one short line per
// event, oldest first, and the result event's text when it is in the window.
// The text is the provider's, unescaped; whoever shows it escapes it.
type Activity struct {
	Lines  []string `json:"lines"`
	Report string   `json:"report,omitempty"`
	Cut    bool     `json:"cut"` // the window began inside the file: earlier events are not shown
}

// ReadActivity reads at most the last window bytes of path, so its cost is
// bounded however much the provider has written. A leading line the window
// cut into, and a final line still being written, are skipped.
func ReadActivity(path string, window int64) (Activity, error) {
	var a Activity
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return a, nil
	}
	if err != nil {
		return a, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return a, err
	}
	from := max(info.Size()-window, 0)
	data, err := io.ReadAll(io.NewSectionReader(f, from, info.Size()-from))
	if err != nil {
		return a, err
	}
	if from > 0 {
		a.Cut = true
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			data = data[i+1:]
		} else {
			data = nil
		}
	}
	if i := bytes.LastIndexByte(data, '\n'); i >= 0 {
		data = data[:i]
	} else {
		data = nil
	}
	for l := range bytes.SplitSeq(data, []byte("\n")) {
		for _, text := range a.note(l) {
			a.Lines = append(a.Lines, clip(text))
		}
	}
	if len(a.Lines) > maxActivity {
		a.Lines, a.Cut = a.Lines[len(a.Lines)-maxActivity:], true
	}
	return a, nil
}

// note summarizes one event as lines; an event without anything to show
// gives none.
func (a *Activity) note(line []byte) []string {
	var ev struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Model   string `json:"model"`
		IsError bool   `json:"is_error"`
		Result  string `json:"result"`
		Message struct {
			Content []struct {
				Type    string          `json:"type"`
				Text    string          `json:"text"`
				Name    string          `json:"name"`
				Input   json.RawMessage `json:"input"`
				IsError bool            `json:"is_error"`
			} `json:"content"`
		} `json:"message"`
	}
	if json.Unmarshal(line, &ev) != nil {
		return nil
	}
	var out []string
	switch ev.Type {
	case "system":
		if ev.Subtype == "init" {
			return []string{"started: model " + ev.Model}
		}
		return []string{"system: " + ev.Subtype}
	case "assistant":
		for _, c := range ev.Message.Content {
			switch c.Type {
			case "text":
				if first, _, _ := strings.Cut(strings.TrimSpace(c.Text), "\n"); first != "" {
					out = append(out, first)
				}
			case "tool_use":
				out = append(out, "tool: "+c.Name+toolHint(c.Input))
			}
		}
	case "user":
		for _, c := range ev.Message.Content {
			if c.Type == "tool_result" && c.IsError {
				out = append(out, "tool error")
			}
		}
	case "result":
		a.Report = ev.Result
		text := "result: " + ev.Subtype
		if ev.IsError {
			text += " (is_error)"
		}
		out = append(out, text)
	}
	return out
}

// toolHint is the one input field that says what a tool call is about.
func toolHint(input json.RawMessage) string {
	var fields map[string]any
	if json.Unmarshal(input, &fields) != nil {
		return ""
	}
	for _, k := range []string{"command", "file_path", "pattern", "description", "skill"} {
		if s, ok := fields[k].(string); ok && s != "" {
			first, _, _ := strings.Cut(s, "\n")
			return " " + first
		}
	}
	return ""
}

func clip(s string) string {
	if len(s) <= maxActivityLine {
		return s
	}
	return strings.ToValidUTF8(s[:maxActivityLine], "") + "…"
}
