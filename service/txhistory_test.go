package service

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGetTxHistory(t *testing.T) {
	var (
		c    *CosmosService
		hist *TxHistory
		err  error
	)
	c = getCosmosService()
	hist, err = c.GetTxHistory(address)
	if err != nil {
		t.Errorf("error getting tx history for address %s: %s", address, err)
		t.FailNow()
	}
	log.Infof("got hist: %#v", hist)
	if hist == nil {
		t.Errorf("no hist returned")
		t.FailNow()
	}
	if hist.Pubkey != address {
		t.Errorf("expected %s, got %s", address, hist.Pubkey)
	}
}

func TestGetTxByID(t *testing.T) {
	var (
		txid string
		c    *CosmosService
		tx   *HistTx
		err  error
	)
	txid = "FA5B8DDB65CFF3EDDCF4EC1C63622FD213DCC4DF52BEA5FC2A9AB5799BA52741"
	c = getCosmosService()
	tx, err = c.GetTxByID(txid)
	if err != nil {
		t.Errorf("error getting tx for id %s: %s", txid, err)
		t.FailNow()
	}
	log.Infof("got tx %#v", tx)
}

func TestParseAmountDenom(t *testing.T) {
	amountDenom := "1234567uatom"
	amt, denom, err := ParseAmountDenom(amountDenom)
	if err != nil {
		t.Errorf("error parsing %s: %s", amountDenom, err)
		t.FailNow()
	}
	if amt != "1234567" {
		t.Errorf("expected 1234567 got %s", amt)
	}

	if denom != "uatom" {
		t.Errorf("expected uatom got %s", denom)
	}

	amountDenom = "7uosmo"
	amt, denom, err = ParseAmountDenom(amountDenom)
	if err != nil {
		t.Errorf("error parsing %s: %s", amountDenom, err)
		t.FailNow()
	}
	if amt != "7" {
		t.Errorf("expected 1234567 got %s", amt)
	}

	if denom != "uosmo" {
		t.Errorf("expected uosmo got %s", denom)
	}

	amountDenom = "7253221"
	amt, denom, err = ParseAmountDenom(amountDenom)
	if err != nil {
		t.Errorf("error parsing %s: %s", amountDenom, err)
		t.FailNow()
	}
	if amt != "7253221" {
		t.Errorf("expected 7253221 got %s", amt)
	}

	if denom != "" {
		t.Errorf("expected empty denom got %s", denom)
	}

	amountDenom = ""
	amt, denom, err = ParseAmountDenom(amountDenom)
	if err != nil {
		t.Errorf("error parsing %s: %s", amountDenom, err)
		t.FailNow()
	}
	if amt != "" {
		t.Errorf("expected empty amt got %s", amt)
	}

	if denom != "" {
		t.Errorf("expected empty denom got %s", denom)
	}
}
