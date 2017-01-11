package query

import (
	"fmt"
	"time"

	"github.com/mmbros/gnucash-viewer/model"
	"github.com/mmbros/gnucash-viewer/types"
)

/*
// Transaction type
type Transaction struct {
	Currency    *Commodity     `xml:"currency"`
	DatePosted  types.Timespec `xml:"date-posted>date"`
	DateEntered types.Timespec `xml:"date-entered>date"`
	Description string         `xml:"description"`
	Splits      Splits         `xml:"splits>split"`
}


// Split type
type Split struct {
	ReconciledState types.ReconciledState `xml:"reconciled-state"`
	ReconcileDate   types.Timespec        `xml:"reconcile-date"`
	Value           types.Numeric         `xml:"value"`
	Memo            string                `xml:"memo"`
	Quantity        types.Numeric         `xml:"quantity"`
	Account         *Account              `xml:"-"`
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

*/

type Result struct {
	T *model.Transaction
	S *model.Split
}

type (
	TransactionFilter func(*model.Transaction) bool
	SplitFilter       func(*model.Split) bool
	AccountFilter     func(*model.Account) bool
)

type query struct {
	book               *model.Book
	transactionFilters []TransactionFilter
	splitFilters       []SplitFilter
	accountFilters     []AccountFilter
}

func Query(b *model.Book) *query {
	q := &query{
		b,

		[]TransactionFilter{},
		[]SplitFilter{},
		[]AccountFilter{},
	}
	return q
}

// ============================================================================
// TRANSACTION FILTERS

func (q *query) DatePostedRange(afterEqual, before time.Time) *query {

	fn := func(tr *model.Transaction) bool {
		datePosted := time.Time(tr.DatePosted)
		if datePosted.Before(afterEqual) {
			return false
		}
		return datePosted.Before(before)
	}

	q.transactionFilters = append(q.transactionFilters, fn)
	return q
}

func (q *query) DatePostedAfterEqual(afterEqual time.Time) *query {

	fn := func(tr *model.Transaction) bool {
		return !time.Time(tr.DatePosted).Before(afterEqual)
	}

	q.transactionFilters = append(q.transactionFilters, fn)
	return q
}

func (q *query) Currency(currency *model.Commodity) *query {

	fn := func(t *model.Transaction) bool {
		return t.Currency == currency
	}

	q.transactionFilters = append(q.transactionFilters, fn)
	return q
}

// ============================================================================
// SPLIT FILTERS

// ============================================================================
// ACCOUNT FILTERS

func (q *query) AccountType(at types.AccountType) *query {

	fn := func(a *model.Account) bool {
		return a.Type == at
	}

	q.accountFilters = append(q.accountFilters, fn)
	return q
}

func (q *query) AccountPath(path string) *query {

	accounts, err := FindAccounts(path, q.book.Accounts.Root)
	if err != nil {
		msg := fmt.Sprintf("query.AccountPath(%s): %v", path, err.Error())
		panic(msg)
	}
	fn := func(a *model.Account) bool {
		for _, acc := range accounts {
			if a == acc {
				return true
			}
		}
		return false
	}

	q.accountFilters = append(q.accountFilters, fn)
	return q
}

// ============================================================================
func (q *query) Execute() []*Result {
	res := []*Result{}

LOOP_TRANSACTION:
	for _, t := range q.book.Transactions {

		// Check fo transaction filters.
		for _, transactionFilter := range q.transactionFilters {
			if !transactionFilter(t) {
				// next Transaction
				continue LOOP_TRANSACTION
			}
		}
		// The current Transaction passed the transaction filters.

	LOOP_SPLIT:
		for _, s := range t.Splits {

			// Check for split filters.
			for _, splitFilter := range q.splitFilters {
				if !splitFilter(s) {
					// next Split
					continue LOOP_SPLIT
				}

			}

			// Check for account filters.
			for _, accountFilter := range q.accountFilters {
				if !accountFilter(s.Account) {
					// next Split
					continue LOOP_SPLIT
				}

			}

			// Success
			res = append(res, &Result{t, s})

		} // LOOP_SPLIT

	} // LOOP_TRANSACTION
	return res
}
