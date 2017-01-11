package model

import (
	"encoding/xml"
	"errors"
)

/*
Commodity = element gnc:commodity {
  attribute version { "2.0.0" },
  ( ( element cmdty:space { "ISO4217" },    # catégorie (monnaies)
      element cmdty:id { text }    # dénomination
    )
  | ( element cmdty:space { text },
      element cmdty:id { text },
      element cmdty:name { text }?,
      element cmdty:xcode { text }?,
      element cmdty:fraction { text }
    )
  ),
  ( element cmdty:get_quotes { empty },
    element cmdty:quote_source { text }?,
    element cmdty:quote_tz { text | empty }?
  )?,
  element cmdty:slots { KvpSlot+ }?
}

*/

type Commodities []*Commodity

type Commodity struct {
	Space       string `xml:"space"`
	ID          string `xml:"id"`
	Name        string `xml:"name"`
	Xcode       string `xml:"xcode"`
	Fraction    string `xml:"fraction"`
	GetQuote    string `xml:"get_quote"`
	QuoteSource string `xml:"quote_source"`
	QuoteTz     string `xml:"quote_tz"`
}

func (c *Commodity) String() string {
	return c.Space + ":" + c.ID
}

func CommodityUnmarshalXML(decoder *xml.Decoder, start xml.StartElement) (*Commodity, error) {
	// http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields
	// http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/

	cmdty := Commodity{}
	var foundID, foundSpace bool

LOOP:
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break LOOP
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			var ps *string

			switch se.Name.Local {
			case "space":
				ps = &cmdty.Space
				foundSpace = true
			case "id":
				ps = &cmdty.ID
				foundID = true
			case "name":
				ps = &cmdty.Name
			case "xcode":
				ps = &cmdty.Xcode
			case "fraction":
				ps = &cmdty.Fraction
			case "get_quote":
				ps = &cmdty.GetQuote
			case "quote_source":
				ps = &cmdty.QuoteSource
			case "quote_tz":
				ps = &cmdty.QuoteTz
			}

			if ps != nil {
				decoder.DecodeElement(ps, &se)
			}
		case xml.EndElement:
			if se.Name.Local == "commodity" {
				break LOOP
			}
		}
	}
	if !(foundID && foundSpace) {
		return nil, errors.New("Commodity without ID and Space values")
	}

	return &cmdty, nil
}

// Get returns the Commodity identified by space and id
func (cs Commodities) Get(space, id string) *Commodity {

	for _, c := range cs {
		if (c.Space == space) && (c.ID == id) {
			return c
		}
	}
	return nil
}

// Add adds a commodity to the collection
func (cs *Commodities) Add(c *Commodity) {
	if c.Space == "template" && c.ID == "template" {
		return
	}
	*cs = append(*cs, c)
}

// Len returns the number of commodities.
func (cs Commodities) Len() int {
	return len(cs)
}
