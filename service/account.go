package service

import (
	"context"
	"fmt"

	cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	log "github.com/sirupsen/logrus"
)

type Account struct {
	Pubkey        string        `json:"pubkey"`
	AccountNumber uint64        `json:"accountNumber,string"`
	Sequence      uint64        `json:"sequence,string"`
	Balance       string        `json:"balance"`
	Tokens        []TokenAmount `json:"tokens"`
	Delegations   []Delegation  `json:"delegations"`
}

type Delegation struct {
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
	Shares    string `json:"shares"`
}

func (c *CosmosService) readDelegations(address string) ([]stakingtypes.DelegationResponse, error) {

	delegationClient := stakingtypes.NewQueryClient(c.grpcConn)
	delegationRes, err := delegationClient.DelegatorDelegations(
		context.Background(),
		&stakingtypes.QueryDelegatorDelegationsRequest{DelegatorAddr: address},
	)
	if err != nil {
		log.Errorf("grpc balances error", err)
		return nil, err
	}

	return delegationRes.DelegationResponses, nil
}

func (c *CosmosService) readBalances(address string) (cosmos_cosmos_sdk_types.Coins, error) {
	bankClient := banktypes.NewQueryClient(c.grpcConn)
	balanceRes, err := bankClient.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{Address: address},
	)
	if err != nil {
		log.Errorf("grpc balances error", err)
		return nil, err
	}

	return balanceRes.Balances, nil
}

func (c *CosmosService) readAccount(address string) (*authtypes.BaseAccount, error) {
	authClient := authtypes.NewQueryClient(c.grpcConn)
	authRes, err := authClient.Account(
		context.Background(),
		&authtypes.QueryAccountRequest{Address: address},
	)
	if err != nil {
		log.Errorf("grpc account error", err)
		return nil, err
	}
	account := authtypes.BaseAccount{}
	err = c.encodingConfig.Marshaler.UnmarshalBinaryBare(authRes.GetAccount().Value, &account)
	if err != nil {
		log.Errorf("unmarshal error", err)
	}
	return &account, nil
}

func (c *CosmosService) GetAccount(address string) (*Account, error) {
	if !isValidAddress(address) {
		return nil, fmt.Errorf("invalid address supplied: %s", address)
	}

	cosmosAcct, err := c.readAccount(address)
	if err != nil {
		return nil, err
	}

	balance := "0"
	cosmosBals, err := c.readBalances(address)
	tokens := make([]TokenAmount, 0, len(cosmosBals))
	if err != nil {
		log.Errorf("error getting balances for %s: %s", address, err)
	} else {
		for _, b := range cosmosBals {
			if b.Denom == "uatom" {
				balance = string(b.Amount.String())
			}
			t := TokenAmount{Denom: b.Denom, Amount: b.Amount.String()}
			tokens = append(tokens, t)
		}
	}

	delegations := make([]Delegation, 0, 25)

	cosmosDelegations, err := c.readDelegations(address)

	if err != nil {
		log.Errorf("error getting delegations for %s: %s", address, err)
	} else {
		for _, cd := range cosmosDelegations {
			d := Delegation{
				Validator: cd.Delegation.ValidatorAddress,
				Shares:    cd.Delegation.Shares.String(),
				Amount:    cd.Balance.Amount.String(),
			}
			delegations = append(delegations, d)
		}
	}

	acct := Account{
		Pubkey:        cosmosAcct.Address,
		AccountNumber: cosmosAcct.AccountNumber,
		Sequence:      cosmosAcct.Sequence,
		Balance:       balance,
		Tokens:        tokens,
		Delegations:   delegations,
	}
	return &acct, nil
}
