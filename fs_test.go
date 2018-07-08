package bit

import (
	"os"
	"testing"
)

func TestMockSaveAndLoad(t *testing.T) {
	fs := NewMockFileSystem()
	want := "baz"
	err := fs.SaveText("/foo/bar", want, 0644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have, err := fs.LoadText("/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMockOverwriteAndLoad(t *testing.T) {
	fs := NewMockFileSystem()
	want := "baz"
	err := fs.SaveText("/foo/bar", "zzz", 0644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = fs.SaveText("/foo/bar", want, 0644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have, err := fs.LoadText("/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestNoFile(t *testing.T) {
	fs := NewMockFileSystem()
	want := os.ErrNotExist
	_, have := fs.LoadText("/foo/bar")
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMkdir(t *testing.T) {
	fs := NewMockFileSystem()
	fs.Mkdir("/foo/bar", 0755)
	want := true
	mode, err := fs.Mode("/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have := mode.IsDir()
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestMode(t *testing.T) {
	fs := NewMockFileSystem()
	want := os.FileMode(0x123)
	err := fs.SaveText("/foo/bar", "", want)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have, err := fs.Mode("/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestNoMode(t *testing.T) {
	fs := NewMockFileSystem()
	want := os.ErrNotExist
	_, have := fs.Mode("/foo/bar")
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestCwd(t *testing.T) {
	fs := NewMockFileSystem()
	want := "/"
	have, err := fs.Getwd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}

func TestChdir(t *testing.T) {
	fs := NewMockFileSystem()
	want := "/foo/bar"
	if err := fs.Chdir(want); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	have, err := fs.Getwd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want != have {
		t.Errorf("\n want: %v \n have: %v", want, have)
	}
}
