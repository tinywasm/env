# env
<img src="docs/img/badges.svg">

Auto-tagged env access: `!wasm` reads `os.Getenv` (+ `.env` fallback), `wasm` reads Cloudflare `context.env` via `syscall/js`. No injection — the build tag selects the implementation. `tinywasm/fmt` is the only dep, `os` never leaks into wasm.

```go
import "github.com/tinywasm/env"

port := env.Get("PORT")                 // "" if unset
dsn, err := env.Require("DATABASE_URL") // error if unset
host := env.GetOr("HOST", "localhost")
if v, ok := env.Lookup("KEY"); ok { ... }
env.Set("KEY", "value") // !wasm only
flag := env.Arg("port") // !wasm only, -port=8080 / -port 8080
```

Low-level `.env` path override: `env.LookupAt("KEY", "/path/to/.env")`.

Wasm target is Cloudflare Workers (`context.env`); other wasm targets get `""`.
