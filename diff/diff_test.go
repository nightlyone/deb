package diff

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func openFixture(t *testing.T, name string) (io.Reader, func()) {
	r, err := os.Open(filepath.Join("./fixtures", name))
	if err != nil {
		t.Fatalf("Cannot read fixture %q", name)
		return nil, nil
	}
	return r, func() { r.Close() }
}

var listpackages = flag.Bool("listpackages", false, "use listpackages to debug changelog parsing")

func TestParse(t *testing.T) {
	fixture := "old.Packages"
	r, cleanup := openFixture(t, fixture)
	defer cleanup()

	list, err := NewList(r)
	if err != nil {
		t.Fatalf("Cannot parse %q", fixture)
	}

	count := 1656
	if len(list.Package) != count {
		t.Errorf("cannot parse packages. got %d, want %d")

	}
	if len(list.Version) != count {
		t.Errorf("cannot parse versions. got %d, want %d")

	}
	if len(list.Location) != count {
		t.Errorf("cannot parse locations. got %d, want %d")

	}
	if *listpackages {
		for pkg, src := range list.Package {
			t.Logf("%s=%s (%s) in %s", pkg, list.Version[pkg], src, list.Location[pkg])
		}
	}
}

func TestDiff(t *testing.T) {
	lists := map[string]*List{}
	for _, s := range []string{"old", "new"} {
		r, cleanup := openFixture(t, s+".Packages")
		defer cleanup()
		list, err := NewList(r)
		if err != nil {
			t.Fatal(err)
		}
		lists[s] = list
	}

	a, r, u := Changes(lists["new"], lists["old"])
	var changes = []struct {
		what string
		got  []*Package
		want int
	}{
		{
			what: "added",
			got:  a,
			want: 1,
		},
		{
			what: "removed",
			got:  r,
			want: 0,
		},
		{
			what: "updated",
			got:  u,
			want: 29,
		},
	}
	summary := []interface{}{"summary:"}
	for _, c := range changes {
		if len(c.got) != c.want {
			t.Errorf("expect %d packages %s, got %d", c.want, c.what, len(c.got))
		} else {
			summary = append(summary, fmt.Sprintf("%s %d", c.what, c.want))
		}
	}
	t.Log(summary...)
	cu := CompressUpdates(u)
	if len(cu) != 5 {
		for _, ud := range cu {
			t.Logf("update %s to %s", ud.Package, ud.Version)
			if len(ud.Related) > 1 {
				summary := []interface{}{"related", len(ud.Related), "are"}
				for _, rel := range ud.Related {
					summary = append(summary, rel.Name)
				}
				t.Log(summary...)
			}
		}
		t.Error("compressed from", len(u), "to", len(cu), "expected", 5)
	}
}
