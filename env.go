// Package env reads configuration from the process environment, falling
// back to a .env file when the variable isn't set. This makes one call
// work both in production (real env vars) and local dev under tools like
// the tinywasm daemon, which runs the app binary as a child process
// without exporting .env into its environment.
package env

import (
	"os"

	. "github.com/tinywasm/fmt"
)

// Get returns the value of key from the process environment, or — if
// unset — the value of key=... in path (typically ".env" in the working
// directory). Values wrapped in double quotes have them stripped. Returns
// "" if the key is set nowhere.
func Get(key string) string {
	return GetFrom(key, ".env")
}

// GetFrom is Get with an explicit .env path, for callers not running from
// the project root.
func GetFrom(key, path string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	prefix := key + "="
	for _, line := range Split(string(data), "\n") {
		if !HasPrefix(line, prefix) {
			continue
		}
		v := Convert(line).TrimPrefix(prefix).String()
		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			v = v[1 : len(v)-1]
		}
		return v
	}
	return ""
}

// GetRequired returns key from env, .env, or fallback paths in order.
// Returns an error if key is not found anywhere.
func GetRequired(key string, fallbackPaths ...string) (string, error) {
	if v := Get(key); v != "" {
		return v, nil
	}
	for _, path := range fallbackPaths {
		if v := GetFrom(key, path); v != "" {
			return v, nil
		}
	}
	return "", Err(key + " is not set")
}

// Arg returns the value for -key=value or -key value in os.Args (skipping
// argv[0]). Unset flags return "" — callers apply their own default, same
// as Get. Common pattern: tools like the tinywasm dev daemon pass runtime
// config as CLI flags to the child process instead of env vars.
func Arg(key string) string {
	prefix := "-" + key + "="
	args := os.Args[1:]
	for i, arg := range args {
		if HasPrefix(arg, prefix) {
			return Convert(arg).TrimPrefix(prefix).String()
		}
		if arg == "-"+key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
