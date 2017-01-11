package query

import (
	"fmt"
	"time"

	"github.com/mmbros/gnucash-viewer/model"
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

type Filters struct {
	transactionFilters []TransactionFilter
	splitFilters       []SplitFilter
	accountFilters     []AccountFilter
}

func NewFilters() *Filters {
	f := Filters{
		[]TransactionFilter{},
		[]SplitFilter{},
		[]AccountFilter{},
	}
	return &f
}

func (f *Filters) String() string {
	return fmt.Sprintf("Filters{Transaction:%d, Split:%d, Account:%d}",
		len(f.transactionFilters),
		len(f.splitFilters),
		len(f.accountFilters))
}

// ============================================================================
// TRANSACTION FILTERS

func (f *Filters) DatePostedRange(afterEqual, before time.Time) {

	fn := func(tr *model.Transaction) bool {
		datePosted := time.Time(tr.DatePosted)
		if datePosted.Before(afterEqual) {
			return false
		}
		return datePosted.Before(before)
	}

	f.transactionFilters = append(f.transactionFilters, fn)
}

func (f *Filters) DatePostedAfterEqual(afterEqual time.Time) {

	fn := func(tr *model.Transaction) bool {
		return !time.Time(tr.DatePosted).Before(afterEqual)
	}

	f.transactionFilters = append(f.transactionFilters, fn)
}

func (f *Filters) Currency(currency *model.Commodity) {

	fn := func(t *model.Transaction) bool {
		return t.Currency == currency
	}

	f.transactionFilters = append(f.transactionFilters, fn)
}

// ============================================================================
// SPLIT FILTERS

// ============================================================================
func Where(b *model.Book, f *Filters) []*Result {
	res := []*Result{}

LOOP_TRANSACTION:
	for _, t := range b.Transactions {

		// Check fo transaction filters.
		for _, transactionFilter := range f.transactionFilters {
			if !transactionFilter(t) {
				// next Transaction
				continue LOOP_TRANSACTION
			}
		}
		// The current Transaction passed the transaction filters.

	LOOP_SPLIT:
		for _, s := range t.Splits {

			// Check for split filters.
			for _, splitFilter := range f.splitFilters {
				if !splitFilter(s) {
					// next Split
					continue LOOP_SPLIT
				}

			}

			// Check for account filters.
			for _, accountFilter := range f.accountFilters {
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
