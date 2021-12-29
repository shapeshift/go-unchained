package service

import (
	"os"
	"testing"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	log "github.com/sirupsen/logrus"
)

const (
	address = "cosmos1ry8xqeyl7v02k6tlugjkzgvpz9ccpj6s4wt3nc"
)

func getCosmosService() *CosmosService {
	clt, _ := NewCosmosService(ChainConfig{
		CosmosBase:          "https://cosmoshub-4--lcd--archive.datahub.figment.io",
		TendermintBase:      "https://cosmoshub-4--rpc--archive.datahub.figment.io",
		ApiKey:              os.Getenv("DATAHUB_API_KEY"),
		Bech32PrefixAccAddr: "cosmos",
		Bech32PrefixAccPub:  "cosmospub",
		RegisterTypes: func(registry codectypes.InterfaceRegistry) {

		},
		EventHandler: func(tx *HistTx, evt TendermintEvent) []TxAction {
			return nil
		},
	})
	return clt
}

func TestWsConnect(t *testing.T) {
	cc := ChainConfig{
		CosmosBase:          "https://cosmoshub-4--lcd--archive.datahub.figment.io",
		TendermintBase:      "https://cosmoshub-4--rpc--archive.datahub.figment.io",
		TendermintWSBase:    "wss://cosmoshub-4--rpc--archive.datahub.figment.io",
		ApiKey:              os.Getenv("DATAHUB_API_KEY"),
		Bech32PrefixAccAddr: "cosmos",
		Bech32PrefixAccPub:  "cosmospub",
	}
	wsConnect(cc)
}

func TestGetAccount(t *testing.T) {
	c := getCosmosService()
	acct, err := c.GetAccount(address)
	if err != nil {
		t.Errorf("error getting account for address %s: %s", address, err)
		t.FailNow()
	}
	log.Infof("got acct: %#v", acct)
	if acct.Pubkey != address {
		t.Errorf("expected %s, got %s", address, acct.Pubkey)
	}
	log.Infof("acct: %d, seq: %d", acct.AccountNumber, acct.Sequence)
}
