microdata - a microdata parser in Go

INSTALLATION
============

Simply run

	go get github.com/iand/microdata

Documentation is at [http://go.pkgdoc.org/github.com/iand/microdata](http://go.pkgdoc.org/github.com/iand/microdata)


USAGE
=====

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
		"github.com/iand/microdata"
		"io/ioutil"
		"net/http"
		"net/url"
		"os"
	)

	func main() {

		baseUrl, _ := url.Parse("http://tagger.steve.museum/steve/object/44863?offset=6")

		resp, _ := http.Get(baseUrl.String())
		defer resp.Body.Close()

		html, _ := ioutil.ReadAll(resp.Body)

		p := microdata.NewParser(bytes.NewReader(html), baseUrl)

		data, _ := p.Parse()

		json, _ := data.Json()
		os.Stdout.Write(json)
	}		


LICENSE
=======
This code and associated documentation is in the public domain.

To the extent possible under law, Ian Davis has waived all copyright
and related or neighboring rights to this file. This work is published 
from the United Kingdom. 
