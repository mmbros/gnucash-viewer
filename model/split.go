package model

import (
	"encoding/xml"
	"fmt"

	"github.com/mmbros/gnucash-viewer/types"
)

// Splits type
type Splits []*Split

// Split type
type Split struct {
	ReconciledState types.ReconciledState `xml:"reconciled-state"`
	ReconcileDate   types.Timespec        `xml:"reconcile-date"`
	Value           types.Numeric         `xml:"value"`
	Memo            string                `xml:"memo"`
	Quantity        types.Numeric         `xml:"quantity"`
	Account         *Account              `xml:"-"`
}

// Add adds a split to the collection
func (ss *Splits) Add(s *Split) {
	*ss = append(*ss, s)
}

// Len returns the number of splits.
func (ss Splits) Len() int {
	return len(ss)
}

func SplitsUnmarshalXML(decoder *xml.Decoder, splits *Splits, accMap AccountMap) error {

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
			if se.Name.Local == "split" {
				s, err := SplitUnmarshalXML(decoder, accMap)
				if err != nil {
					return err
				}
				splits.Add(s)
			}
		case xml.EndElement:
			if se.Name.Local == "splits" {
				break LOOP
			}
		}
	}

	return nil
}

func SplitUnmarshalXML(decoder *xml.Decoder, accMap AccountMap) (*Split, error) {

	var (
		split     Split
		accountID string
	)

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
			case "reconciled-state":
				v = &split.ReconciledState
			case "reconcile-date":
				v = &split.ReconcileDate
			case "value":
				v = &split.Value
			case "memo":
				v = &split.Memo
			case "quantity":
				v = &split.Quantity
			case "account":
				v = &accountID
			}

			if v != nil {
				if err := decoder.DecodeElement(v, &se); err != nil {
					return nil, err
				}
				switch se.Name.Local {
				case "account":
					acc, ok := accMap[accountID]
					if !ok {
						return nil, fmt.Errorf("Account not found: %s", accountID)
					}
					split.Account = acc
				}
			}

		case xml.EndElement:
			if se.Name.Local == "split" {
				break LOOP
			}
		}
	}

	return &split, nil
}
