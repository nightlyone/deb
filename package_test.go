package deb

import (
	"archive/tar"
	"io/ioutil"
	"strings"
	"testing"
)

func TestExtractPackageName(t *testing.T) {
	want := "golang-mode"
	got := NewPackage("fixtures/golang-mode_1.1.2-2_all.deb").Name
	if want != got {
		t.Errorf("want = '%v' got = '%v'", want, got)
	}
}

func TestExtractPackageArch(t *testing.T) {
	want := "all"
	got := NewPackage("fixtures/golang-mode_1.1.2-2_all.deb").Arch
	if want != got {
		t.Errorf("want = '%v' got = '%v'", want, got)
	}
}

func TestExtractPackageVersion(t *testing.T) {
	want := "1.1.2-2"
	got := NewPackage("fixtures/golang-mode_1.1.2-2_all.deb").Version
	if want != got {
		t.Errorf("want = '%v' got = '%v'", want, got)
	}
}

func TestGetArchive(t *testing.T) {
	filename := "fixtures/golang-mode_1.1.2-2_all.deb"

	// this variable is just used for type checking
	var archive *tar.Reader
	archive, err := NewPackage(filename).archiveOpen()
	if err != nil {
		t.Fatal("Cannot load fixture", filename, "reason", err)
		return
	}

	_ = archive
}

func TestExtractPackageSource(t *testing.T) {
	want := "golang"
	got, err := NewPackage("fixtures/golang-mode_1.1.2-2_all.deb").Source()
	if err != nil {
		t.Errorf("got error '%v'", err)
	}
	if want != got {
		t.Errorf("want = '%v' got = '%v'", want, got)
	}
}

func TestFindChangelog(t *testing.T) {
	filename := "fixtures/golang-mode_1.1.2-2_all.deb"
	pack := NewPackage(filename)

	clname := "./usr/share/doc/" + pack.Name + "/changelog.Debian.gz"
	r, err := pack.FindFile(clname)
	if err != nil {
		t.Error("Cannot load package changelog, reason", err)
		return
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error("Cannot load package changelog content, reason", err)
		return
	}
	want := "golang (2:1.1.2-2) unstable; urgency=low"
	got := strings.Split(string(b), "\n")[0]
	if got != want {
		t.Errorf("unexpected changelog content: want '%s', got '%s'", want, got)
	}
}
