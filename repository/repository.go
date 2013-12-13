// Package repository supports a debian repository according to https://wiki.debian.org/RepositoryFormat
package repository

import "fmt"

// Repository models a debian repository Repository
type Repository struct {
	Sections,
	Archs,
	Dists []string
	Translations []string
	DebMirrorTag string // support local directory layout for debmirror
}

// New inits repository for usage
func New() (*Repository, error) {
	return &Repository{
		Sections:     []string{"main", "contrib", "non-free", "main/debian-installer"},
		Archs:        []string{"none"},
		Dists:        []string{"stable"},
		Translations: []string{"en"},
	}, nil
}

// Extensions
var ListExtensions = []string{
	".bz2",
	".gz",
	"",
}

// PackageListNames lists all possible locations for package lists
func (r *Repository) PackageListNames() (names []string) {
	debmirrorTag := r.DebMirrorTag
	if debmirrorTag != "" {
		debmirrorTag = "/" + debmirrorTag
	}
	for _, section := range r.Sections {
		for _, arch := range r.Archs {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					///security.debian.org/dists/squeeze/updates/non-free/binary-amd64/Packages.bz2
					name := fmt.Sprintf("/dists/%s%s/%s/binary-%s/Packages%s", dist, debmirrorTag, section, arch, ext)
					names = append(names, name)
				}

			}
		}
	}
	return
}

// TranslationListNames lists all possible locations for translation lists
func (r *Repository) TranslationListNames() (names []string) {
	debmirrorTag := r.DebMirrorTag
	if debmirrorTag != "" {
		debmirrorTag = "/" + debmirrorTag
	}
	for _, section := range r.Sections {
		for _, translation := range r.Translations {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					//ftp.uk.debian.org/debian/dists/wheezy/contrib/i18n/Translation-en.bz2
					name := fmt.Sprintf("/dists/%s%s/%s/i18n/Translation-%s%s", dist, debmirrorTag, section, translation, ext)
					names = append(names, name)
				}

			}
		}
	}
	return
}
