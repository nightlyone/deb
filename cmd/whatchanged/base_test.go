package main

import (
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestNewBase(t *testing.T) {
	for i, b := range bases {
		got, err := NewBase(b.input, b.tag)

		// check error handling
		switch {
		case err != nil && b.err != err.Error():
			t.Errorf("%d: got err %s, want err %s", i, err, b.err)
			continue
		case err != nil && b.err == err.Error():
			t.Logf("%d: got expect err %s", i, err)
			continue
		case err == nil && b.err != "":
			t.Errorf("%d: got NO err, want err %s", i, b.err)
			continue
		case err == nil && b.err == "":
			// ok
		}

		if reflect.DeepEqual(b.output, got) {
			t.Logf("%d: got url %q, fs %q", i, got.u, got.fs)
		} else {
			t.Errorf("%d: got url %q, fs %q; want url %q, fs %q", i, got.u, got.fs, b.output.u, b.output.fs)
		}
	}
}

var pwd, _ = os.Getwd()

var bases = []struct {
	input  string
	tag    string
	output *Base
	err    string
}{
	{
		input: "/",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   "/",
			},
			fs: LocalDir("/"),
		},
	},
	{
		input: ".",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   pwd,
			},
			fs: LocalDir(pwd),
		},
	},
	{
		input: "",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   pwd,
			},
			fs: LocalDir(pwd),
		},
	},
	{
		input: "file:///var/lib/debmirror/ftp.uk.debian.org",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   "/var/lib/debmirror/ftp.uk.debian.org",
			},
			fs: LocalDir("/var/lib/debmirror/ftp.uk.debian.org"),
		},
	},
	{
		input: "/var/lib/debmirror/ftp.uk.debian.org",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   "/var/lib/debmirror/ftp.uk.debian.org",
			},
			fs: LocalDir("/var/lib/debmirror/ftp.uk.debian.org"),
		},
	},
	{
		input: "/var/lib/debmirror/ftp.uk.debian.org",
		tag:   "latest",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   "/var/lib/debmirror/ftp.uk.debian.org",
			},
			fs:  LocalDir("/var/lib/debmirror/ftp.uk.debian.org"),
			tag: "latest",
		},
	},
	{
		input: "file:///var/lib/debmirror/ftp.uk.debian.org",
		tag:   "latest",
		output: &Base{
			u: &url.URL{
				Scheme: "file",
				Path:   "/var/lib/debmirror/ftp.uk.debian.org",
			},
			fs:  LocalDir("/var/lib/debmirror/ftp.uk.debian.org"),
			tag: "latest",
		},
	},
	{
		input: "http://localhost/var/lib",
		output: &Base{
			u: &url.URL{
				Scheme: "http",
				Path:   "/var/lib",
				Host:   "localhost",
			},
			fs: WebDir("http://localhost/var/lib"),
		},
	},
	{
		input: "https://localhost/var/lib",
		output: &Base{
			u: &url.URL{
				Scheme: "https",
				Path:   "/var/lib",
				Host:   "localhost",
			},
			fs: WebDir("https://localhost/var/lib"),
		},
	},
	{
		input: "https://debmirror.tld/var/lib",
		output: &Base{
			u: &url.URL{
				Scheme: "https",
				Path:   "/var/lib",
				Host:   "debmirror.tld",
			},
			fs: WebDir("https://debmirror.tld/var/lib"),
		},
	},
	{
		input: "https://debmirror.tld/var/lib",
		tag:   "latest",
		output: &Base{
			u: &url.URL{
				Scheme: "https",
				Path:   "/var/lib",
				Host:   "latest.debmirror.tld",
			},
			fs: WebDir("https://latest.debmirror.tld/var/lib"),
		},
	},
	{
		input: "http://debmirror.prod.aws.jimdo-server.com/packages.vpn.jimdo-server.com",
		tag:   "prod",
		output: &Base{
			u: &url.URL{
				Scheme: "http",
				Path:   "/packages.vpn.jimdo-server.com",
				Host:   "prod.debmirror.prod.aws.jimdo-server.com",
			},
			fs: WebDir("http://prod.debmirror.prod.aws.jimdo-server.com/packages.vpn.jimdo-server.com"),
		},
	},
	// negative cases
	{
		input: "var/lib",
		err:   "parse var/lib: invalid URI for request",
	},
	{
		input: "file:///var/lib#Fragment",
		err:   "invalid url \"file:///var/lib#Fragment\", fragments are not supported",
	},
	{
		input: "/var/lib#Fragment",
		err:   "invalid url \"/var/lib#Fragment\", fragments are not supported",
	},
	{
		input: "/var/lib?query=param",
		err:   "invalid url \"/var/lib?query=param\", query parameters are not supported",
	},
	{
		input: "file:///var/lib?query=param",
		err:   "invalid url \"file:///var/lib?query=param\", query parameters are not supported",
	},
	// not yet implemented
	{
		input: "ftp:///localhost/var/lib",
		err:   "invalid url \"ftp:///localhost/var/lib\", scheme \"ftp\" not supported (only file://)",
	},
}
