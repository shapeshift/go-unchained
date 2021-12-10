package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handles getting an account by pubkey (address)
func (s *UnchainedRestServer) GetAccount(writer http.ResponseWriter, req *http.Request) {
	var (
		pubkey string
		ok     bool
	)

	header := writer.Header()
	header.Set("Content-Type", "application/json")

	pubkey, ok = mux.Vars(req)["pubkey"]
	if !ok || pubkey == "" {
		writeError(writer, ErrorResponse{"pubkey is required"})
		return
	}

	acct, err := s.cosmosService.GetAccount(pubkey)
	if err != nil {
		log.Errorf("error reading account for %s: %s", pubkey, err)
		writeError(writer, ErrorResponse{fmt.Sprintf("error reading account for %s", pubkey)})
		return
	}

	res, err := json.Marshal(acct)
	if err != nil {
		log.Errorf("error marshaling account for %s: %s", pubkey, err)
		writeError(writer, ErrorResponse{fmt.Sprintf("error marshaling account to json for %s", pubkey)})
		return
	}

	if _, err := writer.Write([]byte(res)); err != nil {
		log.Errorf("error writing account response: %s", err)
	}
}
