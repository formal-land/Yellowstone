# Conventions

Working agreements for this project. Keep this file up to date as decisions are made.

## Monitoring annotations in `go/example/`

Monitoring is added **directly in the source** with lightweight `// 🦬` markers. A companion tool strips the markers to recover the bare, unmonitored code. There is no parallel tree.

### Layout

```
go/
├── example/
│   ├── go.mod            # require+replace yellowstone (both marked // 🦬)
│   ├── main.go
│   ├── instructions.go
│   └── vm.go             # contains annotations
├── yellowstone/          # monitoring helper module
│   ├── go.mod
│   └── yellowstone.go
└── tools/erase/          # strips annotations
    ├── go.mod
    └── main.go
```

Three independent Go modules. `go/example/` depends on `go/yellowstone/` via a local `replace` directive in `example/go.mod`. Both the `require yellowstone` and `replace yellowstone => ../yellowstone` lines carry `// 🦬` so `erase` removes them from the stripped tree (which has no yellowstone import anymore).

### Monitoring helper

`yellowstone.Monitor(name string, inputs []any, outputs ...any) func()` logs inputs immediately and returns a closure that logs outputs when called. Pair it with `defer` and named return values:

```go
func (vm *VM) pop() (v int64, err error) {
    defer yellowstone.Monitor("VM.pop", []any{vm}, &v, &err)()
    // ...
}
```

Monitoring is off by default and enabled via the env var `YELLOWSTONE_MONITOR` (truthy: `1`, `true`, `yes`, `on`). Output goes to stderr; program stdout is unaffected.

### Annotation markers

Two marker syntaxes, both using the exact byte sequence `// 🦬`:

| Marker | Placement | Effect of `erase` |
|---|---|---|
| `// 🦬` | end of line | delete the entire line |
| `// 🦬: <text>` | own line | delete marker line AND replace the next line with `<text>` |

Typical changes when annotating a function:

1. Add `import "yellowstone" // 🦬` to the file. Keep it as a separate `import` statement rather than merging into a block — one-liner imports are independently deletable via the trailing marker.
2. Add `defer yellowstone.Monitor(...)() // 🦬` as the first statement of the function.

**Named returns are permanent, not annotation-only.** When a function's outputs need to be captured by `Monitor`, we name its return values (`(v int64, err error)` rather than `(int64, error)`) and keep that signature in the default code. Any body adjustments this forces — e.g. `v := expr` → `v = expr` because the named return pre-declares `v` — are also kept in the default. This way adding or removing monitoring only touches the marker lines, not the function signature or body.

### The `erase` tool and CI guarantee

The erase tool (`go/tools/erase/`) walks a source directory and writes stripped copies to an output directory. `.go` and `go.mod` files are stripped; all other files are copied verbatim. Build it once and invoke as `erase -src go/example -out <dir>`.

CI (`.github/workflows/go-example.yml`) enforces:

1. Build with annotations in place.
2. `go run .` prints `20` (monitor off).
3. `YELLOWSTONE_MONITOR=1 go run .` prints `20` on stdout (traces on stderr).
4. Erase, then `go run .` in the erased tree, still prints `20`.

This is a **semantic** check (same observable behavior), not a byte-for-byte check against a baseline. A byte-for-byte check would require storing a separate "pre-annotation" snapshot of the source — reintroducing the duplication we just removed.

### Currently annotated

- `vm.go`: `push`, `pop`, `Run`
