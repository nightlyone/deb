// Package repository supports a debian repository according to https://wiki.debian.org/RepositoryFormat
package repository

import "fmt"

// Repository models a debian repository Repository
type Repository struct {
	Sections,
	Archs,
	Dists []string
	Translations []string
}

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

func (r *Repository) PackageListNames(tag string) (names []string) {
	for _, section := range r.Sections {
		for _, arch := range r.Archs {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					///security.debian.org/dists/squeeze/updates/non-free/binary-amd64/Packages.bz2
					name := fmt.Sprintf("/dists/%s/%s/binary-%s/Packages%s", dist, section, arch, ext)
					names = append(names, name)
				}

			}
		}
	}
	return
}

func (r *Repository) TranslationListNames(tag string) (names []string) {
	for _, section := range r.Sections {
		for _, translation := range r.Translations {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					///security.debian.org/dists/squeeze/updates/non-free/binary-amd64/Packages.bz2
					//ftp.uk.debian.org/debian/dists/wheezy/contrib/i18n/Translation-en.bz2
					name := fmt.Sprintf("/dists/%s/%s/i18n/Translation-%s%s", dist, section, translation, ext)
					names = append(names, name)
				}

			}
		}
	}
	return
}
