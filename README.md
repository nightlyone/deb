deb
===========

Various toolings around debian package lists, mirrors and changelogs

[![Build Status][1]][2]
[![GoDoc][3]][4]




LICENSE
-------
BSD

documentation
-------------
package documentation at [godoc.org](http://godoc.org/github.com/nightlyone/deb)
or [gowalker.org](http://gowalker.org/github.com/nightlyone/deb)


build and install
=================

install from source
-------------------

Install [Go 1][5], either [from source][6] or [with a prepackaged binary][7].

Then run

	go get github.com/nightlyone/deb


LICENSE
-------
BSD

documentation
-------------

contributing
============

Contributions are welcome. Please open an issue or send me a pull request for a dedicated branch.
Make sure the git commit hooks show it works.

git commit hooks
-----------------------
enable commit hooks via

        cd .git ; rm -rf hooks; ln -s ../git-hooks hooks ; cd ..

References:
[1]: https://secure.travis-ci.org/nightlyone/deb.png
[2]: http://travis-ci.org/nightlyone/deb
[3]: https://godoc.org/github.com/nightlyone/deb?status.png
[4]: https://godoc.org/github.com/nightlyone/deb
[5]: http://golang.org
[6]: http://golang.org/doc/install/source
[7]: http://golang.org/doc/install

