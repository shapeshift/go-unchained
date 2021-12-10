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
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
