package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handles getting TX History for an account by pubkey (address)
func (s *UnchainedRestServer) GetTxHistory(writer http.ResponseWriter, req *http.Request) {
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

	hist, err := s.cosmosService.GetTxHistory(pubkey)
	if err != nil {
		log.Errorf("error reading tx history for %s: %s", pubkey, err)
		writeError(writer, ErrorResponse{fmt.Sprintf("error reading tx history for %s", pubkey)})
		return
	}

	res, err := json.Marshal(hist)
	if err != nil {
		log.Errorf("error marshaling tx history for %s: %s", pubkey, err)
		writeError(writer, ErrorResponse{fmt.Sprintf("error marshaling tx history to json for %s", pubkey)})
		return
	}

	if _, err := writer.Write([]byte(res)); err != nil {
		log.Errorf("error writing tx history response: %s", err)
	}
}
