package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
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
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CosmosService struct {
	chainConfig         ChainConfig
	encodingConfig      *params.EncodingConfig
	cosmosLightClient   *resty.Client
	tendermintRPCClient *resty.Client
	grpcConn            *grpc.ClientConn
}

type ChainConfig struct {
	CosmosBase          string
	TendermintBase      string
	GRPCBase            string
	WebsocketBase       string
	ApiKey              string // apiKey for node if required
	RestListenAddr      string // localhost:1660
	Bech32PrefixAccAddr string // cosmos
	Bech32PrefixAccPub  string // cosmospub
	RegisterTypes       func(registry types.InterfaceRegistry)
	EventHandler        func(tx *HistTx, evt TendermintEvent) []TxAction
}

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
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

	grpcConn, err := grpc.DialContext(
		context.Background(),
		net.JoinHostPort(chainConfig.GRPCBase, "443"),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // The Cosmos SDK doesn't support any transport security mechanism.
		grpc.WithPerRPCCredentials(tokenAuth{
			token: chainConfig.ApiKey,
		}),
	)

	if err != nil {
		log.Fatal("error setting up grpc connection: %s", err)
	}

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
	c := &CosmosService{
		chainConfig:    chainConfig,
		encodingConfig: &encodingConfig,
		cosmosLightClient: resty.New().SetBaseURL(chainConfig.CosmosBase).SetHeader("Accept", "application/json").
			SetHeader("Authorization", chainConfig.ApiKey).SetTimeout(10 * time.Second),
		tendermintRPCClient: resty.New().SetBaseURL(chainConfig.TendermintBase).SetHeader("Accept", "application/json").
			SetHeader("Authorization", chainConfig.ApiKey).SetTimeout(10 * time.Second),
		grpcConn: grpcConn,
		// TODO - websockets
	}

	return c, nil
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
