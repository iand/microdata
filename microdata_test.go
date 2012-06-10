package microdata

import (
	"strings"
	"testing"
)


func ParseData(html string, t *testing.T) *Microdata {
	p := NewParser(strings.NewReader(html))

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
	return data.items[0]
}

func TestParse(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["name"][0].(string) != "Elizabeth" {
		t.Errorf("Property value not found")
	}

}


func TestParseActuallyParses(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Daniel</span>.</p>
	</div>`
	item := ParseOneItem(html, t)

	if item.properties["name"][0].(string) != "Daniel" {
		t.Errorf("Property value not found")
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

	if item.properties["name"][0].(string) != "Neil" {
		t.Errorf("Property value not found")
	}

	if item.properties["band"][0].(string) != "Four Parts Water" {
		t.Errorf("Property value not found")
	}

	if item.properties["nationality"][0].(string) != "British" {
		t.Errorf("Property value not found")
	}
}


func TestParseImgSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <img itemprop="image" src="google-logo.png" alt="Google">
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["image"][0].(string) != "google-logo.png" {
		t.Errorf("Property value not found")
	}
}

func TestParseAHref(t *testing.T) {
	html := `
	<div itemscope>
	 <a itemprop="image" href="google-logo.png">foo</a>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["image"][0].(string) != "google-logo.png" {
		t.Errorf("Property value not found")
	}
}

func TestParseAreaHref(t *testing.T) {
	html := `
	<div itemscope><map name="shapes">
	 <area itemprop="foo" href="target.html" shape=rect coords="50,50,100,100">
	
	</map></div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target.html" {
		t.Errorf("Property value not found")
	}
}

func TestParseLinkHref(t *testing.T) {
	html := `
	<div itemscope>
		<link itemprop="foo" rel="author" href="target.html">
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target.html" {
		t.Errorf("Property value not found")
	}
}

func TestParseAudioSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <audio itemprop="foo" src="target"></audio>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestParseSourceSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <source itemprop="foo" src="target"></source>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}


func TestParseVideoSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <video itemprop="foo" src="target"></video>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestParseEmbedSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <embed itemprop="foo" src="target"></embed>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestParseTrackSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <track itemprop="foo" src="target"></track>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestParseIFrameSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <iframe itemprop="foo" src="target"></iframe>
	</div>`

	item := ParseOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestParseDataValue(t *testing.T) {
	html := `
	<h1 itemscope>
 		<data itemprop="product-id" value="9678AOU879">The Instigator 2000</data>
	</h1>`

	item := ParseOneItem(html, t)

	if item.properties["product-id"][0].(string) != "9678AOU879" {
		t.Errorf("Property value not found")
	}
}

func TestParseTimeDatetime(t *testing.T) {
	html := `
	<h1 itemscope>
 		I was born on <time itemprop="birthday" datetime="2009-05-10">May 10th 2009</time>.
	</h1>`

	item := ParseOneItem(html, t)

	if item.properties["birthday"][0].(string) != "2009-05-10" {
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
	if len(item.properties["flavor"]) != 2 {
		t.Errorf("Expecting 2 values but got %d",len(item.properties["flavor"]) )
	}
	if item.properties["flavor"][0].(string) != "Lemon sorbet" {
		t.Errorf("Property value 'Lemon sorbet' not found")
	}
	if item.properties["flavor"][1].(string) != "Apricot sorbet" {
		t.Errorf("Property value 'Apricot sorbet' not found")
	}


}

func TestParseTwoPropertiesOneValue(t *testing.T) {
	html := `
	<div itemscope>
	 <span itemprop="favorite-color favorite-fruit">orange</span>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.properties) != 2 {
		t.Errorf("Expecting 2 properties but got %d",len(item.properties) )
	}
	if len(item.properties["favorite-color"]) != 1 {
		t.Errorf("Expecting 1 value but got %d",len(item.properties["favorite-color"]) )
	}
	if len(item.properties["favorite-fruit"]) != 1 {
		t.Errorf("Expecting 1 value but got %d",len(item.properties["favorite-fruit"]) )
	}
	if item.properties["favorite-color"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-color'")
	}
	if item.properties["favorite-fruit"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-fruit'")
	}
}

func TestParseTwoPropertiesOneValueMultispaced(t *testing.T) {
	html := `
	<div itemscope>
	 <span itemprop="   favorite-color    favorite-fruit   ">orange</span>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.properties) != 2 {
		t.Errorf("Expecting 2 properties but got %d",len(item.properties) )
	}

	if len(item.properties["favorite-color"]) != 1 {
		t.Errorf("Expecting 1 value but got %d",len(item.properties["favorite-color"]) )
	}
	if len(item.properties["favorite-fruit"]) != 1 {
		t.Errorf("Expecting 1 value but got %d",len(item.properties["favorite-fruit"]) )
	}
	if item.properties["favorite-color"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-color'")
	}
	if item.properties["favorite-fruit"][0].(string) != "orange" {
		t.Errorf("Property value 'orange' not found for 'favorite-fruit'")
	}
}

func TestParseItemType(t *testing.T) {
	html := `
	<div itemscope itemtype="http://example.org/animals#cat">
 		<h1 itemprop="name">Hedral</h1>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.types) != 1 {
		t.Errorf("Expecting 1 type but got %d",len(item.types) )	
	}

	if item.types[0] != "http://example.org/animals#cat" {
		t.Errorf("Expecting type of 'http://example.org/animals#cat' but got %d",item.types[0]) 
	}
}

func TestParseMultipleItemTypes(t *testing.T) {
	html := `
	<div itemscope itemtype=" http://example.org/animals#mammal  http://example.org/animals#cat  ">
 		<h1 itemprop="name">Hedral</h1>
	</div>`

	item := ParseOneItem(html, t)
	if len(item.types) != 2 {
		t.Errorf("Expecting 2 types but got %d",len(item.types) )	
	}

	if item.types[0] != "http://example.org/animals#mammal" {
		t.Errorf("Expecting type of 'http://example.org/animals#mammal' but got %d",item.types[0]) 
	}
	if item.types[1] != "http://example.org/animals#cat" {
		t.Errorf("Expecting type of 'http://example.org/animals#cat' but got %d",item.types[1]) 
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

	if item.id != "urn:isbn:0-330-34032-8" {
		t.Errorf("Expecting id of 'urn:isbn:0-330-34032-8' but got %d",item.id) 
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


	if len(item.properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d",len(item.properties) )
	}

	if item.properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
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

	if len(data.items) != 2 {
		t.Errorf("Expecting 2 items but got %d",len(data.items) )
	}
	if len(data.items[0].properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d",len(data.items[0].properties) )
	}
	if len(data.items[1].properties) != 3 {
		t.Errorf("Expecting 3 properties but got %d",len(data.items[1].properties) )
	}

	if data.items[0].properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
		t.Errorf("Property value 'http://www.opensource.org/licenses/mit-license.php' not found for 'license'")
	}

	if data.items[1].properties["license"][0].(string) != "http://www.opensource.org/licenses/mit-license.php" {
		t.Errorf("Property value 'http://www.opensource.org/licenses/mit-license.php' not found for 'license'")
	}

}