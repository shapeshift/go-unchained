package docs

import (
	"github.com/shapeshift/unchained-cosmos/service"
)

// swagger:route GET /account/{pubkey} account-tag accountEndpointId
// Gets account information
// responses:
//   200: accountResponse
//   500: accountResponseError

// Account response contains informations about account (balance, sequence, etc)
// swagger:response accountResponse
type accountResponseWrapper struct {
	// in:body
	Body service.Account
}

// Account response error
// swagger:response accountResponseError
type accountResponseErrorWrapper struct {
	// in:body
	Body struct {
		// Example: error reading account for cosmos1fx4jwv3aalxqwmrpymn34l582lnehr3eqwuz9e
		Error string `json:"error"`
	}
}

// Account url param to specify pubkey (address) to get account info (balance, sequence, etc)
// swagger:parameters accountEndpointId
type accountParams struct {
	// Pubkey (address) to get account details for
	// in:path
	Pubkey string `json:"pubkey"`
}

// swagger:route GET /account/{pubkey}/txs txs-tag txsEndpointId
// Handles getting TX History for an account by pubkey (address)
// responses:
//   200: txsResponse
//   500: txsResponseError

// Tsxxs response contains tx history
// swagger:response txsResponse
type txsResponseWrapper struct {
	// in:body
	Body service.TxHistory
}

// Param to specify pubkey (address) for tx history
// swagger:parameters txsEndpointId
type txsParams struct {
	// Pubkey (address) to get tx history for
	// in:path
	Pubkey string `json:"pubkey"`
}

// Txs response error
// swagger:response txsResponseError
type txsResponseErrorWrapper struct {
	// in:body
	Body struct {
		// Example: error reading txhistory for recipient cosmos1fx4jwv3aalxqwmrpymn34l582lnehr3eqwuz9e
		Error string `json:"error"`
	}
}
