package changelog

import (
	"os"
	"testing"
)

func TestParseChangelog(t *testing.T) {
	r, err := os.Open("fixtures/changelog.Debian")
	if err != nil {
		t.Fatal("Cannot load fixture, reason", err)
		return
	}
	defer r.Close()

	cl, err := Parse(r, "")
	if err != nil {
		t.Error("unexpected error", err)
		return
	}
	if len(cl) == 0 {
		t.Error("empty changelog")
	}
}

func TestParseChangelogUntil(t *testing.T) {
	r, err := os.Open("fixtures/changelog.Debian")
	if err != nil {
		t.Fatal("Cannot load fixture, reason", err)
		return
	}
	defer r.Close()

	cl, err := Parse(r, "2:1.1.2-2")
	if err != nil {
		t.Fatal("unexpected error", err)
		return
	}
	if len(cl) == 0 {
		t.Fatal("empty changelog")
	}
	want := 1
	got := len(cl)

	if got != want {
		t.Errorf("got %d changelog entries, want %d", got, want)
	}
}
