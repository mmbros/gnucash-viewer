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
)

var gnucashPath = flag.String("gnucash-file", "data-crypt/mau.gnucash", "GnuCash file path")

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
		fmt.Printf("%3d) %s - splits=%d\n", j+1, t.Description, t.Splits.Len())
	}
	//book.Accounts.PrintTree("   ")

}
