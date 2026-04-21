// Package yellowstone provides runtime monitoring helpers for annotated
// functions. Monitoring is off by default and enabled via the
// YELLOWSTONE_MONITOR env var (truthy: 1, true, yes, on).
package yellowstone

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
)

var (
	enableOnce sync.Once
	enableVal  bool
)

func enabled() bool {
	enableOnce.Do(func() {
		switch strings.ToLower(os.Getenv("YELLOWSTONE_MONITOR")) {
		case "1", "true", "yes", "on":
			enableVal = true
		}
	})
	return enableVal
}

// Monitor records the inputs to a function immediately and returns a closure
// that, when invoked (typically via defer), records the outputs. Outputs must
// be passed as pointers to the function's named return values so the
// post-execution values are observed.
//
// Typical usage:
//
//	func (vm *VM) pop() (v int64, err error) {
//	    defer yellowstone.Monitor("VM.pop", []any{vm}, &v, &err)()
//	    // ...
//	}
//
// When monitoring is disabled the returned closure is a no-op.
func Monitor(name string, inputs []any, outputs ...any) func() {
	if !enabled() {
		return func() {}
	}
	fmt.Fprintf(os.Stderr, "[🦬 →] %s inputs=%s\n", name, format(inputs))
	return func() {
		derefed := make([]any, len(outputs))
		for i, o := range outputs {
			derefed[i] = deref(o)
		}
		fmt.Fprintf(os.Stderr, "[🦬 ←] %s outputs=%s\n", name, format(derefed))
	}
}

func deref(p any) any {
	if p == nil {
		return nil
	}
	rv := reflect.ValueOf(p)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		return rv.Elem().Interface()
	}
	return p
}

func format(vs []any) string {
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = fmt.Sprintf("%+v", v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
