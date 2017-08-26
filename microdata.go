/*
  This is free and unencumbered software released into the public domain. For more
  information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

// Package microdata provides types and functions for paring microdata from web pages.
// See http://www.w3.org/TR/microdata/ for more information about Microdata
package microdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type ValueList []interface{}
type PropertyMap map[string]ValueList

// Item represents a microdata item
type Item struct {
	Properties PropertyMap `json:"properties"`
	Types      []string    `json:"type,omitempty"`
	ID         string      `json:"id,omitempty"`
}

// NewItem creates a new microdata item
func NewItem() *Item {
	return &Item{
		Properties: make(PropertyMap, 0),
		Types:      make([]string, 0),
	}
}

// AddString adds a string type item property value
func (i *Item) AddString(property string, value string) {
	i.Properties[property] = append(i.Properties[property], value)
}

// AddItem adds an Item type item property value
func (i *Item) AddItem(property string, value *Item) {
	i.Properties[property] = append(i.Properties[property], value)
}

// AddType adds a type to the item
func (i *Item) AddType(value string) {
	i.Types = append(i.Types, value)
}

// Microdata represents a set of microdata items
type Microdata struct {
	Items []*Item `json:"items"`
}

// NewMicrodata creates a new microdata set
func NewMicrodata() *Microdata {
	return &Microdata{
		Items: make([]*Item, 0),
	}
}

// AddItem adds an item to the microdata set
func (m *Microdata) AddItem(value *Item) {
	m.Items = append(m.Items, value)
}

// JSON converts the microdata set to JSON
func (m *Microdata) JSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Parser is an HTML parser that extracts microdata
type Parser struct {
	r               io.Reader
	data            *Microdata
	base            *url.URL
	identifiedNodes map[string]*html.Node
}

// NewParser creates a new parser for extracting microdata
// r is a reader over an HTML document
// base is the base URL for resolving relative URLs
func NewParser(r io.Reader, base *url.URL) *Parser {
	return &Parser{
		r:    r,
		data: NewMicrodata(),
		base: base,
	}
}

// Parse the document and return a Microdata set
func (p *Parser) Parse() (*Microdata, error) {
	tree, err := html.Parse(p.r)
	if err != nil {
		return nil, err
	}

	topLevelItemNodes := make([]*html.Node, 0)
	p.identifiedNodes = make(map[string]*html.Node, 0)

	walk(tree, func(n *html.Node) {
		if n.Type == html.ElementNode {
			if _, exists := getAttr("itemscope", n); exists {
				if _, exists := getAttr("itemprop", n); !exists {
					topLevelItemNodes = append(topLevelItemNodes, n)
				}
			}

			if id, exists := getAttr("id", n); exists {
				p.identifiedNodes[id] = n
			}
		}
	})

	for _, node := range topLevelItemNodes {
		item := NewItem()
		p.data.Items = append(p.data.Items, item)
		if itemtypes, exists := getAttr("itemtype", node); exists {
			for _, itemtype := range strings.Split(strings.TrimSpace(itemtypes), " ") {
				itemtype = strings.TrimSpace(itemtype)
				if itemtype != "" {
					item.Types = append(item.Types, itemtype)
				}
			}
			// itemid only valid when itemscope and itemtype are both present
			if itemid, exists := getAttr("itemid", node); exists {
				if parsedUrl, err := p.base.Parse(itemid); err == nil {
					item.ID = parsedUrl.String()
				}
			}

		}

		if itemrefs, exists := getAttr("itemref", node); exists {
			for _, itemref := range strings.Split(strings.TrimSpace(itemrefs), " ") {
				itemref = strings.TrimSpace(itemref)

				if refnode, exists := p.identifiedNodes[itemref]; exists {
					p.readItem(item, refnode)
				}
			}
		}

		for child := node.FirstChild; child != nil; {
			p.readItem(item, child)
			child = child.NextSibling
		}
	}

	return p.data, nil
}

func (p *Parser) readItem(item *Item, node *html.Node) {
	if itemprop, exists := getAttr("itemprop", node); exists {
		if _, exists := getAttr("itemscope", node); exists {
			subitem := NewItem()

			if itemrefs, exists := getAttr("itemref", node); exists {
				for _, itemref := range strings.Split(strings.TrimSpace(itemrefs), " ") {
					itemref = strings.TrimSpace(itemref)

					if refnode, exists := p.identifiedNodes[itemref]; exists {
						p.readItem(subitem, refnode)
					}
				}
			}

			for child := node.FirstChild; child != nil; {
				p.readItem(subitem, child)
				child = child.NextSibling
			}

			for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
				propertyName = strings.TrimSpace(propertyName)
				if propertyName != "" {
					item.AddItem(propertyName, subitem)
				}
			}

			return

		}

		var propertyValue string

		switch node.DataAtom {
		case atom.Meta:
			if val, exists := getAttr("content", node); exists {
				propertyValue = val
			}
		case atom.Audio, atom.Embed, atom.Iframe, atom.Img, atom.Source, atom.Track, atom.Video:
			if urlValue, exists := getAttr("src", node); exists {
				if parsedUrl, err := p.base.Parse(urlValue); err == nil {
					propertyValue = parsedUrl.String()
				}

			}
		case atom.A, atom.Area, atom.Link:
			if urlValue, exists := getAttr("href", node); exists {
				if parsedUrl, err := p.base.Parse(urlValue); err == nil {
					propertyValue = parsedUrl.String()
				}
			}
		case atom.Object:
			if urlValue, exists := getAttr("data", node); exists {
				propertyValue = urlValue
			}
		case atom.Data, atom.Meter:
			if urlValue, exists := getAttr("value", node); exists {
				propertyValue = urlValue
			}
		case atom.Time:
			if urlValue, exists := getAttr("datetime", node); exists {
				propertyValue = urlValue
			}

		default:
			var text bytes.Buffer
			walk(node, func(n *html.Node) {
				if n.Type == html.TextNode {
					text.WriteString(n.Data)
				}

			})
			propertyValue = text.String()
		}

		if len(propertyValue) > 0 {
			for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
				propertyName = strings.TrimSpace(propertyName)
				if propertyName != "" {
					item.AddString(propertyName, propertyValue)
				}
			}
		}

	}

	for child := node.FirstChild; child != nil; {
		p.readItem(item, child)
		child = child.NextSibling
	}

}

func getAttr(name string, node *html.Node) (string, bool) {
	for _, a := range node.Attr {
		if a.Key == name {
			return a.Val, true
		}
	}
	return "", false
}

func walk(parent *html.Node, fn func(n *html.Node)) {
	if parent == nil {
		return
	}
	fn(parent)

	for child := parent.FirstChild; child != nil; {
		walk(child, fn)
		child = child.NextSibling
	}
}
