package model

import (
	"encoding/xml"
	"errors"
	"fmt"
)

/*
Book = element gnc:book {
  attribute version { "2.0.0" },

# from write_book_parts in src/backend/xml/gnc-book-xml-v2.c

  element book:id { attribute type { "guid" }, GUID },
  element book:slots { KvpSlot+ }?,

# from write_book in src/backend/xml/io-gncxml-v2.c

  element gnc:count-data { attribute cd:type { "commodity" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "account" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "transaction" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "schedxaction" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "budget" }, xsd:int }?,

# plugins (those with a get_count slot)

  element gnc:count-data { attribute cd:type { "gnc:GncBillTerm" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncCustomer" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncEmployee" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncEntry" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncInvoice" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncJob" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncOrder" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncTaxTable" }, xsd:int }?,
  element gnc:count-data { attribute cd:type { "gnc:GncVendor" }, xsd:int }?,

  Commodity*,
  PriceDb?,
  Account*,
  Transaction*,
  TemplateTransactions*,
  ScheduledTransaction*,
  Budget*,

# plugins (those with a write slot)

  BillTerm*,
  Customer*,
  Employee*,
  Entry*,
  Invoice*,
  Job*,
  Order*,
  TaxTable*,
  Vendor*
}


*/

// Book type
type Book struct {
	Commodities Commodities
	//PriceDb?,
	Accounts     Accounts
	Transactions Transactions
}

// UnmarshalXML implements xml.Unmarshaler interface
func (b *Book) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields
	// http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/

	var (
		book   Book
		accMap AccountMap
	)
	accMap = AccountMap{}

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "commodity":
				var cmdty Commodity
				decoder.DecodeElement(&cmdty, &se)
				book.Commodities.Add(&cmdty)
			case "account":
				account, id, parentId, err := AccountUnmarshalXML(decoder, book.Commodities)
				if err != nil {
					return err
				}
				// add new account in book's accounts
				book.Accounts.Add(account)
				// update map from id -> account
				accMap[id] = account
				if parentId == "" {
					// root account
					if book.Accounts.Root != nil {
						return errors.New("Invalid multi root accounts")
					}
					book.Accounts.Root = account
				} else {
					// setup parent and children values
					parent, ok := accMap[parentId]
					if ok {
						account.Parent = parent
						parent.Children = append(parent.Children, account)
					} else {
						return fmt.Errorf("Parent account not (yes) found: parent=%s", parentId)
					}
				}
			case "transaction":
				trn, err := TransactionUnmarshalXML(decoder, book.Commodities, accMap)
				if err != nil {
					return err
				}
				book.Transactions.Add(trn)

			}

		case xml.EndElement:
			if se.Name.Local == "book" {
				break
			}
		}
	}
	*b = book

	return nil
}
