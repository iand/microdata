package microdata

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

type ValueList []interface{}
type PropertyMap map[string]ValueList

type Item struct {
	Properties PropertyMap `json:"properties"`
	Types      []string    `json:"type,omitempty"`
	ID         string      `json:"id,omitempty"`
}

func NewItem() *Item {
	return &Item{
		Properties: make(PropertyMap, 0),
		Types:      make([]string, 0),
	}
}

func (self *Item) SetString(property string, value string) {
	self.Properties[property] = append(self.Properties[property], value)
}

func (self *Item) SetItem(property string, value *Item) {
	self.Properties[property] = append(self.Properties[property], value)
}

func (self *Item) AddType(value string) {
	self.Types = append(self.Types, value)
}

type Microdata struct {
	Items []*Item `json:"items"`
}

func NewMicrodata() *Microdata {
	return &Microdata{
		Items: make([]*Item, 0),
	}
}

func (self *Microdata) AddItem(value *Item) {
	self.Items = append(self.Items, value)
}

func (self *Microdata) Json() ([]byte, error) {
	b, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type Parser struct {
	p               *h5.Parser
	data            *Microdata
	base            *url.URL
	identifiedNodes map[string]*h5.Node
}

func NewParser(r io.Reader, base *url.URL) *Parser {
	return &Parser{
		p:    h5.NewParser(r),
		data: NewMicrodata(),
		base: base,
	}
}

func (self *Parser) Parse() (*Microdata, error) {
	err := self.p.Parse()
	if err != nil {
		return nil, err
	}
	tree := self.p.Tree()

	topLevelItemNodes := make([]*h5.Node, 0)
	self.identifiedNodes = make(map[string]*h5.Node, 0)

	tree.Walk(func(n *h5.Node) {
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

		if len(node.Children) > 0 {
			for _, child := range node.Children {
				self.readItem(item, child)
			}
		}
	}

	return self.data, nil
}

func (self *Parser) readItem(item *Item, node *h5.Node) {
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

			if len(node.Children) > 0 {
				for _, child := range node.Children {
					self.readItem(subitem, child)
				}
			}

			for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
				propertyName = strings.TrimSpace(propertyName)
				if propertyName != "" {
					item.SetItem(propertyName, subitem)
				}
			}

			return

		} else {
			var propertyValue string

			switch node.Data() {

			case "img", "audio", "source", "video", "embed", "iframe", "track":
				if urlValue, exists := getAttr("src", node); exists {
					if parsedUrl, err := self.base.Parse(urlValue); err == nil {
						propertyValue = parsedUrl.String()
					}

				}
			case "a", "area", "link":
				if urlValue, exists := getAttr("href", node); exists {
					if parsedUrl, err := self.base.Parse(urlValue); err == nil {
						propertyValue = parsedUrl.String()
					}
				}
			case "data":
				if urlValue, exists := getAttr("value", node); exists {
					propertyValue = urlValue
				}
			case "time":
				if urlValue, exists := getAttr("datetime", node); exists {
					propertyValue = urlValue
				}

			default:
				var text bytes.Buffer
				node.Walk(func(n *h5.Node) {
					if n.Type == h5.TextNode {
						text.WriteString(n.Data())
					}

				})
				propertyValue = text.String()
			}

			if len(propertyValue) > 0 {
				for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
					propertyName = strings.TrimSpace(propertyName)
					if propertyName != "" {
						item.SetString(propertyName, propertyValue)
					}
				}
			}

		}

	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			self.readItem(item, child)
		}
	}

}

func getAttr(name string, node *h5.Node) (string, bool) {
	for _, a := range node.Attr {
		if a.Name == name {
			return a.Value, true
		}
	}
	return "", false
}
