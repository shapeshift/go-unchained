package main

import (
	"encoding/json"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

// encode to Stdout
var enc = json.NewEncoder(os.Stdout)

func TestGetAccount(t *testing.T) {
	pubkeys := []string{"Fabc123", "B123abc"}
	for _, pk := range pubkeys {
		acct := GetAccount(pk)

		// cast when client code needs specialized treatment
		switch v := acct.(type) {
		case *FoochainAccount:
			log.Infof("foo: %s", v.Foo)
		case BarchainAccount:
			log.Infof("bar: %s", v.Bar)
		default:
			t.Errorf("expected *FoochainAccount or BarchainAccount got %#v", v)
		}

		// in case of REST API, no cast just marshal leveraging built in encoding
		if err := enc.Encode(acct); err != nil {
			t.Errorf("error encoding acct %#v: %s", acct, err)
		}
	}
}
