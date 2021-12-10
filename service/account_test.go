package service

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestReadBalances(t *testing.T) {
	c := getCosmosService()
	bals, err := c.readBalances(address)
	if err != nil {
		t.Errorf("error reading balances: %s", err)
	}
	log.Infof("bals = %#v", bals)
}
