package model

import (
	"compress/gzip"
	"encoding/xml"
	"os"
)

/*
GnuCashXml = element gnc-v2 {

	# from gnc_book_write_to_xml_filehandle_v2 in src/backend/xml/io-gncxml-v2.c

	  ( ( element gnc:count-data { attribute cd:type { "book"  }, "1"  },
        Book
	    )

		# from gnc_book_write_accounts_to_xml_filehandle_v2 in src/backend/xml/io-gncxml-v2.c

	  | ( element gnc:count-data { attribute cd:type { "commodity"  }, xsd:int  }?,
        element gnc:count-data { attribute cd:type { "account"  }, xsd:int  }?,
	      Commodity*,
        Account*
		    )
		  )

}




<gnc-v2>
	<gnc:book>
		<gnc:account>
		<gnc:transaction>
			<trn:slots>
				<slot>
      <trn:splits>
        <trn:split>

// Transactions type
type Transactions []*Transaction

// Transaction type
type Transaction struct {
	ID          string
	Currency    string
	DatePosted  time.Time
	DateEntered time.Time
	Description string
	Splits      []*Split
}

// Split type
type Split struct {
	ID              string
	ReconciledState string
	ReconcileDate   time.Time
	Value           *numeric.Numeric
	Memo            string
	Quantity        *numeric.Numeric
	Account         *Account
}
*/

// Gnc type
type Gnc struct {
	XMLName xml.Name `xml:"gnc-v2"`
	Book    *Book    `xml:"book"`
}

// Slot type : integer | string | frame | gdate | numeric
type Slot struct {
	XMLName  xml.Name `xml:"slot"`
	Key      string   `xml:"key"`
	Value    string   `xml:"value"`
	Children []*Slot
}

// ReadFile read the gnucash file in XML format
func ReadFile(path string) (*Gnc, error) {

	// open gnucash file
	gnucashFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer gnucashFile.Close()

	// decompress gnucash file
	reader, err := gzip.NewReader(gnucashFile)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// unmarshall xml
	gnc := Gnc{}
	err = xml.NewDecoder(reader).Decode(&gnc)
	if err != nil {
		return nil, err
	}

	return &gnc, nil
}
