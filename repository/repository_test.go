package repository

import (
	"reflect"
	"sort"
	"testing"
)

func TestRepositoryPackageListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	testListing(t, "", r.PackageListNames(), defaultRepositoryPackageListNamesUntagged)
	r.DebMirrorTag = "0"
	testListing(t, "0", r.PackageListNames(), defaultRepositoryPackageListNames)
}

func TestRepositoryTranslationListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	testListing(t, "", r.TranslationListNames(), defaultRepositoryTranslationListNamesUntagged)
	r.DebMirrorTag = "0"
	testListing(t, "0", r.TranslationListNames(), defaultRepositoryTranslationListNames)
}

func TestRepositoryKeySets(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}
	r.Archs = []string{"amd64", "i386"}
	r.DebMirrorTag = "0"

	want := x86KeySets

	got := r.KeySets()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got\n%+v\nwant\n%+v\n", got, want)
	}
}

func TestRepositoryUniqeKeys(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}
	r.Archs = []string{"amd64", "i386"}
	r.DebMirrorTag = "0"

	want := []string{
		"stable:amd64",
		"stable:i386",
	}

	sort.Strings(want)

	got := r.UniqueKeys()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got\n%+v\nwant\n%+v\n", got, want)
	}
}

func testListing(t *testing.T, debmirrorTag string, gotNames, wantNames []string) {
	if len(gotNames) == 0 {
		t.Errorf("List of package names for tag %q is empty", debmirrorTag)
	}
	if len(gotNames) != len(wantNames) {
		t.Errorf("length mismatch, got %d, expected %d", len(gotNames), len(wantNames))
	} else {
		t.Logf("%d package lists will by tried", len(gotNames))
	}
	for i, got := range gotNames {
		want := wantNames[i]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		} else {
			t.Log(got)
		}
	}
}

var defaultRepositoryPackageListNamesUntagged = []string{
	"/dists/stable/main/binary-none/Packages.bz2",
	"/dists/stable/main/binary-none/Packages.gz",
	"/dists/stable/main/binary-none/Packages",
	"/dists/stable/contrib/binary-none/Packages.bz2",
	"/dists/stable/contrib/binary-none/Packages.gz",
	"/dists/stable/contrib/binary-none/Packages",
	"/dists/stable/non-free/binary-none/Packages.bz2",
	"/dists/stable/non-free/binary-none/Packages.gz",
	"/dists/stable/non-free/binary-none/Packages",
	"/dists/stable/main/debian-installer/binary-none/Packages.bz2",
	"/dists/stable/main/debian-installer/binary-none/Packages.gz",
	"/dists/stable/main/debian-installer/binary-none/Packages",
}

var defaultRepositoryPackageListNames = []string{
	"/dists/stable/0/main/binary-none/Packages.bz2",
	"/dists/stable/0/main/binary-none/Packages.gz",
	"/dists/stable/0/main/binary-none/Packages",
	"/dists/stable/0/contrib/binary-none/Packages.bz2",
	"/dists/stable/0/contrib/binary-none/Packages.gz",
	"/dists/stable/0/contrib/binary-none/Packages",
	"/dists/stable/0/non-free/binary-none/Packages.bz2",
	"/dists/stable/0/non-free/binary-none/Packages.gz",
	"/dists/stable/0/non-free/binary-none/Packages",
	"/dists/stable/0/main/debian-installer/binary-none/Packages.bz2",
	"/dists/stable/0/main/debian-installer/binary-none/Packages.gz",
	"/dists/stable/0/main/debian-installer/binary-none/Packages",
}

var x86KeySets = map[string][]string{
	"stable:amd64": {
		"/dists/stable/0/main/binary-amd64/Packages.bz2",
		"/dists/stable/0/main/binary-amd64/Packages.gz",
		"/dists/stable/0/main/binary-amd64/Packages",
		"/dists/stable/0/contrib/binary-amd64/Packages.bz2",
		"/dists/stable/0/contrib/binary-amd64/Packages.gz",
		"/dists/stable/0/contrib/binary-amd64/Packages",
		"/dists/stable/0/non-free/binary-amd64/Packages.bz2",
		"/dists/stable/0/non-free/binary-amd64/Packages.gz",
		"/dists/stable/0/non-free/binary-amd64/Packages",
		"/dists/stable/0/main/debian-installer/binary-amd64/Packages.bz2",
		"/dists/stable/0/main/debian-installer/binary-amd64/Packages.gz",
		"/dists/stable/0/main/debian-installer/binary-amd64/Packages",
	},
	"stable:i386": {
		"/dists/stable/0/main/binary-i386/Packages.bz2",
		"/dists/stable/0/main/binary-i386/Packages.gz",
		"/dists/stable/0/main/binary-i386/Packages",
		"/dists/stable/0/contrib/binary-i386/Packages.bz2",
		"/dists/stable/0/contrib/binary-i386/Packages.gz",
		"/dists/stable/0/contrib/binary-i386/Packages",
		"/dists/stable/0/non-free/binary-i386/Packages.bz2",
		"/dists/stable/0/non-free/binary-i386/Packages.gz",
		"/dists/stable/0/non-free/binary-i386/Packages",
		"/dists/stable/0/main/debian-installer/binary-i386/Packages.bz2",
		"/dists/stable/0/main/debian-installer/binary-i386/Packages.gz",
		"/dists/stable/0/main/debian-installer/binary-i386/Packages",
	},
}

var defaultRepositoryTranslationListNamesUntagged = []string{
	"/dists/stable/main/i18n/Translation-en.bz2",
	"/dists/stable/main/i18n/Translation-en.gz",
	"/dists/stable/main/i18n/Translation-en",
	"/dists/stable/contrib/i18n/Translation-en.bz2",
	"/dists/stable/contrib/i18n/Translation-en.gz",
	"/dists/stable/contrib/i18n/Translation-en",
	"/dists/stable/non-free/i18n/Translation-en.bz2",
	"/dists/stable/non-free/i18n/Translation-en.gz",
	"/dists/stable/non-free/i18n/Translation-en",
	"/dists/stable/main/debian-installer/i18n/Translation-en.bz2",
	"/dists/stable/main/debian-installer/i18n/Translation-en.gz",
	"/dists/stable/main/debian-installer/i18n/Translation-en",
}

var defaultRepositoryTranslationListNames = []string{
	"/dists/stable/0/main/i18n/Translation-en.bz2",
	"/dists/stable/0/main/i18n/Translation-en.gz",
	"/dists/stable/0/main/i18n/Translation-en",
	"/dists/stable/0/contrib/i18n/Translation-en.bz2",
	"/dists/stable/0/contrib/i18n/Translation-en.gz",
	"/dists/stable/0/contrib/i18n/Translation-en",
	"/dists/stable/0/non-free/i18n/Translation-en.bz2",
	"/dists/stable/0/non-free/i18n/Translation-en.gz",
	"/dists/stable/0/non-free/i18n/Translation-en",
	"/dists/stable/0/main/debian-installer/i18n/Translation-en.bz2",
	"/dists/stable/0/main/debian-installer/i18n/Translation-en.gz",
	"/dists/stable/0/main/debian-installer/i18n/Translation-en",
}
