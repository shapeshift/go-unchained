package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	ibclighttypes "github.com/cosmos/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/go-resty/resty/v2"
	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type CosmosService struct {
	chainConfig         ChainConfig
	encodingConfig      *params.EncodingConfig
	cosmosLightClient   *resty.Client
	tendermintRPCClient *resty.Client
	grpcConn            *grpc.ClientConn
	wsConn              *ws.Conn
}

type ChainConfig struct {
	CosmosBase          string
	TendermintBase      string
	TendermintWSBase    string
	GRPCBase            string
	WebsocketBase       string
	ApiKey              string // apiKey for node if required
	RestListenAddr      string // localhost:1660
	Bech32PrefixAccAddr string // cosmos
	Bech32PrefixAccPub  string // cosmospub
	RegisterTypes       func(registry types.InterfaceRegistry)
	EventHandler        func(tx *HistTx, evt TendermintEvent) []TxAction
}

type SubscribeMsg struct {
	JsonRpcMsg
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type NewBlockMsg struct { //result.data.value.block.height
	JsonRpcMsg
	Result BlockMsgResult `json:"result"`
}

type BlockMsgResult struct {
	Query string       `json:"query"`
	Data  BlockMsgData `json:"data"`
}
type BlockMsgData struct {
	Type  string        `json:"type"`
	Value BlockMsgValue `json:"value"`
}

type BlockMsgValue struct {
	Block BlockMsgBlock `json:"block"`
}

type BlockMsgBlock struct {
	Header BlockMsgHeader    `json:"header"`
	Data   BlockMsgBlockData `json:"data"`
}

type BlockMsgBlockData struct {
	Txs []string `json:"txs"`
}
type BlockMsgHeader struct {
	Version BlockMsgHeaderVersion `json:"version"`
	ChainID string                `json:"chain_id"`
	Height  string                `json:"height"`
	Time    string                `json:"time"`
}

type BlockMsgHeaderVersion struct {
	Block string `json:"block"`
	App   string `json:"app"`
}

// var (
// TODO - viper? other config means?
// apiKey         = os.Getenv("DATAHUB_API_KEY")
// cosmosBase     = os.Getenv("COSMOS_BASE")
// tendermintBase = os.Getenv("TENDERMINT_BASE")
// gRPCBase       = "cosmoshub-4--grpc--archive.grpc.datahub.figment.io"
// )

func NewCosmosService(chainConfig ChainConfig) (*CosmosService, error) {
	log.Infof("creating new CosmosClient with lcd: %s, rpc: %s", chainConfig.CosmosBase, chainConfig.TendermintBase)
	encodingConfig := MakeEncodingConfig(chainConfig)
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(chainConfig.Bech32PrefixAccAddr, chainConfig.Bech32PrefixAccPub)

	if chainConfig.RegisterTypes == nil {
		log.Infof("using default cosmos types only")
		chainConfig.RegisterTypes = func(registry types.InterfaceRegistry) {}
	}
	if chainConfig.EventHandler == nil {
		log.Infof("using default cosmos EventHandler only")
		chainConfig.EventHandler = func(tx *HistTx, evt TendermintEvent) []TxAction {
			return nil
		}
	}

	wsConn, err := wsConnect(chainConfig)
	if err != nil {
		log.Errorf("error connecting to tendermint websocket at %s: %s", chainConfig.TendermintWSBase, err)
	} else {
		if err = subscribeBlocks(wsConn); err != nil {
			log.Errorf("error connecting to tendermint websocket at %s: %s", chainConfig.TendermintWSBase, err)
		}
	}

	c := &CosmosService{
		chainConfig:    chainConfig,
		encodingConfig: &encodingConfig,
		cosmosLightClient: resty.New().SetBaseURL(chainConfig.CosmosBase).SetHeader("Accept", "application/json").
			SetHeader("Authorization", chainConfig.ApiKey).SetTimeout(10 * time.Second),
		tendermintRPCClient: resty.New().SetBaseURL(chainConfig.TendermintBase).SetHeader("Accept", "application/json").
			SetHeader("Authorization", chainConfig.ApiKey).SetTimeout(10 * time.Second),
		wsConn: wsConn,
	}

	go c.listen()

	return c, nil
}

func (c *CosmosService) listen() {
	if c.wsConn == nil {
		log.Errorf("no websocket connection to listen on")
		return
	}

	for {
		msg := NewBlockMsg{}
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Errorf("error reading json message: %s", err)
			break
		}
		if msg.Result.Data.Type != "tendermint/event/NewBlock" {
			log.Infof("not a block msg")
			continue
		}
		log.Infof("received block at height %s with %d txs", msg.Result.Data.Value.Block.Header.Height, len(msg.Result.Data.Value.Block.Data.Txs))

		for _, tx := range msg.Result.Data.Value.Block.Data.Txs {
			protoEncoded, err := base64.StdEncoding.DecodeString(tx)
			if err != nil {
				log.Errorf("error decoding tx %s: %s", tx, err)
			}

			sdkTx, err := c.encodingConfig.TxConfig.TxDecoder()(protoEncoded)
			if err != nil {
				log.Errorf("error decoding proto tx: %s", err)
				continue
			}

			msgs := sdkTx.GetMsgs()
			for _, msg := range msgs {
				log.Infof("received %s message", msg.Type())
			}
		}
	}
}

func subscribeBlocks(wsConn *ws.Conn) error {
	subscribeMessage := SubscribeMsg{
		JsonRpcMsg: JsonRpcMsg{
			ID:      0,
			Jsonrpc: "2.0",
		},
		Method: "subscribe",
		Params: []string{"tm.event='NewBlock'"},
	}

	if err := wsConn.WriteJSON(subscribeMessage); err != nil {
		log.Errorf("error sending subscribe message: %s", err)
		return fmt.Errorf("error sending subscribe message: %w", err)
	}
	return nil
}

func wsConnect(chainConfig ChainConfig) (*ws.Conn, error) {
	var wsDialer ws.Dialer
	wsUrl := fmt.Sprintf("%s/apikey/%s/websocket", chainConfig.TendermintWSBase, chainConfig.ApiKey)
	log.Infof("connecting to tendermint websocket at %s", wsUrl)
	wsConn, _, err := wsDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Errorf("error connecting ws: %s", err)
		return nil, err
	}
	log.Infof("successfully connected websocket")
	return wsConn, nil
}

func MakeEncodingConfig(config ChainConfig) params.EncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()

	// register protobuf types
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	disttypes.RegisterInterfaces(interfaceRegistry)
	govtypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	ibctypes.RegisterInterfaces(interfaceRegistry)
	ibccoretypes.RegisterInterfaces(interfaceRegistry)
	ibcchanneltypes.RegisterInterfaces(interfaceRegistry)
	ibclighttypes.RegisterInterfaces(interfaceRegistry)
	if config.RegisterTypes != nil {
		config.RegisterTypes(interfaceRegistry)
	}

	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             cdc,
	}
}

func isValidAddress(address string) bool {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false
	}
	return true
}

func doGet(client *resty.Client, path string, params map[string][]string, outputTo interface{}) error {
	url := path
	response, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParamsFromValues(params).
		Get(url)
	if err != nil {
		return fmt.Errorf("error fetching %s: %s", path, err)
	}

	// log.Debugf("response: %#v", string(response.Body()))
	if response.StatusCode() > 399 {
		log.Warnf("response: %#v", string(response.Body()))
		return fmt.Errorf("error requesting %s: %d", path, response.StatusCode())
	}

	if err = json.Unmarshal(response.Body(), outputTo); err != nil {
		log.Errorf("error unmarshalling response body: %s", string(response.Body()))
		return fmt.Errorf("error unmarshalling response for %s: %s", path, err)
	}
	return nil
}

func (c *CosmosService) GetCurrentHeight() int64 {
	// TODO read blockheight
	return 8626600
}
