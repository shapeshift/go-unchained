package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shapeshift/unchained-cosmos/service"

	log "github.com/sirupsen/logrus"
)

type UnchainedRestServer struct {
	cosmosService *service.CosmosService
	config        service.ChainConfig
}

type RestConfig struct {
	ListenAddr string
	Port       uint16
}

type Account struct {
	PubKey        string `json:"address"`
	AccountNumber uint64 `json:"account_number,string"`
	Sequence      uint64 `json:"sequence,string"`
}

type ErrorResponse struct {
	ErrorMsg string `json:"error"`
}

func New(cosmosService *service.CosmosService, config service.ChainConfig) (*UnchainedRestServer, error) {
	srv := &UnchainedRestServer{cosmosService, config}
	srv.start()
	return srv, nil
}

// start the server
func (s *UnchainedRestServer) start() {
	log.Info("starting http server...")
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/account/{pubkey}", s.GetAccount)
	router.HandleFunc("/account/{pubkey}/txs", s.GetTxHistory)

	listenAddr := s.config.RestListenAddr
	if listenAddr == "" {
		listenAddr = "localhost:1660"
	}

	log.Infof("invoking ListenAndServe on %s", listenAddr)
	log.Errorf("ListenAndServe returned: %s", http.ListenAndServe(listenAddr, router))
}

func writeError(writer http.ResponseWriter, errRes ErrorResponse) error {
	var (
		bytes []byte
		err   error
	)

	writer.WriteHeader(500)
	if bytes, err = json.Marshal(errRes); err != nil {
		log.Errorf("error marshaling error response: %s", errRes)
		writer.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", errRes.ErrorMsg)))
		return nil
	}

	if _, err := writer.Write(bytes); err != nil {
		log.Errorf("error writing error response: %s", errRes)
		return err
	}
	return nil
}
