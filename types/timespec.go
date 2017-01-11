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

/*

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// This comparison is different from using t == u, which also compares
// the locations.
func (ts Timespec) Equal(u time.Time) bool {
	return time.Time(ts).Equal(u)
}

// After reports whether the time instant t is after u.
func (ts Timespec) After(u time.Time) bool {
	return time.Time(ts).After(u)
}

// Before reports whether the time instant t is before u.
func (ts Timespec) Before(u time.Time) bool {
	return time.Time(ts).Before(u)
}

*/
