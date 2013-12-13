package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileSystem abstracts an collection which can return io.ReadCloser by a unqiue name
type FileSystem interface {
	Open(name string) (io.ReadCloser, error)
}

// LocalDir is a directory in the local filesystem
type LocalDir string

// Open a local file.
// NOTE: It will never return a directory.
func (d LocalDir) Open(name string) (io.ReadCloser, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.ContainsAny(name, "\x00?#") {
		return nil, errors.New("invalid character in file path")
	}
	dir := string(d)
	if dir == "" {
		dir = "."
	}
	f, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
	if err != nil {
		return nil, err
	}
	if fi, err := f.Stat(); err != nil || fi.IsDir() {
		f.Close()
		return nil, err
	}
	return f, nil
}

// WebDir is a directory in via http/https
type WebDir string

// Open streams directly from http/htpps.
// NOTE: It will never return a directory
func (d WebDir) Open(name string) (io.ReadCloser, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.ContainsAny(name, "\x00?#:") {
		return nil, errors.New("invalid character in file path")
	}
	base, err := url.ParseRequestURI(string(d))
	if err != nil {
		return nil, err
	}

	u, err := base.Parse(path.Join(base.Path, path.Clean("/"+name)))
	if err != nil {
		return nil, err
	}

	log.Println("downloading", u)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 300 {
		return resp.Body, nil
	}
	return nil, fmt.Errorf("download of %q failed with %d: %s", u, resp.StatusCode, http.StatusText(resp.StatusCode))

}
