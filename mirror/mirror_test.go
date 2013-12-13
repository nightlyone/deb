package mirror

import "testing"

var defaultMirrorPackageListNames = []string{
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

func TestMirrorPackageListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	names := r.PackageListNames("stable")
	if len(names) == 0 {
		t.Error("List of package names is empty")
	}
	if len(names) != len(defaultMirrorPackageListNames) {
		t.Errorf("length mismatch, got %d, expected %d", len(names), len(defaultMirrorPackageListNames))
	} else {
		t.Logf("%d package lists will by tried", len(names))
	}
	for i, got := range names {
		want := defaultMirrorPackageListNames[i]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		} else {
			t.Log(got)
		}
	}
}

var defaultMirrorTranslationListNames = []string{
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

func TestMirrorTranslationListNames(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
		return
	}

	names := r.TranslationListNames("stable")
	if len(names) == 0 {
		t.Error("List of package names is empty")
	}
	if len(names) != len(defaultMirrorTranslationListNames) {
		t.Errorf("length mismatch, got %d, expected %d", len(names), len(defaultMirrorTranslationListNames))
	} else {
		t.Logf("%d package lists will by tried", len(names))
	}
	for i, got := range names {
		want := defaultMirrorTranslationListNames[i]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		} else {
			t.Log(got)
		}
	}
}
