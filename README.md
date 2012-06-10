microdata - a microdata parser in Go

INSTALLATION
============

Simply run

	go get github.com/iand/microdata

Documentation is at [http://go.pkgdoc.org/github.com/iand/microdata](http://go.pkgdoc.org/github.com/iand/microdata)


USAGE
=====

Example of parsing a string containing HTML:

	include (
		"net/url"
		"strings"
	)
	html = `<div itemscope>
	 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
	</div>`

	baseUrl, _ := url.Parse("http://example.com/")
	p := NewParser(strings.NewReader(html), baseUrl)

	data, err := p.Parse()
	if err != nil {
		t.Errorf("Expected no error but got %d", err)
	}

	println("Name: ", data.items[0].properties["name"][0]