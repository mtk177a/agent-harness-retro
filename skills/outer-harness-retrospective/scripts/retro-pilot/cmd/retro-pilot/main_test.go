package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrivateOutputStaysOutsideRepository(t *testing.T) {
	root := t.TempDir()
	working := filepath.Join(root, "skills", "outer-harness-retrospective", "scripts", "retro-pilot")
	if err := os.MkdirAll(working, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	inside := filepath.Join(root, "tmp")
	if err := os.Mkdir(inside, 0o700); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	for _, test := range []struct {
		name string
		path string
		want bool
	}{
		{"inside repository", inside, false},
		{"outside repository", outside, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := outsideProtectedRoot(test.path, working)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("outsideProtectedRoot() = %v, want %v", got, test.want)
			}
		})
	}
}
