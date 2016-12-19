package types

import (
	"encoding/xml"
	"fmt"
)

// ReconciledState enum type
type ReconciledState int

// ReconciledState constants
const (
	ReconciledStateY ReconciledState = iota
	ReconciledStateN
	ReconciledStateC
	ReconciledStateF
	ReconciledStateV
)

func ReconciledStateFromString(v string) (ReconciledState, error) {

	var s2i = map[string]ReconciledState{
		"y": ReconciledStateY,
		"n": ReconciledStateN,
		"c": ReconciledStateC,
		"f": ReconciledStateF,
		"v": ReconciledStateV,
	}
	i, ok := s2i[v]
	if ok {
		return i, nil
	}
	return ReconciledStateN, fmt.Errorf("Invalid ReconciledState: %q", v)
}

func (rs ReconciledState) String() string {

	defer func() {
		// in case of invalid ReconciledState returns ""
		if r := recover(); r != nil {
			return
		}
	}()

	var i2s = []string{"y", "n", "c", "f", "v"}
	return i2s[int(rs)]
}

// UnmarshalXML implements xml.Unmarshaler interface
func (rs *ReconciledState) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var v string
	d.DecodeElement(&v, &start)
	*rs, err = ReconciledStateFromString(v)
	return err
}
