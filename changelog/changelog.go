// Package changelog provides helpers to parse a typical
// Debian changelog file found in packages generated from source packages.
// see http://www.debian.org/doc/debian-policy/ch-source.html for details on the format
package changelog

import (
	"bufio"
	"bytes"
	"io"
	"time"
	"unicode"
)

// Change describes a single changeset between 2 versions
type Change struct {
	Source  string
	Version string
	Dists   []string
	Urgency string
	Author  string
	Email   string
	Changed time.Time
	Changes []byte
}

// Changelog describes a full or partial changelog
type Changelog []Change

func (change *Change) parseVersionLine(b []byte) (ok bool) {
	split := bytes.SplitN(b, []byte{';'}, 2)
	f := bytes.Fields(split[0])
	if len(f) < 3 || len(split) != 2 {
		return false
	}
	u := split[1]
	s, v, d := f[0], f[1], f[2:]

	change.Version = string(bytes.Trim(v, "()"))
	change.Source = string(s)
	change.Urgency = string(bytes.TrimPrefix(bytes.TrimSpace(u), []byte("urgency=")))

	change.Dists = nil
	for _, dist := range d {
		change.Dists = append(change.Dists, string(dist))
	}

	return change.Version != "" &&
		change.Source != "" &&
		len(change.Dists) > 0 &&
		change.Urgency != ""
}

func (change *Change) parseChangedByLine(b []byte) (ok bool) {
	if !bytes.HasPrefix(b, []byte(" -- ")) {
		return false
	}
	line := bytes.TrimPrefix(b, []byte(" -- "))
	start := bytes.IndexByte(line, '<')
	end := bytes.IndexByte(line, '>')
	if start < 0 || start >= end || end >= len(line) {
		return false
	}
	update, err := time.Parse(time.RFC1123Z, string(bytes.TrimSpace(line[end+1:])))
	if err != nil {
		return false
	}
	change.Author = string(bytes.TrimSpace(line[:start]))
	change.Email = string(line[start+1 : end-2])
	change.Changed = update

	return true
}

// Parse a debian changelog from r for any changes happening later than afterVersion
func Parse(r io.Reader, afterVersion string) (Changelog, error) {
	scanner := bufio.NewScanner(r)
	changelog := make(Changelog, 0, 5)
	change := Change{}
	inside := false

	for scanner.Scan() {
		b := bytes.TrimRightFunc(scanner.Bytes(), unicode.IsSpace)
		if b2 := bytes.TrimSpace(b); len(b2) < len(b) && !inside {
			b = b2
		}
		if len(b) == 0 {
			if inside {
				change.Changes = append(change.Changes, '\n')
			}
			continue
		}

		if !inside && change.parseVersionLine(b) {
			if len(afterVersion) > 0 && change.Version == afterVersion {
				break
			}
			inside = true
			continue
		}

		if inside && change.parseChangedByLine(b) {
			changelog = append(changelog, change)
			change = Change{}
			inside = false
			continue
		}

		change.Changes = append(change.Changes, b...)
		change.Changes = append(change.Changes, '\n')
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return changelog, nil
}
