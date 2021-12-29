package main

import (
	"strings"
)

// minimal requirements for cosmos account
type IAccount interface {
	GetPubkey() string
	GetBalance() string
	GetAcctNumber() string
	GetSequence() string
}

// base implementation of cosmos account
type BaseAccount struct {
	Pubkey     string `json:"pubkey"`
	Balance    string `json:"balance"`
	Sequence   string `json:"sequence"`
	AcctNumber string `json:"acctNumber"`
}

// Accessors implement IAccount interface on BaseAccount struct
func (acct BaseAccount) GetPubkey() string {
	return acct.Pubkey
}

func (acct BaseAccount) GetBalance() string {
	return acct.Balance
}

func (acct BaseAccount) GetAcctNumber() string {
	return acct.AcctNumber
}

func (acct BaseAccount) GetSequence() string {
	return acct.AcctNumber
}

// "extend" BaseAccount via composition
type FoochainAccount struct {
	BaseAccount
	Foo string `json:"foo"`
}

// "extend" BaseAccount via composition
type BarchainAccount struct {
	BaseAccount
	Bar string `json:"bar"`
}

// To me the value is here; GetAccount can return various concrete types
// appending data specific to the chain
func GetAccount(pubkey string) IAccount {
	if strings.HasPrefix(pubkey, "F") {
		return &FoochainAccount{
			BaseAccount: BaseAccount{
				Pubkey:     pubkey,
				Balance:    "123456789",
				AcctNumber: "456",
				Sequence:   "660",
			},
			Foo: "Ima Foo",
		}
	}
	return BarchainAccount{
		BaseAccount: BaseAccount{
			Pubkey:     pubkey,
			Balance:    "987654321",
			AcctNumber: "654",
			Sequence:   "17",
		},
		Bar: "Ima Bar",
	}
}

func main() {

}
