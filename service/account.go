package service

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

type CosmosAccount struct {
	Address       string     `json:"address"`
	PubKey        *types.Any `json:"public_key"`
	AccountNumber uint64     `json:"account_number,string"`
	Sequence      uint64     `json:"sequence,string"`
}

type CosmosAccountResponse struct {
	Account *CosmosAccount `json:"account"`
}

type CosmosBalance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type CosmosBalancesResponse struct {
	Balances   []CosmosBalance  `json:"balances"`
	Pagination CosmosPagination `json:"pagination"`
}

type CosmosDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

type CosmosDelegationWrapper struct {
	Delegation CosmosDelegation
	Balance    CosmosBalance
}

type CosmosDelegationsResponse struct {
	DelegationResponses []CosmosDelegationWrapper `json:"delegation_responses"`
	Pagination          CosmosPagination          `json:"pagination"`
}

type Delegation struct {
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
	Shares    string `json:"shares"`
}

func (c *CosmosService) readDelegations(address string) ([]CosmosDelegationWrapper, error) {
	res := &CosmosDelegationsResponse{}
	if err := doGet(c.cosmosLightClient, fmt.Sprintf("/cosmos/staking/v1beta1/delegations/%s", address), nil, res); err != nil {
		return nil, fmt.Errorf("error requesting delegations for %s: %s", address, err)
	}
	return res.DelegationResponses, nil
}

func (c *CosmosService) readBalances(address string) ([]CosmosBalance, error) {
	bankClient := banktypes.NewQueryClient(c.grpcConn)
	balanceRes, err := bankClient.AllBalances(
		context.Background(),
		&banktypes.QueryAllBalancesRequest{Address: address},
	)
	if err != nil {
		log.Errorf("grpc account error", err)
		return nil, err
	}

	cosmoBalances := make([]CosmosBalance, len(balanceRes.Balances))

	for i, bal := range balanceRes.Balances {

		cosmoBalances[i] = CosmosBalance{
			Denom:  bal.Denom,
			Amount: bal.Amount.String(),
		}
	}

	return cosmoBalances, nil
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
				balance = string(b.Amount)
			}
			t := TokenAmount{Denom: b.Denom, Amount: b.Amount}
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
				Shares:    cd.Delegation.Shares,
				Amount:    cd.Balance.Amount,
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
