package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/mascah/grove/internal/handoff"
)

// runContext writes only a completed bundle: a refused source leaves stdout empty.
func runContext(root string, a invocation, out, errOut io.Writer) int {
	bundle, err := handoff.Build(context.Background(), root, a.ids, a.options)
	if err != nil {
		report(errOut, err)
		return 1
	}
	if !a.json {
		return writeResult(out, errOut, handoff.Text(bundle))
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	encoder.Encode(bundle) // strings, numbers, and booleans cannot fail to encode
	return writeResult(out, errOut, buffer.Bytes())
}
