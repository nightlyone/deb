package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"

	"github.com/stapelberg/godebiancontrol"

	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Package struct {
	Name, Version, Arch string

	filename   string
	readone    bool
	archive    *tar.Reader
	fieldcache godebiancontrol.Paragraph
}

func NewPackage(filename string) *Package {
	basename := filepath.Base(filename)
	split := strings.Split(basename, "_")
	if len(split) != 3 {
		return nil
	}
	return &Package{
		Name:     split[0],
		Version:  strings.Replace(split[1], "%", ":", -1),
		Arch:     strings.TrimSuffix(split[2], ".deb"),
		filename: filename,
	}

}

func (p *Package) updateFieldCache() error {
	cmd := exec.Command("dpkg-deb", "--field", p.filename)

	o, err := cmd.Output()
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(o)
	pp, err := godebiancontrol.Parse(b)
	if err != nil {
		return err
	}

	if len(pp) != 1 {
		return os.ErrInvalid

	}

	p.fieldcache = pp[0]
	return nil
}

func (p *Package) Source() (string, error) {
	if p.fieldcache == nil {
		err := p.updateFieldCache()
		if err != nil {
			return "", err
		}
	}
	return p.fieldcache["Source"], nil
}

func (p *Package) archiveOpen() (*tar.Reader, error) {
	if p.archive != nil {
		return p.archive, nil
	}
	cmd := exec.Command("dpkg-deb", "--fsys-tarfile", p.filename)
	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func(c *exec.Cmd) {
		if err := c.Wait(); err != nil {
			log.Println("Cannot stream package in tar mode.", err)
		}
	}(cmd)
	p.archive = tar.NewReader(r)
	p.readone = false
	return p.archive, nil
}

func (p *Package) FindFile(filename string) (io.ReadCloser, error) {
	archive, err := p.archiveOpen()
	if err != nil {
		return nil, err
	}

	rewind := !p.readone
	for {
		hdr, err := archive.Next()
		p.readone = true
		// end of tar archive
		if err == io.EOF {
			// if we started in the middle, try to start from the beginning again
			if rewind {
				p.archive = nil

				archive, err = p.archiveOpen()
				if err != nil {
					return nil, err
				}

				// and now back to searching...
				rewind = false
				continue
			}
			return nil, os.ErrNotExist
		}
		if err != nil {
			// b0rken tar archive
			return nil, err
		}
		if fi := hdr.FileInfo(); fi.Name() == filename {
			if fi.Mode().IsRegular() {
				return gzip.NewReader(archive)
			} else {
				return nil, os.ErrInvalid
			}
		}
	}
	panic("Never reached")
}
