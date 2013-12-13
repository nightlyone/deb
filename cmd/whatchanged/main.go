package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nightlyone/deb/repository"
)

func main() {
	repository.ListExtensions = []string{".gz", ""}
	log.SetFlags(0)
	_, e := parseFlags()

	switch {
	default:
		log.Fatalf("FATAL: argument parsing error, %s", e)
		fallthrough
	case e == errFlagError:
		os.Exit(2)
	case e == errFlagHelp:
		// help, which is displayed by the flags package already
		// we are done
		return
	case e == nil:
	}

	log.Println("Lets start")

	older, err := NewBase(opts.Base, opts.Old)
	if err != nil {
		log.Fatalf("FATAL: invalid base parameter %q, leading to %s", opts.Base, err)
	}
	newer, err := NewBase(opts.Base, opts.New)
	if err != nil {
		log.Fatalf("FATAL: invalid base parameter %q, leading to %s", opts.Base, err)
	}
	repo, err := repository.New()
	if err != nil {
		log.Fatalf("FATAL: cannot create repository %s", err)
	}
	repo.Archs = opts.Archs
	repo.Dists = opts.Dists
	repo.Sections = opts.Sections

	cs, err := CalcUpdates(newer, older, repo)
	if err != nil {
		log.Fatalf("FATAL: cannot calculate updates: %s", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(cs)
	if err != nil {
		log.Fatalf("FATAL: cannot write result: %s", err)
	}
}
