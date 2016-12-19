package types

import "testing"

func TestAccountType_String(t *testing.T) {
	const totAT = 20
	for j := 0; j < totAT; j++ {
		at := AccountType(j)
		name := at.String()
		t.Logf("%d - %s", j, name)

	}
}
