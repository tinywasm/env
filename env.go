// Package env reads and writes configuration through an injected source,
// agnostic of where that source lives. This keeps the package wasm-safe
// (no "os" import): a server binary injects a Reader/Writer backed by the
// process environment (see the osenv subpackage), while an edge binary
// injects one backed by its own platform binding (e.g. Cloudflare env
// vars) — the same Get/GetRequired call works against either.
package env

import "github.com/tinywasm/fmt"

// Reader returns the value of key from whatever source it wraps, or "" if
// key is unset there. "" always means "not found" — no separate ok bool,
// matching the rest of the ecosystem's convention.
type Reader func(key string) string

// Writer sets key to value in whatever source it wraps. Not every source
// can persist a write (e.g. platform env vars configured outside the
// running process); such a Writer returns a descriptive error rather than
// silently discarding the value.
type Writer func(key, value string) error

// Get returns the value of key from read, or "" if unset there.
func Get(read Reader, key string) string {
	return read(key)
}

// GetRequired tries read first, then each fallback Reader in order.
// Returns an error if key is not found anywhere.
func GetRequired(read Reader, key string, fallbacks ...Reader) (string, error) {
	if v := read(key); v != "" {
		return v, nil
	}
	for _, fb := range fallbacks {
		if v := fb(key); v != "" {
			return v, nil
		}
	}
	return "", fmt.Err(key + " is not set")
}
