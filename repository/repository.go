// Package repository supports a debian repository according to https://wiki.debian.org/RepositoryFormat
package repository

import (
	"fmt"
	"sort"
)

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

// PackageWalker can be iterates over all possible positions of packages
type PackageVisitor func(dist, tag, section, arch, ext string)

// ForEachPackageSpec iterates over all possible positions of packages
func (r *Repository) ForEachPackageSpec(v PackageVisitor) {
	debmirrorTag := r.DebMirrorTag
	if debmirrorTag != "" {
		debmirrorTag = "/" + debmirrorTag
	}
	for _, section := range r.Sections {
		for _, arch := range r.Archs {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					v(dist, debmirrorTag, section, arch, ext)
				}

			}
		}
	}
}

// UniqueKeys returns the unqiue package list sets of a repository
func (r *Repository) UniqueKeys() []string {
	key := map[string]bool{}
	r.ForEachPackageSpec(func(dist, tag, section, arch, ext string) {
		key[dist+":"+arch] = true
	})

	keys := []string{}
	for k := range key {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys

}

// KeySets returns how package lists can be merged
// They are indexed by the values returned in UniqueKeys
func (r *Repository) KeySets() map[string][]string {
	set := map[string][]string{}
	r.ForEachPackageSpec(func(dist, tag, section, arch, ext string) {
		key := dist + ":" + arch
		name := fmt.Sprintf("/dists/%s%s/%s/binary-%s/Packages%s", dist, tag, section, arch, ext)
		set[key] = append(set[key], name)
	})

	return set

}

// TranslationWalker can be iterates over all possible positions of translations
type TranslationVisitor func(dist, tag, section, translation, ext string)

// ForEachTranslationSpec iterates over all possible positions of translations
func (r *Repository) ForEachTranslationSpec(v TranslationVisitor) {
	debmirrorTag := r.DebMirrorTag
	if debmirrorTag != "" {
		debmirrorTag = "/" + debmirrorTag
	}
	for _, section := range r.Sections {
		for _, translation := range r.Translations {
			for _, dist := range r.Dists {
				for _, ext := range ListExtensions {
					v(dist, debmirrorTag, section, translation, ext)
				}

			}
		}
	}
}

// PackageListNames lists all possible locations for package lists
func (r *Repository) PackageListNames() (names []string) {
	r.ForEachPackageSpec(func(dist, tag, section, arch, ext string) {
		///security.debian.org/dists/squeeze/updates/non-free/binary-amd64/Packages.bz2
		name := fmt.Sprintf("/dists/%s%s/%s/binary-%s/Packages%s", dist, tag, section, arch, ext)
		names = append(names, name)
	})
	return
}

// TranslationListNames lists all possible locations for translation lists
func (r *Repository) TranslationListNames() (names []string) {
	r.ForEachTranslationSpec(func(dist, tag, section, translation, ext string) {
		//ftp.uk.debian.org/debian/dists/wheezy/contrib/i18n/Translation-en.bz2
		name := fmt.Sprintf("/dists/%s%s/%s/i18n/Translation-%s%s", dist, tag, section, translation, ext)
		names = append(names, name)
	})
	return
}
