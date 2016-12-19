package types

import "strconv"

// base type of Numeric
type numint int64

// atoi converts a string to a numint
func atoi(s string) (numint, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	return numint(i), err
}

// abs returns the absolute value of i.
func abs(i numint) numint {
	if i < 0 {
		return -i
	}
	return i
}

// gcd returns the greatest common divisor of a and b.
//   gcd(0, b) -> |b|
//   gcd(a, 0) -> |a|
func gcd(a, b numint) numint {
	a = abs(a)
	b = abs(b)
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm returns the least common multiple of a and b.
//   lcm(0, 0) -> 0
//   lcm(a, 0) -> 0
//   lcm(0, b) -> 0
func lcm(a, b numint) numint {
	if a == 0 || b == 0 {
		return 0
	}
	a = abs(a)
	b = abs(b)
	g := gcd(a, b)
	l := (a / g) * b

	return l
}
