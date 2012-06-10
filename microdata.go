package microdata

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"io"
	"strings"
)

type ValueList []interface{}
type PropertyMap map[string]ValueList

type Item struct {
	properties PropertyMap
	types      []string
	id         string
}

func NewItem() *Item {
	return &Item{
		properties: make(PropertyMap, 0),
		types:      make([]string, 0),
	}
}

func (self *Item) SetString(property string, value string) {
	self.properties[property] = append(self.properties[property], value)
}

func (self *Item) SetItem(property string, value *Item) {
	self.properties[property] = append(self.properties[property], value)
}


type Microdata struct {
	items []*Item
}

func NewMicrodata() *Microdata {
	return &Microdata{
		items: make([]*Item, 0),
	}
}

type Parser struct {
	p               *h5.Parser
	data            *Microdata
	identifiedNodes map[string]*h5.Node
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		p:    h5.NewParser(r),
		data: NewMicrodata(),
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
		self.data.items = append(self.data.items, item)
		if itemtypes, exists := getAttr("itemtype", node); exists {
			for _, itemtype := range strings.Split(strings.TrimSpace(itemtypes), " ") {
				itemtype = strings.TrimSpace(itemtype)
				if itemtype != "" {
					item.types = append(item.types, itemtype)
				}
			}
			// itemid only valid when itemscope and itemtype are both present
			if itemid, exists := getAttr("itemid", node); exists {
				item.id = strings.TrimSpace(itemid)
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
					propertyValue = urlValue
				}
			case "a", "area", "link":
				if urlValue, exists := getAttr("href", node); exists {
					propertyValue = urlValue
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

			for _, propertyName := range strings.Split(strings.TrimSpace(itemprop), " ") {
				propertyName = strings.TrimSpace(propertyName)
				if propertyName != "" {
					item.SetString(propertyName, propertyValue)
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
