package thorchain

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/shapeshift/unchained-cosmos/service"
	thortypes "gitlab.com/thorchain/thornode/x/thorchain/types"
)

func RegisterTypes(registry codectypes.InterfaceRegistry) {
	thortypes.RegisterInterfaces(registry)
}

func ThorchainEventHandler(tx *service.HistTx, evt service.TendermintEvent) []service.TxAction {
	return nil
}
