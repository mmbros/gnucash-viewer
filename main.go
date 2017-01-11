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
)

var gnucashPath = flag.String("gnucash-file", "data-crypt/mau.gnucash", "GnuCash file path")

func testQuery(b *model.Book) {

	loc, _ := time.LoadLocation("Local")
	afterEq := time.Date(2016, 1, 2, 0, 0, 0, 0, loc)
	before := afterEq.Add(24 * time.Hour)

	filters := query.NewFilters()
	filters.DatePostedRange(afterEq, before)

	fmt.Println(filters.String())

	results := query.Where(b, filters)
	fmt.Printf("results LEN = %+v\n", len(results))
	for j, r := range results {
		fmt.Printf("%02d) %+v\n", j+1, r.T)
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
	accounts, err := query.FindAccounts(".//Benzina", book.Accounts.Root)
	for j, a := range accounts {
		fmt.Printf("%d) %s\n", j+1, a.FullName())
	}

	testQuery(book)

}
