package query

import (
	"strings"

	"github.com/mmbros/gnucash-viewer/model"
	"github.com/mmbros/treepath"
)

// FindAccounts returns the Accounts that matches the path string
func FindAccounts(path string, root *model.Account) ([]*model.Account, error) {
	// compile the path
	compiledpath, err := treepath.CompilePath(path)
	if err != nil {
		return nil, err
	}
	// find the elements
	elements := compiledpath.FindElements(&accountElement{root})
	if elements == nil || len(elements) == 0 {
		return nil, nil
	}
	// returns the accounts
	accounts := make([]*model.Account, len(elements))
	for j, e := range elements {
		accounts[j] = e.(*accountElement).Account
	}
	return accounts, nil
}

// accountElement implements treepath.Element interface for Account
type accountElement struct{ *model.Account }

// Parent returns the parent element.
// It returns nil in case of root element.
func (e accountElement) Parent() treepath.Element {
	par := accountElement{e.Account.Parent}
	return treepath.Element(&par)
}

// Children returns the children element of the current node
func (e accountElement) Children() []treepath.Element {
	elements := make([]treepath.Element, len(e.Account.Children))
	for j, c := range e.Account.Children {
		child := accountElement{c}
		elements[j] = &child
	}
	return elements
}

// MatchTag returns true if ...
func (e accountElement) MatchTag(tag string) bool {
	return strings.Contains(e.Name, tag)
}

// MatchTagText returns true if ...
func (e accountElement) MatchTagText(tag, text string) bool {
	return false
}

// MatchAttr returns true if ...
func (e accountElement) MatchAttr(attr string) bool {
	return false
}

// MatchAttrText returns true if ...
func (e accountElement) MatchAttrText(attr, text string) bool {
	return false
}
