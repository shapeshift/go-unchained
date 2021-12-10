package osmosis

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/shapeshift/unchained-cosmos/service"
	log "github.com/sirupsen/logrus"
)

const (
	address = "osmo1ry8xqeyl7v02k6tlugjkzgvpz9ccpj6sa4cp92"
)

func getOsmoService() *service.CosmosService {
	clt, _ := service.NewCosmosService(service.ChainConfig{
		CosmosBase:          "https://osmosis-1--lcd--archive.datahub.figment.io",
		TendermintBase:      "https://osmosis-1--rpc--archive.datahub.figment.io",
		GRPCBase:            os.Getenv("GRPC_BASE"),
		ApiKey:              os.Getenv("DATAHUB_OSMO_API_KEY"),
		Bech32PrefixAccAddr: "osmo",
		Bech32PrefixAccPub:  "osmopub",
		RegisterTypes:       RegisterTypes,
		EventHandler:        OsmosisEventHandler,
	})
	return clt
}

func TestGetAccount(t *testing.T) {
	var (
		c    *service.CosmosService
		acct *service.Account
		err  error
	)

	c = getOsmoService()
	acct, err = c.GetAccount(address)
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

func TestGetTxHistory(t *testing.T) {
	var (
		c    *service.CosmosService
		hist *service.TxHistory
		err  error
	)
	c = getOsmoService()
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
		c    *service.CosmosService
		tx   *service.HistTx
		err  error
	)

	txid = "4DF5CCCB36E88FFDCB40D2AA072C497934DF0BD06F37577297B2BA4987F8781B"
	c = getOsmoService()
	tx, err = c.GetTxByID(txid)
	if err != nil {
		t.Errorf("error getting tx for id %s: %s", txid, err)
		t.FailNow()
	}
	raw, err := json.Marshal(tx)
	if err != nil {
		t.Errorf("error marshaling tx: %s", err)
		t.FailNow()
	}
	log.Infof("got tx %s", string(raw))
}
