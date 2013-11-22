package changelog

import (
	"bufio"
	"bytes"
	"io"
	"time"
	"unicode"
)

type Change struct {
	Source  string
	Version string
	Dist    string
	Urgency string
	Author  string
	Email   string
	Changed time.Time
	Changes []byte
}

type Changelog []Change

func Parse(r io.Reader, beforeVersion string) (Changelog, error) {
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
		if !inside {
			f := bytes.Fields(b)
			if len(f) < 4 {
				continue
			}
			s, v, d, u := f[0], f[1], f[2], f[3]
			change.Version = string(bytes.Trim(v, "()"))

			if len(beforeVersion) > 0 && change.Version == beforeVersion {
				break
			}

			change.Source = string(s)
			change.Dist = string(bytes.TrimRight(d, ";"))
			change.Urgency = string(bytes.TrimPrefix(u, []byte("urgency=")))

			inside = true
			continue
		}

		// try to parse author line
		if line := bytes.TrimPrefix(b, []byte(" -- ")); inside && len(line) < len(b) {
			start := bytes.IndexByte(line, '<')
			end := bytes.IndexByte(line, '>')
			if start < 0 || start >= end || end >= len(line) {
				continue
			}
			update, err := time.Parse(time.RFC1123Z, string(bytes.TrimSpace(line[end+1:])))
			if err != nil {
				continue
			}
			change.Author = string(bytes.TrimSpace(line[:start]))
			change.Email = string(line[start+1 : end-2])
			change.Changed = update

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
