package attempt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

// ActivityWindow is how much of the end of events.jsonl ReadActivity reads:
// enough that a final result line of MaxLine always fits whole.
const ActivityWindow = MaxLine + 64<<10

// maxActivity bounds the entries kept, and the bytes of each.
const maxActivity, maxActivityLine = 200, 300

// Activity is what the end of an attempt's events shows: one short entry per
// thing that happened, oldest first, and the result event's text when it is
// in the window. The text is the provider's, unescaped; whoever shows it
// escapes it.
type Activity struct {
	Entries []Entry `json:"entries"`
	Report  string  `json:"report,omitempty"`
	Model   string  `json:"model,omitempty"` // as the provider reported it at the start
	Cut     bool    `json:"cut"`             // the window began inside the file: earlier events are not shown or counted
	Metrics Metrics `json:"metrics"`
}

// Entry is one row of activity. Consecutive identical rows are one entry
// with their count, so a provider's repeated notices take one row.
type Entry struct {
	Time  time.Time `json:"time,omitzero"` // the first event's own timestamp; zero when the provider gave none
	Kind  string    `json:"kind"`          // start, text, tool, error, notice or result
	Text  string    `json:"text"`
	Count int       `json:"count"`
}

// Metrics is what a run has used, in terms any provider's reader can fill,
// so a display depends on these fields and not on one provider's events.
// Counts come from the window read, and are lower bounds when the Activity
// is Cut. Total says the tokens are the result event's totals for the run;
// without it the input is summed over the window's messages and the output
// is unknown, since Claude reports a message's output tokens before writing
// it. A run resumed in one session has several result events, and each
// counts turns for its own query only, so Turns and Subagents are the larger
// of the result's figure and the window's count. Context and Window are 0
// when unknown.
type Metrics struct {
	Turns        int  `json:"turns"`
	InputTokens  int  `json:"input_tokens"`   // cached input included
	OutputTokens int  `json:"output_tokens"`  // 0 without Total
	Context      int  `json:"context"`        // tokens in the context at the latest top-level message
	Window       int  `json:"context_window"` // the model's context size
	Subagents    int  `json:"subagents"`
	Compactions  int  `json:"compactions"`
	Tools        int  `json:"tool_calls"`
	ToolErrors   int  `json:"tool_errors"`
	Total        bool `json:"total"`
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
	c := counter{in: map[string]int{}, turns: map[string]bool{}}
	for l := range bytes.SplitSeq(data, []byte("\n")) {
		c.note(&a, l)
	}
	m := &a.Metrics
	if !m.Total {
		for _, n := range c.in {
			m.InputTokens += n
		}
	}
	m.Turns, m.Subagents = max(len(c.turns), c.queries), max(c.agents, c.spawned)
	if len(a.Entries) > maxActivity {
		a.Entries, a.Cut = a.Entries[len(a.Entries)-maxActivity:], true
	}
	return a, nil
}

// counter is what summing a run's messages needs while reading: a message
// arrives as one event per content block, each repeating its usage.
type counter struct {
	in               map[string]int  // input tokens by message
	turns            map[string]bool // top-level messages
	agents           int             // subagents started
	spawned, queries int             // the result events' own counts, the largest seen
}

type usage struct {
	Input       int `json:"input_tokens"`
	CacheRead   int `json:"cache_read_input_tokens"`
	CacheCreate int `json:"cache_creation_input_tokens"`
	Output      int `json:"output_tokens"`
}

// note reads one Claude stream-json event into a's entries and metrics.
func (c *counter) note(a *Activity, line []byte) {
	var ev struct {
		Type      string  `json:"type"`
		Subtype   string  `json:"subtype"`
		Model     string  `json:"model"`
		TaskType  string  `json:"task_type"`
		IsError   bool    `json:"is_error"`
		Result    string  `json:"result"`
		Timestamp string  `json:"timestamp"`
		Parent    *string `json:"parent_tool_use_id"`
		Turns     int     `json:"num_turns"`
		Message   struct {
			ID      string `json:"id"`
			Usage   *usage `json:"usage"`
			Content []struct {
				Type    string          `json:"type"`
				Text    string          `json:"text"`
				Name    string          `json:"name"`
				Input   json.RawMessage `json:"input"`
				IsError bool            `json:"is_error"`
			} `json:"content"`
		} `json:"message"`
		ModelUsage map[string]struct {
			Input       int `json:"inputTokens"`
			Output      int `json:"outputTokens"`
			CacheRead   int `json:"cacheReadInputTokens"`
			CacheCreate int `json:"cacheCreationInputTokens"`
			Window      int `json:"contextWindow"`
		} `json:"modelUsage"`
		Subagents *struct {
			Spawned int `json:"spawned"`
		} `json:"subagent_stats"`
	}
	if json.Unmarshal(line, &ev) != nil {
		return
	}
	at, _ := time.Parse(time.RFC3339Nano, ev.Timestamp) // zero when absent or unreadable
	add := func(kind, text string) {
		text = clip(text)
		if n := len(a.Entries); n > 0 && a.Entries[n-1].Kind == kind && a.Entries[n-1].Text == text {
			a.Entries[n-1].Count++
			return
		}
		a.Entries = append(a.Entries, Entry{Time: at, Kind: kind, Text: text, Count: 1})
	}
	m := &a.Metrics
	switch ev.Type {
	case "system":
		switch ev.Subtype {
		case "init":
			a.Model = ev.Model
			add("start", "started: model "+ev.Model)
			return
		case "compact_boundary":
			m.Compactions++
		case "task_started":
			if ev.TaskType == "local_agent" {
				c.agents++
			}
		}
		add("notice", "system: "+ev.Subtype)
	case "assistant":
		if u := ev.Message.Usage; u != nil && ev.Message.ID != "" {
			c.in[ev.Message.ID] = u.Input + u.CacheRead + u.CacheCreate
			if ev.Parent == nil {
				m.Context = u.Input + u.CacheRead + u.CacheCreate
				c.turns[ev.Message.ID] = true
			}
		}
		for _, b := range ev.Message.Content {
			switch b.Type {
			case "text":
				if first, _, _ := strings.Cut(strings.TrimSpace(b.Text), "\n"); first != "" {
					add("text", first)
				}
			case "tool_use":
				m.Tools++
				add("tool", b.Name+toolHint(b.Input))
			}
		}
	case "user":
		for _, b := range ev.Message.Content {
			if b.Type == "tool_result" && b.IsError {
				m.ToolErrors++
				add("error", "tool error")
			}
		}
	case "result":
		a.Report = ev.Result
		text := "result: " + ev.Subtype
		if ev.IsError {
			text += " (is_error)"
		}
		add("result", text)
		if len(ev.ModelUsage) > 0 {
			m.Total, m.InputTokens, m.OutputTokens = true, 0, 0
			for _, u := range ev.ModelUsage {
				m.InputTokens += u.Input + u.CacheRead + u.CacheCreate
				m.OutputTokens += u.Output
				m.Window = max(m.Window, u.Window)
			}
		}
		c.queries = max(c.queries, ev.Turns)
		if ev.Subagents != nil {
			c.spawned = max(c.spawned, ev.Subagents.Spawned)
		}
	}
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
