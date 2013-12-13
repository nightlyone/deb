package main

import (
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"

	"github.com/nightlyone/deb/diff"
	"github.com/nightlyone/deb/repository"
)

// Base is a base for futher repository retrieval actions
type Base struct {
	u   *url.URL
	fs  FileSystem
	tag string
}

// NewBase creates a base for further repository retrieval actions
func NewBase(base string, tag string) (*Base, error) {
	if base == "." || base == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		base = pwd
	}
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	if u.Fragment != "" {
		return nil, fmt.Errorf("invalid url %q, fragments are not supported", base)
	}

	// Need to parse a second time, bacause we want it absolute, but url.ParseRequestURI supports no fragments
	u, err = url.ParseRequestURI(base)
	if err != nil {
		return nil, err
	}

	if u.RawQuery != "" {
		return nil, fmt.Errorf("invalid url %q, query parameters are not supported", base)
	}

	u, _ = u.Parse(u.Path)

	switch u.Scheme {
	case "":
		u.Scheme = "file"
		fallthrough
	case "file":
		return &Base{
			u:   u,
			fs:  LocalDir(u.Path),
			tag: tag,
		}, nil
	case "http", "https":
		if u.Host == "" {
			return nil, fmt.Errorf("invalid remote url %q, host part should not be empty", base)
		}
		if tag != "" {
			u.Host = tag + "." + u.Host
		}
		return &Base{
			u:  u,
			fs: WebDir(u.String()),
		}, nil
	default:
		return nil, fmt.Errorf("invalid url %q, scheme %q not supported (only file://)", base, u.Scheme)
	}
}

func (b *Base) taggedRepository(r *repository.Repository) *repository.Repository {
	rep := *r
	rep.DebMirrorTag = b.tag
	return &rep
}

func ReadList(r io.ReadCloser, ext string) (list *diff.List, err error) {
	var reader io.Reader
	defer r.Close()
	switch ext {
	case ".gz":
		reader, err = gzip.NewReader(r)
		if err != nil {
			return
		}
		defer reader.(io.Closer).Close()
	case ".bz2":
		reader = bzip2.NewReader(r)
	case "":
		reader = r
	}
	return diff.NewList(reader)
}

func (b *Base) loadMergableLists(mergable []string) (combined *diff.List, err error) {
	lists := map[string]*diff.List{}
	for _, fn := range mergable {
		key := path.Join(path.Dir(fn), "Packages")
		list := lists[key]
		if list != nil {
			continue
		}
		f, err := b.fs.Open(fn)
		if err != nil {
			return nil, err
		}
		list, err = ReadList(f, path.Ext(fn))
		if err != nil {
			return nil, err
		}
		lists[key] = list
	}

	for _, v := range lists {
		if combined == nil {
			combined = v
		} else {
			combined.Merge(v)
		}
	}

	return
}
func (b *Base) Packages(repo *repository.Repository) (combined map[string]*diff.List, err error) {
	tagged := b.taggedRepository(repo)
	combined = map[string]*diff.List{}
	for key, mergable := range tagged.KeySets() {
		c, err := b.loadMergableLists(mergable)
		if err != nil {
			return nil, err
		}
		combined[key] = c
	}
	return combined, nil
}

type Changeset struct {
	Added,
	Removed []*diff.Package
	Updated []*diff.Update
}

var errIncompatible = errors.New("Incompatible package lists")

func CalcUpdates(newer, older *Base, repo *repository.Repository) (map[string]Changeset, error) {
	o, err := older.Packages(repo)
	if err != nil {
		return nil, err
	}
	n, err := newer.Packages(repo)
	if err != nil {
		return nil, err
	}
	if len(n) != len(o) {
		return nil, errIncompatible
	}
	changes := map[string]Changeset{}
	for key, list := range n {
		a, r, u := diff.Changes(list, o[key])
		if len(a) == 0 && len(r) == 0 && len(u) == 0 {
			continue
		}
		changes[key] = Changeset{
			Added:   a,
			Removed: r,
			Updated: diff.CompressUpdates(u),
		}
	}
	return changes, nil
}
