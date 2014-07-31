This is a Go package for accessing data about an [Opal](http://opal.com.au) card.

It is not anything particularly useful for users (yet);
it is currently intended for programmers to use as a building block.

To get it, you need `docutils`, `mercurial (hg command)` and `go programming language`.

* docutils

https://pypi.python.org/pypi/docutils

	mkdir docutilsetup
	cd docutilsetup
	curl -o docutils-docutils.tar.gz https://pypi.python.org/packages/source/d/docutils/docutils-0.12.tar.gz
	gunzip docutils-docutils.tar.gz 
	tar -xf docutils-docutils.tar 
	cd docutils
	sudo python setup.py install

* mercurial

http://mercurial.selenic.com/downloads

	 make            # see install targets
	 make install    # do a system-wide install
	 hg debuginstall # sanity-check setup
	 hg              # see help

* Go lanaguage

download and install GO programming lanaguage 

http://golang.org/dl/

	mkdir $HOME/go
	export GOPATH=$HOME/go	
	export GOROOT=/usr/local/go    	# default path, adjust it with your environment, if not installed in default path

	go get github.com/bw57899/opal

You will need to put auth information in $HOME/.opal that looks like

	cat ~/.opal
	{"Username":"xxx","Password":"yyy"}

Don't forget to `chmod 600 $HOME/.opal`.

I plan to make it easier to pass to this package programmatically.

To do a quick test that will access your card information
and print out its balance and recent transactions,

	go test -v github.com/dsymonds/opal
