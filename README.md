deb
===========

Query puppetdb sane and dependency free.

[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/nightlyone/deb.png
[2]: http://travis-ci.org/nightlyone/deb


LICENSE
-------
BSD

documentation
-------------
[package documentation at go.pkgdoc.org](http://go.pkgdoc.org/github.com/nightlyone/deb)


build and install
=================

install from source
-------------------

Install [Go 1][3], either [from source][4] or [with a prepackaged binary][5].

Then run

	go get github.com/nightlyone/deb

[3]: http://golang.org
[4]: http://golang.org/doc/install/source
[5]: http://golang.org/doc/install

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

