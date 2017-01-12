package model

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/mmbros/gnucash-viewer/types"
)

/*
Account = element gnc:account {
  attribute version { "2.0.0" },
  element act:name { text },
  element act:id { attribute type { "guid" }, GUID },

  # from xaccAccountTypeEnumAsString in src/engine/Account.c

	element act:type { "NONE"
					| "MONEYMRKT"
					| "CREDITLINE" },

	( element act:commodity {
			element cmdty:space { text },
			element cmdty:id { text }
		},
		element act:commodity-scu { xsd:int },  // SCU = Smallest Commodity Unit, in most cases 100
		element act:non-standard-scu { empty }?
	)?,
	element act:code { text }?,
	element act:description { text }?,
	element act:slots { KvpSlot+ }?,
	element act:parent { attribute type { "guid" }, GUID }?,
	element act:lots { Lot+ }?
}

*/

// Accounts represents the Accounts tree.
type Accounts struct {
	Root *Account
	//Map  map[types.GUID]*Account
	List []*Account
}

// Account type
type Account struct {
	Type        types.AccountType
	Name        string
	Description string
	Currency    *Commodity

	// SCU = Smallest Commodity Unit, in most cases 100
	CommodityScu   int
	NonStandardScu bool

	Parent   *Account
	Children []*Account
}

// AccountMap type represents a mapping from a string key to an account item.
type AccountMap map[string]*Account

func AccountUnmarshalXML(decoder *xml.Decoder, commodities Commodities) (a *Account, ID string, ParentID string, err error) {

	acc := Account{}

LOOP:
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break LOOP
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			var v interface{}

			switch se.Name.Local {
			case "commodity":
				var cmdty Commodity
				if err = decoder.DecodeElement(&cmdty, &se); err != nil {
					return
				}
				acc.Currency = commodities.Get(cmdty.Space, cmdty.ID)
			case "name":
				v = &acc.Name
			case "description":
				v = &acc.Description
			case "id":
				v = &ID
			case "parent":
				v = &ParentID
			case "type":
				v = &acc.Type
			case "commodity-scu":
				v = &acc.CommodityScu
			case "non-standard-scu":
				v = &acc.NonStandardScu
			}

			if v != nil {
				if err = decoder.DecodeElement(v, &se); err != nil {
					return
				}
			}
		case xml.EndElement:
			if se.Name.Local == "account" {
				break LOOP
			}
		}
	}
	if ID == "" || acc.Name == "" {
		err = fmt.Errorf("Account without ID or Name")
		return
	}

	a = &acc
	return
}

/*
func (accounts *Accounts) Init(acclist []*Account) error {
	if accounts == nil {
		return errors.New("Accounts must not be nil")
	}

	accounts.Map = map[types.GUID]*Account{}
	var root *Account
	m := accounts.Map

	// populate accounts map
	for _, a := range acclist {
		m[a.ID] = a
		//a.Children = []*Account{}
	}
	// set the root and the parent element of every account
	for _, a := range acclist {
		if p, ok := m[a.ParentID]; ok {
			a.Parent = p
			p.Children = append(p.Children, a)
			continue
		}
		if a.Type != types.AccountTypeRoot {
			return errors.New("Account without parent and type != root")
		}
		if root != nil {
			return errors.New("Multiple root accounts")
		}
		root = a
	}
	if root == nil {
		return errors.New("Root not found")
	}
	accounts.Root = root

	return nil
}
func NewAccounts(acclist []*Account) (*Accounts, error) {
	var root *Account
	m := map[types.GUID]*Account{}

	// populate accounts map
	for _, a := range acclist {
		m[a.ID] = a
		//a.Children = []*Account{}
	}
	// set the root and the parent element of every account
	for _, a := range acclist {
		if p, ok := m[a.ParentID]; ok {
			a.Parent = p
			p.Children = append(p.Children, a)
			continue
		}
		if a.Type != types.AccountTypeRoot {
			return nil, errors.New("Account without parent and type != root")
		}
		if root != nil {
			return nil, errors.New("Multiple root accounts")
		}
		root = a
	}
	if root == nil {
		return nil, errors.New("Root not found")
	}
	return &Accounts{Map: m, Root: root}, nil
}

*/
// PrintTree prints account tree
func (accounts *Accounts) PrintTree(indent string) {
	if indent == "" {
		indent = "  "
	}

	if accounts == nil {
		fmt.Println("<nil>")
		return
	}

	if accounts.Root == nil {
		fmt.Println("<root-nil>")
		return
	}
	var pr func(*Account, int, string)

	pr = func(a *Account, level int, indent string) {
		fmt.Printf("%s[%s] %s (%s)\n", strings.Repeat(indent, level), a.Type, a.Name, a.Currency)
		for _, child := range a.Children {
			pr(child, level+1, indent)
		}
	}

	pr(accounts.Root, 0, indent)

}

// Add append an account to the list of accounts.
// It doesn't check if the account is already present in the list.
func (accounts *Accounts) Add(acc *Account) {
	accounts.List = append(accounts.List, acc)
}

// Len returns the number of the accounts.
func (accounts *Accounts) Len() int {
	return len(accounts.List)
}

// ancestors returns the list of the ancestors of the given account.
// The first item of the list is the given account itself.
// The last item of the list is the root account.
//    -> {account, parent, grandparent, ..., root}
func (a *Account) ancestors() []*Account {
	accounts := []*Account{}
	for a != nil {
		accounts = append(accounts, a)
		a = a.Parent
	}
	return accounts
}

func (a *Account) getFullName(sep string) string {
	if a == nil {
		return "<nil>"
	}
	if a.Parent == nil {
		return "<root>"
	}

	accounts := a.ancestors()
	buf := new(bytes.Buffer)
	for j := len(accounts) - 2; j >= 0; j-- {
		buf.WriteString(accounts[j].Name)
		if j > 0 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

func (a *Account) FullName() string {
	return a.getFullName(" > ")
}

// BasicType returns the account type of the ancestor tha is a child of the root account.
// returns None in case of nil account.
// returns Root in case of root account.
func (a *Account) BasicType() types.AccountType {
	if a == nil {
		return types.AccountTypeNone
	}
	p := a.Parent
	if p == nil {
		return a.Type // Root
	}
	for {
		if p.Parent == nil {
			return a.Type
		}
		p, a = p.Parent, p
	}
}
