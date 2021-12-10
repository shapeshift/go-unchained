package osmosis

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	gammtypes "github.com/osmosis-labs/osmosis/x/gamm/types"
	lockuptypes "github.com/osmosis-labs/osmosis/x/lockup/types"
	"github.com/shapeshift/unchained-cosmos/service"
	log "github.com/sirupsen/logrus"
)

func RegisterTypes(registry codectypes.InterfaceRegistry) {
	gammtypes.RegisterInterfaces(registry)
	lockuptypes.RegisterInterfaces(registry)
}

// Osmosis specific event (log) handlers
func OsmosisEventHandler(tx *service.HistTx, evt service.TendermintEvent) []service.TxAction {
	txActions := make([]service.TxAction, 0, 5)
	attribs := service.ToAttributeMap(evt.Attributes)
	switch evt.Type {
	case "transfer":
		log.Infof("transfer: %v", attribs)
		if tx.Type != "swap_exact_amount_in" { // swap_exact_amount_out
			// don't mess with other types of transfer events
			break
		}

		xferIn := service.TransferTxAction{BaseTxAction: service.BaseTxAction{Type: evt.Type}}
		xferOut := service.TransferTxAction{BaseTxAction: service.BaseTxAction{Type: evt.Type}}
		// 2 xfers, in and out
		for _, attrib := range evt.Attributes {
			if attrib.Key == "sender" {
				if xferOut.FromAddress == "" {
					xferOut.FromAddress = attrib.Value
				} else {
					xferIn.FromAddress = attrib.Value
				}
			}

			if attrib.Key == "recipient" {
				if xferOut.ToAddress == "" {
					xferOut.ToAddress = attrib.Value
				} else {
					xferIn.ToAddress = attrib.Value
				}
			}

			if attrib.Key == "amount" {
				amount, denom, err := service.ParseAmountDenom(attrib.Value)
				if err != nil {
					log.Errorf("error parsing transfer amount (%s) for swap: %s", attrib.Value, err)
					amount, denom = "", ""
				}
				if xferOut.Amount == "" {
					xferOut.Amount = amount
					xferOut.Asset = denom
				} else {
					xferIn.Amount = amount
					xferIn.Asset = denom
				}
			}
		}

		txActions = append(txActions, xferIn, xferOut)
	case "swap_exact_amount_out":
		log.Infof("swap_exact_amount_out: %v", attribs)
		tx.Type = evt.Type
	case "join_pool": // believe this falls under the custom transfer handler
		log.Infof("swap_exact_amount_out: %v", attribs)
		tx.Type = evt.Type
	case "claim":
		log.Infof("claim: %v", attribs)
		tx.Type = evt.Type
	default:
	}
	return txActions
}
