package repository

import "testing"

func TestRepositoryPackageListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	testListing(t, "0", r.PackageListNames("0"), defaultRepositoryPackageListNames)
	testListing(t, "", r.PackageListNames(""), defaultRepositoryPackageListNamesUntagged)
}

func TestRepositoryTranslationListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	testListing(t, "0", r.TranslationListNames("0"), defaultRepositoryTranslationListNames)
	testListing(t, "", r.TranslationListNames(""), defaultRepositoryTranslationListNamesUntagged)
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
