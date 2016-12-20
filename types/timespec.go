package types

import (
	"encoding/xml"
	"time"
)

/*
TimeSpec = ( TimeStamp,
             element ts:ns { xsd:int }?
           )

TimeStamp = element ts:date { xsd:string { pattern = "[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} (\+|-)[0-9]{4}" } }
*/

// Timespec represent a gnucash Timespec value.
type Timespec time.Time

// UnmarshalXML implements xml.Unmarshaler interface
func (ts *Timespec) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields
	const gncForm = "2006-01-02 15:04:05 -0700"
	//var v string
	var v struct {
		Date string `xml:"date"`
		Ns   int    `xml:"ns"`
	}
	d.DecodeElement(&v, &start)

	x, err := time.Parse(gncForm, v.Date)
	if err != nil {
		return err
	}
	*ts = Timespec(x)

	return nil
}

// String returns the time formatted using the format string
//	"2006-01-02 15:04:05.999999999 -0700 MST"
// Return "" in case of zero Timespec.
func (ts Timespec) String() string {
	t := time.Time(ts)
	if t.IsZero() {
		return ""
	}
	return t.String()
}
