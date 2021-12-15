package service

type CosmosPagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

type JsonRpcMsg struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
}

// Unchained API Types
type TokenAmount struct {
	// Example: uatom
	Denom string `json:"denom"`
	// Example: 420
	Amount string `json:"amount"`
}
