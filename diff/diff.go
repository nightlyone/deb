// Package diff calculates changes in Debian package lists.
package diff

import (
	"io"
	"sort"

	"github.com/stapelberg/godebiancontrol"
)

// A List of packages in a format useful for postprocessing
type List struct {
	Version map[string]string // Version for each package
	Package map[string]string // Package to source mapping
}

// NewList reads a package list from r
func NewList(r io.Reader) (*List, error) {
	pp, err := godebiancontrol.Parse(r)
	if err != nil {
		return nil, err
	}
	l := &List{
		Version: map[string]string{},
		Package: map[string]string{},
	}
	for _, e := range pp {
		p := e["Package"]
		v := e["Version"]
		if p == "" || v == "" {
			continue
		}
		s, ok := e["Source"]
		if !ok {
			s = e["Package"]
		}
		l.Package[p] = s
		l.Version[p] = v
	}
	return l, nil
}

// Package succintly describes a package
type Package struct {
	Name, Source, Version string
}

// Sort package list by Source package of Package
type bySource []*Package

func (p bySource) Len() int           { return len(p) }
func (p bySource) Less(i, j int) bool { return p[i].Source < p[j].Source }
func (p bySource) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Changes are returned as list of added, removed and updated packages.
func Changes(newer, older *List) (added, removed, updated []*Package) {
	for pkg, src := range newer.Package {
		old_src, exists := older.Package[pkg]
		after := newer.Version[pkg]
		if !exists {
			added = append(added, &Package{
				Name:    pkg,
				Source:  src,
				Version: after,
			})
			continue
		}
		before := older.Version[pkg]
		if after != before || src != old_src {
			updated = append(updated, &Package{
				Name:    pkg,
				Source:  src,
				Version: after,
			})
		}
	}
	for pkg, old_src := range older.Package {
		_, exists := newer.Package[pkg]
		before := older.Version[pkg]
		if !exists {
			removed = append(removed, &Package{
				Name:    pkg,
				Source:  old_src,
				Version: before,
			})
		}
	}
	return
}

// Update describe an updated Package and updates related to it
type Update struct {
	Package string
	Version string
	Related []*Package
}

// CompressUpdates detects updates of the same source package
// to the same version and considers them related.
// This is useful for display of changes to the user.
// e.g an update of 29 packages is actually just 5 different sets of updates,
// where each has a bunch of related packages compiled from the same source package update.
func CompressUpdates(changes []*Package) (update []*Update) {
	if len(changes) == 0 {
		return nil
	}
	if len(changes) == 1 {
		return append(update, &Update{
			Package: changes[0].Name,
			Version: changes[0].Version,
		})
	}
	sort.Sort(bySource(changes))
	start := 0
	p := changes[start]
	s := p
	for i, u := range changes {
		if u.Source == p.Source && u.Version == p.Version {
			if u.Name == s.Source {
				s = u
			}
			continue
		}
		update = append(update, &Update{
			Package: s.Name,
			Version: s.Version,
			Related: changes[start:i],
		})
		start = i
		p = changes[start]
		s = p
	}
	update = append(update, &Update{
		Package: s.Name,
		Version: s.Version,
		Related: changes[start:],
	})
	return update
}
