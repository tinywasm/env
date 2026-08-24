package env

import "testing"

func staticReader(values map[string]string) Reader {
	return func(key string) string { return values[key] }
}

func TestGet_ReadsFromReader(t *testing.T) {
	read := staticReader(map[string]string{"KEY": "value"})
	if got := Get(read, "KEY"); got != "value" {
		t.Errorf("got %q, want %q", got, "value")
	}
}

func TestGet_MissingReturnsEmpty(t *testing.T) {
	read := staticReader(nil)
	if got := Get(read, "MISSING"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestGetRequired_PrimaryWins(t *testing.T) {
	read := staticReader(map[string]string{"KEY": "primary"})
	fallback := staticReader(map[string]string{"KEY": "fallback"})
	got, err := GetRequired(read, "KEY", fallback)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != "primary" {
		t.Errorf("got %q, want %q", got, "primary")
	}
}

func TestGetRequired_FallsBackInOrder(t *testing.T) {
	read := staticReader(nil)
	fallback1 := staticReader(nil)
	fallback2 := staticReader(map[string]string{"KEY": "from-fallback-2"})
	got, err := GetRequired(read, "KEY", fallback1, fallback2)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got != "from-fallback-2" {
		t.Errorf("got %q, want %q", got, "from-fallback-2")
	}
}

func TestGetRequired_NotFoundAnywhere(t *testing.T) {
	read := staticReader(nil)
	got, err := GetRequired(read, "KEY", staticReader(nil))
	if err == nil {
		t.Error("got nil error, want error for missing key")
	}
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestGetRequired_NoFallbacksNotFound(t *testing.T) {
	read := staticReader(nil)
	_, err := GetRequired(read, "KEY")
	if err == nil {
		t.Error("got nil error, want error for missing key with no fallbacks")
	}
}
