// Package numeric defines the Numeric type.
package types

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Numeric type is used to represents a GnuCash numeric type.
// Numeric{} equals 0 numeric number.
type Numeric struct {
	// Numerator
	num numint

	// Denominator
	// if den == 0 then the Numeric is 0
	// den is always >= 0
	den numint
}

// String returns the string representation of Numeric value.
func (n Numeric) String() string {
	switch n.den {
	case 0: // den == 0
		return "0"
	case 1: // den == 1
		return strconv.FormatInt(int64(n.num), 10)
	default:
		if n.num == 0 {
			return "0"
		}
		return fmt.Sprintf("%d/%d", n.num, n.den)
	}
}

// New creates a new numeric with numerator num and denominator den.
func New(num, den numint) *Numeric {
	if den < 0 {
		num, den = -num, -den
	}
	return &Numeric{num: num, den: den}
}

// Copy returns a new Numeric equals to x.
func Copy(x *Numeric) *Numeric {
	return &Numeric{x.num, x.den}
}

// Copy sets the existing Numeric n to the value of Numeric x.
func (n *Numeric) Copy(x *Numeric) {
	n.num, n.den = x.num, x.den
}

// FromString creates a new Numeric from string.
func FromString(v string) (*Numeric, error) {
	var n Numeric

	idx := strings.IndexByte(v, '/')
	if idx < 0 {
		num1, err := atoi(v)
		if err != nil {
			return nil, err
		}
		n.num = num1
		n.den = 1
	} else {
		num1, err := atoi(v[0:idx])
		if err != nil {
			return nil, err
		}
		den1, err := atoi(v[idx+1:])
		if err != nil {
			return nil, err
		}
		if den1 < 0 {
			num1, den1 = -num1, -den1
		}
		n.num = num1
		n.den = den1
	}
	return &n, nil
}

// IsZero return true if Numeric is zero.
//
// NOTE: Numeric{1, 0} == Numeric{10, 0} == Numeric{0, 0} == Numeric{0, 1}.
func (n *Numeric) IsZero() bool {
	// must be consistent con Sign func
	return (n.num == 0) || (n.den == 0)
}

// Equals returns true if z == x.
//
// NOTE: Numeric{1, 1} != Numeric{10, 10}.
func (n *Numeric) Equals(x *Numeric) bool {
	if n.IsZero() {
		return x.IsZero()
	}
	return (n.num == x.num) && (n.den == x.den)
}

// Sign returns:
//
//	-1 if z <  0
//	 0 if z == 0
//	+1 if z >  0
//
func (n *Numeric) Sign() int {
	if n.den == 0 || n.den == 0 {
		// must be consistent con IsZero func
		return 0
	}
	// assert(den > 0)
	if n.num > 0 {
		return 1
	}
	return -1
}

// NegEqual sets z to -z
func (n *Numeric) NegEqual() {
	n.num = -n.num
}

// AddEqual function: z.AddEqual(x) -> z += x
func (n *Numeric) AddEqual(x *Numeric) {

	if x.IsZero() {
		// n += 0
		return
	}
	if n.IsZero() {
		// 0 += x
		n.Copy(x)
		return
	}
	if n.den == x.den {
		n.num += x.num
		return
	}
	g := lcm(n.den, x.den)
	n.num = n.num*(g/n.den) + x.num*(g/x.den)
	n.den = g
}

// SubEqual function: z.SubEqual(x) -> z -= x
func (n *Numeric) SubEqual(x *Numeric) {
	y := Neg(x)
	n.AddEqual(y)
}

// Add function returns x+y.
func Add(x *Numeric, y *Numeric) *Numeric {
	z := Copy(x)
	z.AddEqual(y)
	return z
}

// Sub function returns x-y.
func Sub(x *Numeric, y *Numeric) *Numeric {
	z := Copy(x)
	z.SubEqual(y)
	return z
}

// Neg function sets x to -x.
func Neg(x *Numeric) *Numeric {
	return &Numeric{num: -x.num, den: x.den}
}

// Float64 converts the Numeric to a float64 value.
func (n *Numeric) Float64() float64 {
	if n.IsZero() {
		return 0.0
	}
	return float64(n.num) / float64(n.den)
}

// UnmarshalXML implements xml.Unmarshaler interface
func (n *Numeric) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// http://stackoverflow.com/questions/17301149/golang-xml-unmarshal-and-time-time-fields
	var v string
	d.DecodeElement(&v, &start)
	x, err := FromString(v)
	if err != nil {
		return err
	}
	n.Copy(x)

	return nil
}
