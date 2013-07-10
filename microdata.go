/*
  This is free and unencumbered software released into the public domain. For more
  information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

// A package for parsing microdata
// See http://www.w3.org/TR/microdata/ for more information about Microdata
package microdata

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

type ValueList []interface{}
type PropertyMap map[string]ValueList

// Represents a microdata item
type Item struct {
	Properties PropertyMap `json:"properties"`
	Types      []string    `json:"type,omitempty"`
	ID         string      `json:"id,omitempty"`
}

// Create a new microdata item
func NewItem() *Item {
	return &Item{
		Properties: make(PropertyMap, 0),
		Types:      make([]string, 0),
	}
}

// Add a string type item property value
func (self *Item) AddString(property string, value string) {
	self.Properties[property] = append(self.Properties[property], value)
}

// Add an Item type item property value
func (self *Item) AddItem(property string, value *Item) {
	self.Properties[property] = append(self.Properties[property], value)
}

// Add a type to the item
func (self *Item) AddType(value string) {
	self.Types = append(self.Types, value)
}

// Represents a set of microdata items
type Microdata struct {
	Items []*Item `json:"items"`
}

// Create a new microdata set
func NewMicrodata() *Microdata {
	return &Microdata{
		Items: make([]*Item, 0),
	}
}

// Add an item to the microdata set
func (self *Microdata) AddItem(value *Item) {
	self.Items = append(self.Items, value)
}

// Convert the microdata set to JSON
func (self *Microdata) Json() ([]byte, error) {
	b, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// An HTML parser that extracts microdata
type Parser struct {
	p               *h5.Tree
	data            *Microdata
	base            *url.URL
	identifiedNodes map[string]*html.Node
}

// Create a new parser for extracting microdata
// r is a reader over an HTML document
// base is the base URL for resolving relative URLs
func NewParser(r io.Reader, base *url.URL) *Parser {
	p, _ := h5.New(r)

	return &Parser{
		p:    p,
		data: NewMicrodata(),
		base: base,
	}
}

// Parse the document and return a Microdata set
func (self *Parser) Parse() (*Microdata, error) {
	tree := self.p

	topLevelItemNodes := make([]*html.Node, 0)
	self.identifiedNodes = make(map[string]*html.Node, 0)

	tree.Walk(func(n *html.Node) {
		if _, exists := getAttr("itemscope", n); exists {
			if _, exists := getAttr("itemprop", n); !exists {
				topLevelItemNodes = append(topLevelItemNodes, n)
			}
		}

		if id, exists := getAttr("id", n); exists {
			self.identifiedNodes[id] = n
		}
	})

	for _, node := range topLevelItemNodes {
		item := NewItem()
		self.data.Items = append(self.data.Items, item)
		if itemtypes, exists := getAttr("itemtype", node); exists {
			for _, itemtype := range strings.Split(strings.TrimSpace(itemtypes), " ") {
				itemtype = strings.TrimSpace(itemtype)
				if itemtype != "" {
					item.Types = append(item.Types, itemtype)
				}
			}
			// itemid only valid when itemscope and itemtype are both present
			if itemid, exists := getAttr("itemid", node); exists {
				if parsedUrl, err := self.base.Parse(itemid); err == nil {
					item.ID = parsedUrl.String()
				}
			}

		}

		if itemrefs, exists := getAttr("itemref", node); exists {
			for _, itemref := range strings.Split(strings.TrimSpace(itemrefs), " ") {
				itemref = strings.TrimSpace(itemref)

				if refnode, exists := self.identifiedNodes[itemref]; exists {
					self.readItem(item, refnode)
				}
			}
		}

		for child := node.FirstChild; child != nil; {
			self.readItem(item, child)
			child = child.NextSibling
		}
	}

	return self.data, nil
}

func (self *Parser) readItem(item *Item, node *html.Node) {
	if itemprop, exists := getAttr("itemprop", node); exists {
		if _, exists := getAttr("itemscope", node); exists {
			subitem := NewItem()

			if itemrefs, exists := getAttr("itemref", node); exists {
				for _, itemref := range strings.Split(strings.TrimSpace(itemrefs), " ") {
					itemref = strings.TrimSpace(itemref)

					if refnode, exists := self.identifiedNodes[itemref]; exists {
						self.readItem(subitem, refnode)
					}
				}
			}

			for child := node.FirstChild; child != nil; {
				self.readItem(subitem, child)
				child = child.NextSibling
			}

			for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
				propertyName = strings.TrimSpace(propertyName)
				if propertyName != "" {
					item.AddItem(propertyName, subitem)
				}
			}

			return

		} else {
			var propertyValue string

			switch node.DataAtom {
			case atom.Meta:
				if val, exists := getAttr("content", node); exists {
					propertyValue = val
				}
			case atom.Audio, atom.Embed, atom.Iframe, atom.Img, atom.Source, atom.Track, atom.Video:
				if urlValue, exists := getAttr("src", node); exists {
					if parsedUrl, err := self.base.Parse(urlValue); err == nil {
						propertyValue = parsedUrl.String()
					}

				}
			case atom.A, atom.Area, atom.Link:
				if urlValue, exists := getAttr("href", node); exists {
					if parsedUrl, err := self.base.Parse(urlValue); err == nil {
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
				h5.WalkNodes(node, func(n *html.Node) {
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

	}

	for child := node.FirstChild; child != nil; {
		self.readItem(item, child)
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
