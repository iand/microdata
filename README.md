# microdata
A microdata parser in Go

See [http://www.w3.org/TR/microdata/](http://www.w3.org/TR/microdata/) for more information about Microdata

## Installation

Simply run

	go get github.com/iand/microdata

Documentation is at [http://godoc.org/github.com/iand/microdata](http://godoc.org/github.com/iand/microdata)


## Usage

Example of parsing a string containing HTML:

	package main

	import (
		"github.com/iand/microdata"
		"net/url"
		"strings"
	)

	func main() {
		html := `<div itemscope>
		 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
		</div>`

		baseUrl, _ := url.Parse("http://example.com/")
		p := microdata.NewParser(strings.NewReader(html), baseUrl)

		data, err := p.Parse()
		if err != nil {
			panic(err)
		}

		println("Name: ", data.Items[0].Properties["name"][0].(string))
	}

Extract microdata from a webpage and print the result as JSON

	package main

	import (
	    "bytes"
	    "io/ioutil"
	    "net/http"
	    "net/url"
	    "os"

	    "github.com/iand/microdata"
	)

	func main() {

	    baseUrl, _ := url.Parse("http://www.designhive.com/blog/using-schemaorg-microdata")

	    resp, _ := http.Get(baseUrl.String())
	    defer resp.Body.Close()

	    html, _ := ioutil.ReadAll(resp.Body)

	    p := microdata.NewParser(bytes.NewReader(html), baseUrl)

	    data, _ := p.Parse()

	    json, _ := data.Json()
	    os.Stdout.Write(json)
	}


## Authors

* [Ian Davis](http://github.com/iand) - <http://iandavis.com/>


## Contributors


## Contributing

* Do submit your changes as a pull request
* Do your best to adhere to the existing coding conventions and idioms.
* Do run `go fmt` on the code before committing
* Do feel free to add yourself to the [`CREDITS`](CREDITS) file and the
  corresponding Contributors list in the the [`README.md`](README.md).
  Alphabetical order applies.
* Don't touch the [`AUTHORS`](AUTHORS) file. An existing author will add you if
  your contributions are significant enough.
* Do note that in order for any non-trivial changes to be merged (as a rule
  of thumb, additions larger than about 15 lines of code), an explicit
  Public Domain Dedication needs to be on record from you. Please include
  a copy of the statement found in the [`WAIVER`](WAIVER) file with your pull request

## License

This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.
