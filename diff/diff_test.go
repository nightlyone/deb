package diff

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func open_fixture(t *testing.T, name string) (io.Reader, func()) {
	r, err := os.Open(filepath.Join("./fixtures", name))
	if err != nil {
		t.Fatalf("Cannot read fixture %q", name)
		return nil, nil
	}
	return r, func() { r.Close() }
}

func TestParse(t *testing.T) {
	fixture := "old.Packages"
	r, cleanup := open_fixture(t, fixture)
	defer cleanup()

	list, err := NewList(r)
	if err != nil {
		t.Fatalf("Cannot parse %q", fixture)
	}
	for pkg, src := range list.Package {
		t.Logf("%s=%s (%s)", pkg, list.Version[pkg], src)
	}
}

func TestDiff(t *testing.T) {
	fixture_old := "old.Packages"
	r_old, cleanup_old := open_fixture(t, fixture_old)
	defer cleanup_old()
	fixture_new := "new.Packages"
	r_new, cleanup_new := open_fixture(t, fixture_new)
	defer cleanup_new()

	older, err := NewList(r_old)
	if err != nil {
		t.Fatal(err)
	}
	newer, err := NewList(r_new)
	if err != nil {
		t.Fatal(err)
	}
	a, r, u := Changes(newer, older)
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
