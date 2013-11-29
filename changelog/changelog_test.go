package changelog

import (
	"os"
	"reflect"
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

func TestParseVersionLine(t *testing.T) {
	for i, test := range versionLines {
		change := Change{}
		ok := change.parseVersionLine([]byte(test.line))
		if ok && test.bad {
			t.Errorf("%d: did parse %q, but shouldn't", i, test.line)
		} else if !ok && !test.bad {
			t.Errorf("%d: cannot parse %q", i, test.line)
		}
		if ok && !reflect.DeepEqual(test.want, change) {
			t.Errorf("%d: want:\n%v\n, got:\n%v\n", i, test.want, change)
		}
	}
}

var versionLines = []struct {
	line string
	bad  bool
	want Change
}{
	{

		line: "golang (2:1.1.2-2ubuntu1) saucy; urgency=low",
		want: Change{
			Source:  "golang",
			Version: "2:1.1.2-2ubuntu1",
			Dist:    "saucy",
			Urgency: "low",
		},
	},
	{

		line: "golang (2:1.1.2-2ubuntu1) saucy; urgency=low",
		want: Change{
			Source:  "golang",
			Version: "2:1.1.2-2ubuntu1",
			Dist:    "saucy",
			Urgency: "low",
		},
	},
	{

		line: "golang (2:1.1.2-2ubuntu1) ; urgency=low",
		bad:  true,
	},
}
