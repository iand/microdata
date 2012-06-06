package microdata

import (
	"strings"
	"testing"
)


func ReadData(html string, t *testing.T) *Microdata {
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

func ReadOneItem(html string, t *testing.T) *Item {
	data := ReadData(html, t)
	return data.items[0]
}

func TestRead(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Elizabeth</span>.</p>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["name"][0].(string) != "Elizabeth" {
		t.Errorf("Property value not found")
	}

}


func TestReadActuallyParses(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Daniel</span>.</p>
	</div>`
	item := ReadOneItem(html, t)

	if item.properties["name"][0].(string) != "Daniel" {
		t.Errorf("Property value not found")
	}

}


func TestReadThreeProps(t *testing.T) {
	html := `
	<div itemscope>
	 <p>My name is <span itemprop="name">Neil</span>.</p>
	 <p>My band is called <span itemprop="band">Four Parts Water</span>.</p>
	 <p>I am <span itemprop="nationality">British</span>.</p>
	</div>`

	item := ReadOneItem(html, t)

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


func TestReadImgSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <img itemprop="image" src="google-logo.png" alt="Google">
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["image"][0].(string) != "google-logo.png" {
		t.Errorf("Property value not found")
	}
}

func TestReadAHref(t *testing.T) {
	html := `
	<div itemscope>
	 <a itemprop="image" href="google-logo.png">foo</a>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["image"][0].(string) != "google-logo.png" {
		t.Errorf("Property value not found")
	}
}

func TestReadAreaHref(t *testing.T) {
	html := `
	<div itemscope><map name="shapes">
	 <area itemprop="foo" href="target.html" shape=rect coords="50,50,100,100">
	
	</map></div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target.html" {
		t.Errorf("Property value not found")
	}
}

func TestReadLinkHref(t *testing.T) {
	html := `
	<div itemscope>
		<link itemprop="foo" rel="author" href="target.html">
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target.html" {
		t.Errorf("Property value not found")
	}
}

func TestReadAudioSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <audio itemprop="foo" src="target"></audio>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestReadSourceSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <source itemprop="foo" src="target"></source>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}


func TestReadVideoSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <video itemprop="foo" src="target"></video>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestReadEmbedSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <embed itemprop="foo" src="target"></embed>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestReadTrackSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <track itemprop="foo" src="target"></track>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestReadIFrameSrc(t *testing.T) {
	html := `
	<div itemscope>
	 <iframe itemprop="foo" src="target"></iframe>
	</div>`

	item := ReadOneItem(html, t)

	if item.properties["foo"][0].(string) != "target" {
		t.Errorf("Property value not found")
	}
}

func TestReadDataValue(t *testing.T) {
	html := `
	<h1 itemscope>
 		<data itemprop="product-id" value="9678AOU879">The Instigator 2000</data>
	</h1>`

	item := ReadOneItem(html, t)

	if item.properties["product-id"][0].(string) != "9678AOU879" {
		t.Errorf("Property value not found")
	}
}

func TestReadTimeDatetime(t *testing.T) {
	html := `
	<h1 itemscope>
 		I was born on <time itemprop="birthday" datetime="2009-05-10">May 10th 2009</time>.
	</h1>`

	item := ReadOneItem(html, t)

	if item.properties["birthday"][0].(string) != "2009-05-10" {
		t.Errorf("Property value not found")
	}
}



func TestReadTwoValues(t *testing.T) {
	html := `
	<div itemscope>
	 <p>Flavors in my favorite ice cream:</p>
	 <ul>
	  <li itemprop="flavor">Lemon sorbet</li>
	  <li itemprop="flavor">Apricot sorbet</li>
	 </ul>
	</div>`

	item := ReadOneItem(html, t)
	if len(item.properties["flavor"]) != 2 {
		t.Errorf("Expecting 2 items but got %d",len(item.properties["flavor"]) )
	}
	if item.properties["flavor"][0].(string) != "Lemon sorbet" {
		t.Errorf("Property value 'Lemon sorbet' not found")
	}
	if item.properties["flavor"][1].(string) != "Apricot sorbet" {
		t.Errorf("Property value 'Apricot sorbet' not found")
	}


}
