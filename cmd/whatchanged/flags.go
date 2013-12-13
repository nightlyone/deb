package main

import (
	"errors"
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Base     string   `long:"url" default:"." description:"the base url for the repository mirror. file:// can be omitted"`
	Old      string   `long:"old" required:"true" description:"old version tag in repository mirror"`
	New      string   `long:"new" required:"true" description:"new version tag in repository mirror"`
	Output   string   `long:"output" default:"json" description:"format of the output"`
	Archs    []string `long:"arch" required:"true" description:"architecture to check"`
	Dists    []string `long:"dist" required:"true" description:"dist to check"`
	Sections []string `long:"section" required:"true" description:"section to check"`
}

var errFlagError = errors.New("error handled by flag package, should not be displayed")
var errFlagHelp = errors.New("help displayed, not an error")

// FlagConstraintError documents errors due to complex constraints on commandline arguments
type FlagConstraintError struct {
	Constraint string
}

func (c *FlagConstraintError) Error() string {
	return fmt.Sprintf(c.Constraint)
}

func validateOptionConstraints() (err error) {
	if len(opts.Sections) == 0 {
		opts.Sections = []string{"main", "contrib", "non-free", "main/debian-installer"}
	}
	if len(opts.Dists) == 0 {
		opts.Dists = []string{"stable"}
	}
	if len(opts.Archs) == 0 {
		opts.Archs = []string{"all"}
	}
	switch opts.Output {
	case "json":
	default:
		return &FlagConstraintError{Constraint: "output option allows only json"}
	}
	return nil
}

func parseFlags() ([]string, error) {
	p := flags.NewParser(&opts, flags.Default)

	// display nice usage message
	p.Usage = "[OPTIONS]... \n\nDisplay what changed in a repository"

	args, err := p.Parse()
	if err != nil {
		// --help is not an error
		if e, ok := err.(*flags.Error); ok {
			switch e.Type {
			case flags.ErrHelp:
				return nil, errFlagHelp
			case flags.ErrRequired,
				flags.ErrExpectedArgument,
				flags.ErrUnknownFlag,
				flags.ErrUnknownGroup,
				flags.ErrNoArgumentForBool,
				flags.ErrShortNameTooLong:
				return nil, errFlagError
			case flags.ErrUnknown:
				return nil, err
			default:
				// catch the programming error
				panic(err)
			}
			return nil, nil
		}
		return nil, err
	}

	if err := validateOptionConstraints(); err != nil {
		return nil, err
	}

	return args, nil
}
