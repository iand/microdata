/*
  This is free and unencumbered software released into the public domain. For more
  information, see <http://unlicense.org/> or the accompanying UNLICENSE file.
*/

package microdata

import (
	"bytes"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func ParseData(html string, t *testing.T) *Microdata {
	u, _ := url.Parse("http://example.com/")
	p := NewParser(strings.NewReader(html), u)

	data, err := p.Parse()
	if err != nil {
		t.Errorf("Expected no error but got %d", err)
	}

	if data == nil {
		t.Errorf("Expected non-nil data")
	}

	return data
}

func ParseOneItem(html string, t *testing.T) *Item {
	data := ParseData(html, t)
	return data.Items[0]
}

func TestParse(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["name"][0].(string) != "Elizabeth" {
		t.Errorf("Property value not found")
	}

}

func TestParseActuallyParses(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Daniel</span>.</p>
	</div>`
	item := ParseOneItem(html, t)

	if item.Properties["name"][0].(string) != "Daniel" {
		t.Errorf("got %v, wanted %s", item.Properties["name"][0], "Daniel")
	}

}

func TestParseThreeProps(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Neil</span>.</p>
	 <p>My band is called <span itemprop="band">Four Parts Water</span>.</p>
	 <p>I am <span itemprop="nationality">British</span>.</p>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["name"][0].(string) != "Neil" {
		t.Errorf("Property value not found")
	}

	if item.Properties["band"][0].(string) != "Four Parts Water" {
		t.Errorf("Property value not found")
	}

	if item.Properties["nationality"][0].(string) != "British" {
		t.Errorf("Property value not found")
	}
}

func TestParseImgSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <img itemprop="image" src="http://example.com/foo" alt="Google">
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["image"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseAHref(t *testing.T) {
	html := `
	<div itemscope>
	 <a itemprop="image" href="http://example.com/foo">foo</a>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["image"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseAreaHref(t *testing.T) {
	html := `
	<div itemscope><map name="shapes">
	 <area itemprop="foo" href="http://example.com/foo" shape=rect coords="50,50,100,100">

	</map></div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseLinkHref(t *testing.T) {
	html := `
	<div itemscope>
		<link itemprop="foo" rel="author" href="http://example.com/foo">
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseAudioSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <audio itemprop="foo" src="http://example.com/foo"></audio>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseSourceSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <source itemprop="foo" src="http://example.com/foo"></source>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseVideoSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <video itemprop="foo" src="http://example.com/foo"></video>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseEmbedSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <embed itemprop="foo" src="http://example.com/foo"></embed>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseTrackSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <track itemprop="foo" src="http://example.com/foo"></track>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseIFrameSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <iframe itemprop="foo" src="http://example.com/foo"></iframe>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["foo"][0].(string) != "http://example.com/foo" {
		t.Errorf("Property value not found")
	}
}

func TestParseDataValue(t *testing.T) {
	html := `
	<h1 itemscope>
 		<data itemprop="product-id" value="9678AOU879">The Instigator 2000</data>
	</h1>`

	item := ParseOneItem(html, t)

	if item.Properties["product-id"][0].(string) != "9678AOU879" {
		t.Errorf("Property value not found")
	}
}

func TestParseTimeDatetime(t *testing.T) {
	html := `
	<h1 itemscope>
 		I was born on <time itemprop="birthday" datetime="2009-05-10">May 10th 2009</time>.
	</h1>`

	item := ParseOneItem(html, t)

	if item.Properties["birthday"][0].(string) != "2009-05-10" {
		t.Errorf("Property value not found")
	}
}

func TestParseTwoValues(t *testing.T) {
	html := `
	<div itemscope>
	 <p>Flavors in my favorite ice cream:</p>
	 <ul>
	  <li itemprop="flavor">Lemon sorbet</li>
	  <li itemprop="flavor">Apricot sorbet</li>
	 </ul>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.Properties["flavor"]) != 2 {
		t.Errorf("Expecting 2 values but got %d", len(item.Properties["flavor"]))
	}
	if item.Properties["flavor"][0].(string) != "Lemon sorbet" {
		t.Errorf("Property value 'Lemon sorbet' not found")
	}
	if item.Properties["flavor"][1].(string) != "Apricot sorbet" {
		t.Errorf("Property value 'Apricot sorbet' not found")
	}

}

func TestParseTwoPropertiesOneValue(t *testing.T) {
	html := `
	<div itemscope>
	 <span itemprop="favorite-color favorite-fruit">orange</span>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.Properties) != 2 {
		t.Errorf("Expecting 2 properties but got %d", len(item.Properties))
	}
	if len(item.Properties["favorite-color"]) != 1 {
		t.Errorf("Expecting 1 value but got %d", len(item.Properties["favorite-color"]))
	}
	if len(item.Properties["favorite-fruit"]) != 1 {
		t.Errorf("Expecting 1 value but got %d", len(item.Properties["favorite-fruit"]))
	}
	if item.Properties["favorite-color"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-color'")
	}
	if item.Properties["favorite-fruit"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-fruit'")
	}
}

func TestParseTwoPropertiesOneValueMultispaced(t *testing.T) {
	html := `
	<div itemscope>
	 <span itemprop="   favorite-color    favorite-fruit   ">orange</span>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.Properties) != 2 {
		t.Errorf("Expecting 2 properties but got %d", len(item.Properties))
	}

	if len(item.Properties["favorite-color"]) != 1 {
		t.Errorf("Expecting 1 value but got %d", len(item.Properties["favorite-color"]))
	}
	if len(item.Properties["favorite-fruit"]) != 1 {
		t.Errorf("Expecting 1 value but got %d", len(item.Properties["favorite-fruit"]))
	}
	if item.Properties["favorite-color"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-color'")
	}
	if item.Properties["favorite-fruit"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-fruit'")
	}
}

func TestParseItemType(t *testing.T) {
	html := `
	<div itemscope itemtype="http://example.org/animals#cat">
 		<h1 itemprop="name">Hedral</h1>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.Types) != 1 {
		t.Errorf("Expecting 1 type but got %d", len(item.Types))
	}

	if item.Types[0] != "http://example.org/animals#cat" {
		t.Errorf("Expecting type of 'http://example.org/animals#cat' but got %s", item.Types[0])
	}
}

func TestParseMultipleItemTypes(t *testing.T) {
	html := `
	<div itemscope itemtype=" http://example.org/animals#mammal  http://example.org/animals#cat  ">
 		<h1 itemprop="name">Hedral</h1>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.Types) != 2 {
		t.Errorf("Expecting 2 types but got %d", len(item.Types))
	}

	if item.Types[0] != "http://example.org/animals#mammal" {
		t.Errorf("Expecting type of 'http://example.org/animals#mammal' but got %s", item.Types[0])
	}
	if item.Types[1] != "http://example.org/animals#cat" {
		t.Errorf("Expecting type of 'http://example.org/animals#cat' but got %s", item.Types[1])
	}
}

func TestParseItemId(t *testing.T) {
	html := `<dl itemscope
	    itemtype="http://vocab.example.net/book"
	    itemid="urn:isbn:0-330-34032-8">
	 <dt>Title
	 <dd itemprop="title">The Reality Dysfunction
	 <dt>Author
	 <dd itemprop="author">Peter F. Hamilton
	 <dt>Publication date
	 <dd><time itemprop="pubdate" datetime="1996-01-26">26 January 1996</time>
	</dl>`

	item := ParseOneItem(html, t)

	if item.ID != "urn:isbn:0-330-34032-8" {
		t.Errorf("Expecting id of 'urn:isbn:0-330-34032-8' but got %s", item.ID)
	}
}

func TestParseItemRef(t *testing.T) {
	html := `<body><p><figure itemscope itemtype="http://n.whatwg.org/work" itemref="licenses">
   <img itemprop="work" src="images/house.jpeg" alt="A white house, boarded up, sits in a forest.">
   <figcaption itemprop="title">The house I found.</figcaption>
  </figure></p>
   <p id="licenses">All images licensed under the <a itemprop="license"
   href="http://www.opensource.org/licenses/mit-license.php">MIT
   license</a>.</p></body>`

	item := ParseOneItem(html, t)

	if len(item.Properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d", len(item.Properties))
	}

	if item.Properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
		t.Errorf("Property value 'http://www.opensource.org/licenses/mit-license.php' not found for 'license'")
	}

}

func TestParseSharedItemRef(t *testing.T) {
	html := `<!DOCTYPE HTML>
		<html>
		 <head>
		  <title>Photo gallery</title>
		 </head>
		 <body>
		  <h1>My photos</h1>
		  <figure itemscope itemtype="http://n.whatwg.org/work" itemref="licenses">
		   <img itemprop="work" src="images/house.jpeg" alt="A white house, boarded up, sits in a forest.">
		   <figcaption itemprop="title">The house I found.</figcaption>
		  </figure>
		  <figure itemscope itemtype="http://n.whatwg.org/work" itemref="licenses">
		   <img itemprop="work" src="images/mailbox.jpeg" alt="Outside the house is a mailbox. It has a leaflet inside.">
		   <figcaption itemprop="title">The mailbox.</figcaption>
		  </figure>
		  <footer>
		   <p id="licenses">All images licensed under the <a itemprop="license"
		   href="http://www.opensource.org/licenses/mit-license.php">MIT
		   license</a>.</p>
		  </footer>
		 </body>
		</html>`

	data := ParseData(html, t)

	if len(data.Items) != 2 {
		t.Errorf("Expecting 2 items but got %d", len(data.Items))
	}
	if len(data.Items[0].Properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d", len(data.Items[0].Properties))
	}
	if len(data.Items[1].Properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d", len(data.Items[1].Properties))
	}

	if data.Items[0].Properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
		t.Errorf("Property value 'http://www.opensource.org/licenses/mit-license.php' not found for 'license'")
	}

	if data.Items[1].Properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
		t.Errorf("Property value 'http://www.opensource.org/licenses/mit-license.php' not found for 'license'")
	}

}

func TestParseMultiValuedItemRef(t *testing.T) {
	html := `<!DOCTYPE HTML>
		<html>
		 <body>
		 	<div itemscope id="amanda" itemref="a b"></div>
			<p id="a">Name: <span itemprop="name">Amanda</span></p>
			<p id="b">Age: <span itemprop="age">26</span></p>

		 </body>
		</html>`

	data := ParseData(html, t)

	if data.Items[0].Properties["name"][0].(string) != "Amanda" {
		t.Errorf("Property value 'Amanda' not found for 'name'")
	}

	if data.Items[0].Properties["age"][0].(string) != "26" {
		t.Errorf("Property value '26' not found for 'age'")
	}
}

func TestParseEmbeddedItem(t *testing.T) {
	html := `<div itemscope>
			 <p>Name: <span itemprop="name">Amanda</span></p>
			 <p>Band: <span itemprop="band" itemscope> <span itemprop="name">Jazz Band</span> (<span itemprop="size">12</span> players)</span></p>
			</div>`

	data := ParseData(html, t)

	if len(data.Items) != 1 {
		t.Errorf("Expecting 1 item but got %d", len(data.Items))
	}

	if data.Items[0].Properties["name"][0].(string) != "Amanda" {
		t.Errorf("Property value 'Amanda' not found for 'name'")
	}

	subitem := data.Items[0].Properties["band"][0].(*Item)

	if subitem.Properties["name"][0].(string) != "Jazz Band" {
		t.Errorf("Property value 'Jazz Band' not found for 'name'")
	}
}

func TestParseEmbeddedItemWithItemRef(t *testing.T) {
	html := `<body>
			<div itemscope id="amanda" itemref="a b"></div>
		<p id="a">Name: <span itemprop="name">Amanda</span></p>
		<div id="b" itemprop="band" itemscope itemref="c"></div>
		<div id="c">
		 <p>Band: <span itemprop="name">Jazz Band</span></p>
		 <p>Size: <span itemprop="size">12</span> players</p>
		</div></body>`

	data := ParseData(html, t)

	if len(data.Items) != 1 {
		t.Errorf("Expecting 1 item but got %d", len(data.Items))
	}

	if data.Items[0].Properties["name"][0].(string) != "Amanda" {
		t.Errorf("Property value 'Amanda' not found for 'name'")
	}

	subitem := data.Items[0].Properties["band"][0].(*Item)

	if subitem.Properties["name"][0].(string) != "Jazz Band" {
		t.Errorf("Property value 'Jazz Band' not found for 'name'")
	}
}

func TestParseRelativeURL(t *testing.T) {
	html := `
	<div itemscope>
	 <a itemprop="image" href="test.png">foo</a>
	</div>`

	item := ParseOneItem(html, t)

	if item.Properties["image"][0].(string) != "http://example.com/test.png" {
		t.Errorf("Property value not found")
	}
}

func TestParseItemRelativeId(t *testing.T) {
	html := `<dl itemscope
	    itemtype="http://vocab.example.net/book"
	    itemid="foo">
	 <dt>Title
	 <dd itemprop="title">The Reality Dysfunction
	 <dt>Author
	 <dd itemprop="author">Peter F. Hamilton
	 <dt>Publication date
	 <dd><time itemprop="pubdate" datetime="1996-01-26">26 January 1996</time>
	</dl>`

	item := ParseOneItem(html, t)

	if item.ID != "http://example.com/foo" {
		t.Errorf("Expecting id of 'http://example.com/foo' but got %s", item.ID)
	}
}

func TestJSON(t *testing.T) {
	item := NewItem()
	item.AddString("name", "Elizabeth")

	data := NewMicrodata()
	data.AddItem(item)

	expected := []byte(`{"items":[{"properties":{"name":["Elizabeth"]}}]}`)

	actual, _ := data.JSON()

	if !bytes.Equal(actual, expected) {
		t.Errorf("Expecting %s but got %s", expected, actual)
	}
}

func TestJsonWithType(t *testing.T) {
	item := NewItem()
	item.AddType("http://example.org/animals#cat")
	item.AddString("name", "Elizabeth")

	data := NewMicrodata()
	data.AddItem(item)

	expected := []byte(`{"items":[{"properties":{"name":["Elizabeth"]},"type":["http://example.org/animals#cat"]}]}`)

	actual, _ := data.JSON()

	if !bytes.Equal(actual, expected) {
		t.Errorf("Expecting %s but got %s", expected, actual)
	}
}

// This test checks stack overflow doesn't happen as mentioned in
// https://github.com/iand/microdata/issues/3
func TestSkipSelfReferencingItemref(t *testing.T) {
	html := `<body itemscope itemtype="http://schema.org/WebPage">
	  <span id="1" itemscope itemtype="http://data-vocabulary.org/Breadcrumb" itemprop="child" itemref="1">
	    <a title="Foo" itemprop="url" href="/foo/bar"><span itemprop="title">Foo</span></a>
	  </span>
	</body>`

	actual := ParseData(html, t)

	child := NewItem()
	child.AddString("title", "Foo")
	child.AddString("url", "http://example.com/foo/bar")

	item := NewItem()
	item.AddType("http://schema.org/WebPage")
	item.AddItem("child", child)

	expected := NewMicrodata()
	expected.AddItem(item)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expecting %s but got %s", expected, actual)
	}
}
