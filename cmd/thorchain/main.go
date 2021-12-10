package main

import (
	// "net/http"
	// _ "net/http/pprof"
	"os"

	"github.com/shapeshift/unchained-cosmos/server/rest"
	"github.com/shapeshift/unchained-cosmos/service"
	"github.com/shapeshift/unchained-cosmos/thorchain"
	log "github.com/sirupsen/logrus"
)

func main() {
	start()
}

func start() {
	// where to listen
	restListenAddr := os.Getenv("REST_LISTEN_ADDR")

	config := service.ChainConfig{
		CosmosBase:          os.Getenv("COSMOS_BASE"),
		TendermintBase:      os.Getenv("TENDERMINT_BASE"),
		ApiKey:              os.Getenv("DATAHUB_OSMO_API_KEY"),
		RestListenAddr:      restListenAddr,
		Bech32PrefixAccAddr: "thor",
		Bech32PrefixAccPub:  "thorpub",
		RegisterTypes:       thorchain.RegisterTypes,
		EventHandler:        thorchain.ThorchainEventHandler,
	}

	service, err := service.NewCosmosService(config)
	if err != nil {
		log.Errorf("dumping config: %#v", config)
		log.Fatalf("error creating CosmosService: %s", err)
	}
	srv, err := rest.New(service, config)
	if err != nil {
		log.Fatal("error starting http server: %s", err)
	}
	log.Infof("got srv: %#v", srv)
}
