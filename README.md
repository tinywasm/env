# env
<img src="docs/img/badges.svg">

Agnostic config reads/writes: `Get`/`GetRequired` take an injected `Reader`,
so this package never imports `os` and stays wasm-safe. A server binary
injects `osenv.Reader()` (process env vars, falling back to a `.env` file);
an edge binary injects a `Reader` over its own platform binding — the same
call works against either.

```go
import (
    "github.com/tinywasm/env"
    "github.com/tinywasm/env/osenv"
)

dsn, err := env.GetRequired(osenv.Reader(), "DATABASE_URL")
```

`osenv` also carries `Arg(key)` (CLI flag lookup via `os.Args`) and
`Writer()` (sets a value in the current process's own environment) — both
process-only, with no edge equivalent.
