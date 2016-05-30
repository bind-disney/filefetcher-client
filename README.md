### Description

This is a simple FTP-like command-line client written in Go.

### Installation

`go get github.com/bind-disney/filefetcher-client`

### Usage

You need to specify server domain name or IP address through `-H` option, server port through `-P` option
(by default it is 27518) and target directory (by default it is `downloads`), where files will be stored
by default, if not specified while downloading a file.

When the connection to the server is successfully established, you can use simple FTP-like commands:

* `list` / `ls` for listing files in the current remote directory
* `chdir` / `cd` for changing current remote directory
* `get` / `download` / `fetch` for downloading a file
* `help` / `h` for showing this usage help
* `exit` / `quit` / `q` for quiting the application

##### Warning

* Directory downloading is not supported yet
* File names with spaces are not supported yet

### References

#### Libraries

- [A Commander for modern Go CLI interactions](https://github.com/spf13/cobra)
- [A small package for building command line apps in Go](https://github.com/urfave/cli)
- [The easy way to build Golang command-line application](https://github.com/tcnksm/gcli)
- [RGB colors for your terminal](https://github.com/aybabtme/rgbterm)
- [like python-sh, for easy call shell with golang](https://github.com/codeskyblue/go-sh)
- [Pure Go termbox implementation](https://github.com/nsf/termbox-go)
- [A command-line arguments parser that will make you smile](https://github.com/docopt/docopt.go)
- [A Go (golang) command line and flag parser](https://github.com/alecthomas/kingpin)
- [Quake-style colour formatting for Unix terminals](https://github.com/alecthomas/colour)
- [The libs club / development / cli](http://libs.club/golang/developement/cli)

#### Posts

* [CLI: Command Line Programming with Go](http://thenewstack.io/cli-command-line-programming-with-go/)
- [Why codegangsta's cli package is the bomb, and you should use it](http://nathanleclaire.com/blog/2014/08/31/why-codegangstas-cli-package-is-the-bomb-and-you-should-use-it/)
- [A Modern CLI Commander For Go](http://spf13.com/post/announcing-cobra/)
- [Golang: Implementing subcommands for command line applications](http://blog.ralch.com/tutorial/golang-subcommands/)
- [The easy way to start building Golang command-line application](https://coderwall.com/p/_1d7nq/the-easy-way-to-start-building-golang-command-line-application)
- [cli.go — better command-line applications in Go](https://www.progville.com/go/cli-go-better-command-line-applications-go/)
- [Python, Ruby, and Golang: A Command-Line Application Comparison](https://realpython.com/blog/python/python-ruby-and-golang-a-command-line-application-comparison/)
- [On Distributing Command line Applications: Why I switched from Ruby to Go](https://codegangsta.io/blog/2013/07/21/creating-cli-applications-in-go/)
- [Creating console applications in Golang](http://blog.dutchcoders.io/creating-cli-applications-in-golang/)
- [Building Command Line Tools with Golang](http://www.slideshare.net/nanzou/building-command-line-tools-with-golang)
- [Fatally exiting a command line tool, with grace](https://medium.com/@matryer/golang-advent-calendar-day-three-fatally-exiting-a-command-line-tool-with-grace-874befeb64a4#.mj8lnzods)
- [Building helpful CLI tools with Go and Kingpin](https://developer.atlassian.com/blog/2016/03/building-helpful-golang-cli-tools/)
