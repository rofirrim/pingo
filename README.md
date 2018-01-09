# Welcome to go.pinchito.com

[![CircleCI](https://circleci.com/gh/rofirrim/pingo/tree/master.svg?style=svg)](https://circleci.com/gh/rofirrim/pingo/tree/master)

Web app for go.pinchito.com. Done using [Revel](http://revel.github.io/) a
high-productivity web framework for the [Go language](http://www.golang.org/) 

## Quick start

Make sure you have a correct `$GOPATH` set (e.g. add `export GOPATH=$HOME/Go` to your `.bashrc`)

    $ go get github.com/revel/revel
    $ go get github.com/revel/cmd/revel
    $ go get github.com/go-sql-driver/mysql
    $ git clone https://github.com/rofirrim/pingo src/pingo

To set up the DB

    $ cd $GOPATH/src/pingo/conf
    $ cp settings.json.example settings.json

and then edit `settings.json`.

The DB used is MySQL. Ask me for a dump of the DB in SQL format, otherwise
the application will not work.

Local server for development

    $ export PATH=$GOPATH/bin:$PATH
    $ revel run pingo

Now connect to localhost:9000

### Follow the guidelines to start developing your application:

* The README file created within your application.
* The [Getting Started with Revel](http://revel.github.io/tutorial/index.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/samples/index.html).
* The [API documentation](https://godoc.org/github.com/revel/revel).
