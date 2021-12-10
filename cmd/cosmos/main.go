package main

import (
	// "net/http"
	// _ "net/http/pprof"
	"os"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/shapeshift/unchained-cosmos/server/rest"
	"github.com/shapeshift/unchained-cosmos/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Debug("")
	start()
}

func start() {
	// where to listen
	restListenAddr := os.Getenv("REST_LISTEN_ADDR")
	if restListenAddr == "" {
		restListenAddr = "localhost:1660"
	}

	config := service.ChainConfig{
		CosmosBase:          os.Getenv("COSMOS_BASE"),
		TendermintBase:      os.Getenv("TENDERMINT_BASE"),
		ApiKey:              os.Getenv("DATAHUB_API_KEY"),
		GRPCBase:            os.Getenv("GRPC_BASE"),
		RestListenAddr:      restListenAddr,
		Bech32PrefixAccAddr: "cosmos",
		Bech32PrefixAccPub:  "cosmospub",
		RegisterTypes: func(registry codectypes.InterfaceRegistry) {
			// cosmos so all stock protos
		},
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
