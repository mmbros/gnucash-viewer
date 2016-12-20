package model

import (
	"encoding/xml"
	"sort"
	"time"

	"github.com/mmbros/gnucash-viewer/types"
)

/*

Transaction = element gnc:transaction {
  attribute version { "2.0.0" },
  element trn:id { attribute type { "guid" }, GUID },
  element trn:currency {
    element cmdty:space { text },
    element cmdty:id { text }
  },
  element trn:num { text }?,
  element trn:date-posted { TimeSpec },
  element trn:date-entered { TimeSpec },
  element trn:description { text }?,
  element trn:slots { KvpSlot+ }?,
  element trn:splits { Split+ }
}

*/

// Transactions type is a list of Transaction
type Transactions []*Transaction

// Transaction type
type Transaction struct {
	Currency    *Commodity     `xml:"currency"`
	DatePosted  types.Timespec `xml:"date-posted>date"`
	DateEntered types.Timespec `xml:"date-entered>date"`
	Description string         `xml:"description"`
	Splits      Splits         `xml:"splits>split"`
}

func transactionUnmarshalXML(decoder *xml.Decoder, commodities Commodities, accMap AccountMap) (*Transaction, error) {

	var trn Transaction
	var cmdty Commodity

LOOP:
	for {
		// Read tokens from the XML document in a stream.
		token, _ := decoder.Token()
		if token == nil {
			break LOOP
		}
		// Inspect the type of the token just read.
		switch se := token.(type) {
		case xml.StartElement:
			var v interface{}

			switch se.Name.Local {
			case "currency":
				v = &cmdty
			case "description":
				v = &trn.Description
			case "date-posted":
				v = &trn.DatePosted
			case "date-entered":
				v = &trn.DateEntered
			case "splits":
				if err := SplitsUnmarshalXML(decoder, &trn.Splits, accMap); err != nil {
					return nil, err
				}

			}

			if v != nil {
				if err := decoder.DecodeElement(v, &se); err != nil {
					return nil, err
				}
				switch se.Name.Local {
				case "currency":
					trn.Currency = commodities.Get(cmdty.Space, cmdty.ID)
				}
			}

		case xml.EndElement:
			if se.Name.Local == "transaction" {
				break LOOP
			}
		}
	}

	return &trn, nil
}

// Add adds a transaction to the collection
func (ts *Transactions) Add(t *Transaction) {
	*ts = append(*ts, t)
}

// Len returns the number of transactions.
func (ts Transactions) Len() int {
	return len(ts)
}

// used to sort Transactions
type byDatePosted Transactions

func (t byDatePosted) Len() int      { return len(t) }
func (t byDatePosted) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t byDatePosted) Less(i, j int) bool {
	return time.Time(t[i].DatePosted).Before(time.Time(t[j].DatePosted))
}

// Sort sorts Transactions by DatePosted.
func (ts Transactions) Sort() {
	sort.Sort(byDatePosted(ts))
}
