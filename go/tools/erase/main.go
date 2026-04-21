// Command erase strips yellowstone monitoring annotations from a source tree
// and writes the result to an output directory mirroring the input layout.
// Go source files and go.mod files are stripped; other files are copied
// verbatim.
//
// Two markers are recognized:
//
//   - `// 🦬` at end of a line: the entire line is removed.
//   - `// 🦬: <text>` on its own line: the marker line is removed and the
//     following line is replaced verbatim by <text>.
//
// The markers must use the exact byte sequence "// 🦬" (space + bison emoji).
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	deleteMarker  = "// 🦬"
	replaceMarker = "// 🦬: "
)

func strip(src []byte) []byte {
	lines := strings.Split(string(src), "\n")
	out := make([]string, 0, len(lines))
	var pending string
	var hasPending bool
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, replaceMarker) {
			pending = strings.TrimPrefix(trimmed, replaceMarker)
			hasPending = true
			continue
		}
		if hasPending {
			out = append(out, pending)
			hasPending = false
			continue
		}
		if strings.HasSuffix(trimmed, deleteMarker) {
			continue
		}
		out = append(out, line)
	}
	return []byte(strings.Join(out, "\n"))
}

func shouldStrip(path string) bool {
	return strings.HasSuffix(path, ".go") || filepath.Base(path) == "go.mod"
}

func main() {
	src := flag.String("src", "", "annotated source directory")
	out := flag.String("out", "", "destination directory for stripped files")
	flag.Parse()
	if *src == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "usage: erase -src <annotated-dir> -out <destination-dir>")
		os.Exit(2)
	}
	err := filepath.WalkDir(*src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(*src, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result := data
		if shouldStrip(path) {
			result = strip(data)
		}
		target := filepath.Join(*out, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, result, 0o644)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
