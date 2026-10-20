package cli

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"
)

// version is set by the release build with -ldflags -X.
var version string

func printVersion(w io.Writer) {
	v := version
	// go install of a tagged module records the tag here instead.
	if info, ok := debug.ReadBuildInfo(); v == "" && ok {
		v = strings.TrimPrefix(info.Main.Version, "v")
	}
	if v == "" || v == "(devel)" {
		v = "dev"
	}
	fmt.Fprintf(w, "grove %s\n", v)
}
