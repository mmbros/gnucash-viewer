package main

/*
	References:
	* [Parsing huge xml files with go](http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/)
	* [xml:nested namespace](https://groups.google.com/forum/#!topic/golang-nuts/QFDHM7_VFks)
	* [xml unmarshall example](https://golang.org/pkg/encoding/xml/#example_Unmarshal)

	* [Golang XML Unmarshal and time.Time fields](http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields)
*/

import (
	"flag"
	"fmt"
	"time"

	"github.com/mmbros/gnucash-viewer/model"
	"github.com/mmbros/gnucash-viewer/query"
	"github.com/mmbros/gnucash-viewer/types"
)

var gnucashPath = flag.String("gnucash-file", "data-crypt/mau.gnucash", "GnuCash file path")

func InfoTransaction(t *model.Transaction) {

	if t.Splits.Len() > 2 {
		fmt.Printf("*** SPLITS = %d *** \n", t.Splits.Len())
		return
	}
	sPos, sNeg := t.Splits[0], t.Splits[1]
	if sPos.Value.Sign() < 0 {
		sPos, sNeg = sNeg, sPos
	}

	atPos := sPos.Account.BasicType()
	atNeg := sNeg.Account.BasicType()

	// Asset  -> Expense : uscite
	// Income -> Asset   : entrate
	// Asset  -> Asset   : trasferimenti
	// else              : boh

	descr := "OTHER"
	switch atNeg {
	case types.AccountTypeAsset:
		switch atPos {
		case types.AccountTypeExpense:
			descr = "USCITE"
		case types.AccountTypeAsset:
			descr = "MOVIMENTO"
		}
	case types.AccountTypeIncome:
		if atPos == types.AccountTypeAsset {
			descr = "ENTRATE"
		}
	}

	fmt.Printf("  %0.2f %s : %s --> %s\n",
		sPos.Value.Float64(),
		descr,
		atNeg,
		atPos,
	)
}

func PrintTransaction(idx int, t *model.Transaction) {

	fmt.Printf("%02d) %s %s\n",
		idx,
		t.DatePosted.YMD(),
		t.Description,
	)
	for k, s := range t.Splits {
		fmt.Printf("  [%d] %0.2f  %s [%s]\n",
			k,
			s.Value.Float64(),
			s.Account.Name,
			s.Account.Type,
		)
	}
	InfoTransaction(t)
}

func testQuery(b *model.Book) {

	loc, _ := time.LoadLocation("Local")

	afterEq := time.Date(2016, 1, 1, 0, 0, 0, 0, loc)
	before := time.Date(2017, 1, 1, 0, 0, 0, 0, loc)
	//before := afterEq.Add(24 * time.Hour)

	results := query.Query(b).
		DatePostedRange(afterEq, before).
		//AccountType(types.AccountTypeExpense).
		AccountPath(".//Benzina").
		Execute()

	fmt.Printf("results LEN = %+v\n", len(results))
	for j, r := range results {
		fmt.Printf("%02d) %s  %0.2f EUR  %s\n",
			j+1,
			r.T.DatePosted.YMD(),
			r.S.Value.Float64(),
			r.T.Description,
		)

	}

	afterEq = time.Date(2016, 8, 1, 0, 0, 0, 0, loc)
	before = time.Date(2016, 9, 1, 0, 0, 0, 0, loc)
	results = query.Query(b).
		DatePostedRange(afterEq, before).
		Execute()

	fmt.Printf("results LEN = %+v\n", len(results))
	for j, r := range results {

		/*
			var descr string

			if len(r.T.Splits) == 2 {
				// simple transaction
				s0, s1 := r.T.Splits[0], r.T.Splits[1]

				// Entrate:
				//   Entrate -> Attivita
				// Uscite:
				//   Attivita -> Uscite
				//   Passivita -> Uscite
				// Rimborsi:
				//   Uscite -> Entrate
				//   Uscite -> Passivita
				// Altro:
				//   Attivita -> Passivita

				descr = fmt.Sprintf("%s {%s} <--> %s {%s}", s0.Account.Name, s0.Account.Type.String(), s1.Account.Name, s1.Account.Type.String())

			} else {
				descr = fmt.Sprintf("*** %d splits ***", len(r.T.Splits))
			}
		*/

		fmt.Printf("%02d) %s %s\n",
			j+1,
			r.T.DatePosted.YMD(),
			r.T.Description,
		)
		for k, s := range r.T.Splits {
			fmt.Printf("  [%d] %0.2f  %s [%s]\n",
				k,
				s.Value.Float64(),
				s.Account.Name,
				s.Account.Type,
			)
		}

	}
}

func main() {
	defer timeTrack(time.Now(), "task duration:")

	gnc, err := model.ReadFile(*gnucashPath)

	if err != nil {
		panic(err)
	}

	book := gnc.Book

	fmt.Printf("Commodites   (1)   : %d\n", book.Commodities.Len())
	fmt.Printf("Accounts     (161) : %d\n", book.Accounts.Len())
	fmt.Printf("Transactions (2553): %d\n", book.Transactions.Len())
	fmt.Println("")

	for j, cmdty := range book.Commodities {
		fmt.Printf("*** COMMODITY #%d ***\n", j)
		fmt.Printf("ID: %s\n", cmdty.ID)
		fmt.Printf("Space: %s\n", cmdty.Space)
		fmt.Printf("Name: %s\n", cmdty.Name)
		fmt.Printf("Xcode: %s\n", cmdty.Xcode)
		fmt.Printf("Fraction: %s\n", cmdty.Fraction)
	}

	fmt.Println("")

	for j, a := range book.Accounts.List {
		if j >= 10 {
			break
		}
		fmt.Printf("%3d) %v %s\n", j+1, a, a.Currency.ID)
	}

	fmt.Println("")

	for j, t := range book.Transactions {
		if j >= 10 {
			break
		}
		fmt.Printf("%3d) %s - %s - splits=%d\n",
			j+1, t.DatePosted, t.Description, t.Splits.Len())
	}
	book.Accounts.PrintTree("   ")

	// check Transactions is sorted by DatePosted
	var precTime time.Time
	for j, t := range book.Transactions {
		currTime := time.Time(t.DatePosted)
		if currTime.Before(precTime) {
			fmt.Printf("Transactions(%d): currTime < precTime - %s < %s\n", j, currTime.UTC(), precTime.UTC())

			fmt.Println(j-1, book.Transactions[j-1])
			fmt.Println(j, t)
			return
		}
		precTime = currTime
	}

	// find account by path
	//accounts, err := query.FindAccounts(".//Benzina", book.Accounts.Root)
	//for j, a := range accounts {
	//fmt.Printf("%d) %s\n", j+1, a.FullName())
	//}

	//testQuery(book)

	found := 0

	for _, t := range book.Transactions {
		if t.Splits.Len() > 1 {
			found++
			PrintTransaction(found, t)

			if found >= 30 {
				break
			}
		}
	}
}
