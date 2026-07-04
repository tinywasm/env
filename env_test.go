package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFrom_EnvVarWins(t *testing.T) {
	t.Setenv("ENV_TEST_KEY", "from-process")
	if got := GetFrom("ENV_TEST_KEY", "does-not-exist.env"); got != "from-process" {
		t.Errorf("got %q, want %q", got, "from-process")
	}
}

func TestGetFrom_FallsBackToFile(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("OTHER=1\nENV_TEST_KEY=\"from-file\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if got := GetFrom("ENV_TEST_KEY", path); got != "from-file" {
		t.Errorf("got %q, want %q", got, "from-file")
	}
}

func TestGetFrom_MissingEverywhere(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY_MISSING")
	if got := GetFrom("ENV_TEST_KEY_MISSING", "does-not-exist.env"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestArg_EqualsForm(t *testing.T) {
	restore := setArgs("cmd", "-server_port=6060")
	defer restore()
	if got := Arg("server_port"); got != "6060" {
		t.Errorf("got %q, want %q", got, "6060")
	}
}

func TestArg_SpaceForm(t *testing.T) {
	restore := setArgs("cmd", "-server_port", "6060")
	defer restore()
	if got := Arg("server_port"); got != "6060" {
		t.Errorf("got %q, want %q", got, "6060")
	}
}

func TestArg_Missing(t *testing.T) {
	restore := setArgs("cmd")
	defer restore()
	if got := Arg("server_port"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestGetRequired_FromEnv(t *testing.T) {
	t.Setenv("ENV_TEST_KEY", "from-process")
	got, err := GetRequired("ENV_TEST_KEY", "does-not-exist.env")
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != "from-process" {
		t.Errorf("got %q, want %q", got, "from-process")
	}
}

func TestGetRequired_FromFirstFallback(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY")
	path1 := filepath.Join(t.TempDir(), ".env.1")
	if err := os.WriteFile(path1, []byte("ENV_TEST_KEY=\"from-fallback-1\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := GetRequired("ENV_TEST_KEY", path1)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != "from-fallback-1" {
		t.Errorf("got %q, want %q", got, "from-fallback-1")
	}
}

func TestGetRequired_FromSecondFallback(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY")
	path1 := filepath.Join(t.TempDir(), ".env.1")
	path2 := filepath.Join(t.TempDir(), ".env.2")
	if err := os.WriteFile(path2, []byte("ENV_TEST_KEY=\"from-fallback-2\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := GetRequired("ENV_TEST_KEY", path1, path2)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}
	if got != "from-fallback-2" {
		t.Errorf("got %q, want %q", got, "from-fallback-2")
	}
}

func TestGetRequired_NotFound(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY")
	got, err := GetRequired("ENV_TEST_KEY", "does-not-exist-1.env", "does-not-exist-2.env")
	if err == nil {
		t.Errorf("got nil error, want error for missing key")
	}
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func setArgs(args ...string) func() {
	orig := os.Args
	os.Args = args
	return func() { os.Args = orig }
}
