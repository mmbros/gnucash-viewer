package types

import "encoding/xml"

// AccountType enum type
type AccountType int

// AccountType constants
const (
	AccountTypeNone AccountType = iota
	AccountTypeBank
	AccountTypeCash
	AccountTypeCredit
	AccountTypeAsset
	AccountTypeLiability
	AccountTypeStock
	AccountTypeMutual
	AccountTypeCurrency
	AccountTypeIncome
	AccountTypeExpense
	AccountTypeEquity
	AccountTypeReceivable
	AccountTypePayable
	AccountTypeRoot
	AccountTypeTrading
	AccountTypeChecking
	AccountTypeSavings
	AccountTypeMoneyMrkt
	AccountTypeCreditLine
)

type infoAccountType struct {
	label        string
	root         bool
	invertValues bool
	plusLabel    string
	minusLabel   string
}

func AccountTypeFromString(v string) AccountType {
	var s2i = map[string]AccountType{
		"NONE":       AccountTypeNone,
		"BANK":       AccountTypeBank,
		"CASH":       AccountTypeCash,
		"CREDIT":     AccountTypeCredit,
		"ASSET":      AccountTypeAsset,
		"LIABILITY":  AccountTypeLiability,
		"STOCK":      AccountTypeStock,
		"MUTUAL":     AccountTypeMutual,
		"CURRENCY":   AccountTypeCurrency,
		"INCOME":     AccountTypeIncome,
		"EXPENSE":    AccountTypeExpense,
		"EQUITY":     AccountTypeEquity,
		"RECEIVABLE": AccountTypeReceivable,
		"PAYABLE":    AccountTypePayable,
		"ROOT":       AccountTypeRoot,
		"TRADING":    AccountTypeTrading,
		"CHECKING":   AccountTypeChecking,
		"SAVINGS":    AccountTypeSavings,
		"MONEYMRKT":  AccountTypeMoneyMrkt,
		"CREDITLINE": AccountTypeCreditLine,
	}
	at, ok := s2i[v]
	if ok {
		return at
	}
	return AccountTypeNone
}

// UnmarshalXML implements xml.Unmarshaler interface
func (at *AccountType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	*at = AccountTypeFromString(v)
	return nil
}

func (at AccountType) String() string {
	return infoAccountTypes[at].label
}

func (at AccountType) PlusLabel() string {
	return infoAccountTypes[at].plusLabel
}

func (at AccountType) MinusLabel() string {
	return infoAccountTypes[at].minusLabel
}

func (at AccountType) InvertValues() bool {
	return infoAccountTypes[at].invertValues
}

func (at AccountType) Root() bool {
	return infoAccountTypes[at].root
}

// infoAccountTypes is the map of all infoAccountType
var infoAccountTypes = map[AccountType]*infoAccountType{
	AccountTypeNone: &infoAccountType{
		label: "None",
	},
	AccountTypeBank: &infoAccountType{
		label:      "Bank",
		plusLabel:  "Deposit",
		minusLabel: "Withdrawal",
	},
	AccountTypeCash: &infoAccountType{
		label:      "Cash",
		plusLabel:  "Receive",
		minusLabel: "Spend",
	},
	AccountTypeCredit: &infoAccountType{
		label:      "Credit",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeAsset: &infoAccountType{
		label:      "Asset",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeLiability: &infoAccountType{
		label:        "Liability",
		invertValues: true,
		plusLabel:    "Decrease",
		minusLabel:   "Increase",
	},
	AccountTypeStock: &infoAccountType{
		label:      "Stock",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeMutual: &infoAccountType{
		label:      "Mutual",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeCurrency: &infoAccountType{
		label:      "Currency",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeIncome: &infoAccountType{
		label:        "Income",
		invertValues: true,
		plusLabel:    "Charge",
		minusLabel:   "Income",
	},
	AccountTypeExpense: &infoAccountType{
		label:      "Expense",
		plusLabel:  "Expense",
		minusLabel: "Rebate",
	},
	AccountTypeEquity: &infoAccountType{
		label:        "Equity",
		invertValues: true,
		plusLabel:    "Decrease",
		minusLabel:   "Increase",
	},
	AccountTypeReceivable: &infoAccountType{
		label:      "Receivible",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypePayable: &infoAccountType{
		label:      "Payable",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeRoot: &infoAccountType{
		label: "Root",
		root:  true,
	},
	AccountTypeTrading: &infoAccountType{
		label:      "Trading",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeChecking: &infoAccountType{
		label:      "Checking",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeSavings: &infoAccountType{
		label:      "Savings",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeMoneyMrkt: &infoAccountType{
		label:      "MoneyMrkt",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
	AccountTypeCreditLine: &infoAccountType{
		label:      "CreditLine",
		plusLabel:  "Increase",
		minusLabel: "Decrease",
	},
}
