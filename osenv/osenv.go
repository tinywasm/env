// Package osenv provides the process-backed env.Reader/env.Writer: real
// environment variables, falling back to a .env file, plus CLI argument
// lookup. This is the concrete implementation a server binary injects —
// it imports "os" on purpose, so it never compiles into a wasm/edge
// binary that only ever injects a platform-specific Reader instead.
package osenv

import (
	"os"

	"github.com/tinywasm/env"
	"github.com/tinywasm/fmt"
)

const defaultDotEnvPath = ".env"

// Reader returns an env.Reader backed by the process environment, falling
// back to a .env file in the working directory when a key is unset.
func Reader() env.Reader {
	return ReaderFrom(defaultDotEnvPath)
}

// ReaderFrom is Reader with an explicit .env path, for callers not
// running from the project root.
func ReaderFrom(path string) env.Reader {
	return func(key string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return ""
		}
		prefix := key + "="
		for _, line := range fmt.Split(string(data), "\n") {
			if !fmt.HasPrefix(line, prefix) {
				continue
			}
			v := fmt.Convert(line).TrimPrefix(prefix).String()
			if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			return v
		}
		return ""
	}
}

// Writer returns an env.Writer that sets key in the current process's own
// environment (os.Setenv) — it does not persist to a .env file, and it
// only affects this process, not its parent or children already started.
func Writer() env.Writer {
	return func(key, value string) error {
		return os.Setenv(key, value)
	}
}

// Arg returns the value for -key=value or -key value in os.Args (skipping
// argv[0]). Unset flags return "" — callers apply their own default.
// There is no edge/wasm equivalent: a Cloudflare Worker has no argv, so
// this stays osenv-only rather than living in the agnostic env package.
func Arg(key string) string {
	prefix := "-" + key + "="
	args := os.Args[1:]
	for i, arg := range args {
		if fmt.HasPrefix(arg, prefix) {
			return fmt.Convert(arg).TrimPrefix(prefix).String()
		}
		if arg == "-"+key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
