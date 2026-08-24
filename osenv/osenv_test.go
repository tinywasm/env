package osenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReaderFrom_EnvVarWins(t *testing.T) {
	t.Setenv("ENV_TEST_KEY", "from-process")
	read := ReaderFrom("does-not-exist.env")
	if got := read("ENV_TEST_KEY"); got != "from-process" {
		t.Errorf("got %q, want %q", got, "from-process")
	}
}

func TestReaderFrom_FallsBackToFile(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("OTHER=1\nENV_TEST_KEY=\"from-file\"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	read := ReaderFrom(path)
	if got := read("ENV_TEST_KEY"); got != "from-file" {
		t.Errorf("got %q, want %q", got, "from-file")
	}
}

func TestReaderFrom_MissingEverywhere(t *testing.T) {
	os.Unsetenv("ENV_TEST_KEY_MISSING")
	read := ReaderFrom("does-not-exist.env")
	if got := read("ENV_TEST_KEY_MISSING"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestWriter_SetsProcessEnv(t *testing.T) {
	os.Unsetenv("ENV_TEST_WRITE_KEY")
	write := Writer()
	if err := write("ENV_TEST_WRITE_KEY", "written"); err != nil {
		t.Fatalf("write: %v", err)
	}
	if got := os.Getenv("ENV_TEST_WRITE_KEY"); got != "written" {
		t.Errorf("got %q, want %q", got, "written")
	}
}

func TestArg_EqualsForm(t *testing.T) {
	restore := setArgs("cmd", "-server_port=8080")
	defer restore()
	if got := Arg("server_port"); got != "8080" {
		t.Errorf("got %q, want %q", got, "8080")
	}
}

func TestArg_SpaceForm(t *testing.T) {
	restore := setArgs("cmd", "-server_port", "8080")
	defer restore()
	if got := Arg("server_port"); got != "8080" {
		t.Errorf("got %q, want %q", got, "8080")
	}
}

func TestArg_Missing(t *testing.T) {
	restore := setArgs("cmd")
	defer restore()
	if got := Arg("server_port"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func setArgs(args ...string) func() {
	orig := os.Args
	os.Args = args
	return func() { os.Args = orig }
}
