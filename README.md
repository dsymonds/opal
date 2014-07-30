This is a Go package for accessing data about an [Opal](http://opal.com.au) card.

It is not anything particularly useful for users (yet);
it is currently intended for programmers to use as a building block.

To get it, you need download and install GO programming lanaguage first

http://golang.org/dl/

mkdir $HOME/go
export GOPATH=$HOME/go

?? exec: "hg": executable file not found in $PATH


	go get github.com/bw57899/opal

You will need to put auth information in $HOME/.opal that looks like

	{"Username":"xxx","Password":"yyy"}

Don't forget to `chmod 600 $HOME/.opal`.
I plan to make it easier to pass to this package programmatically.

To do a quick test that will access your card information
and print out its balance and recent transactions,

	go test -v github.com/dsymonds/opal
