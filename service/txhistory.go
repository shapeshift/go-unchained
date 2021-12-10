package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibccoretypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	log "github.com/sirupsen/logrus"
)

type TendermintTxResult struct {
	Code      int    `json:"code"`
	Data      string `json:"data"`
	Log       string `json:"log"`
	Info      string `json:"info"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	// Events []TendermintEvent
	Codespace string `json:"codespace"`
}

type TendermintHistTx struct {
	Hash     string             `json:"hash"`
	Height   string             `json:"height"`
	Index    uint               `json:"index"`
	TxResult TendermintTxResult `json:"tx_result"`
	Tx       string             `json:"tx"`
}
type TendermintTxSearchResponse struct {
	JsonRpcMsg
	Result struct {
		Txs []TendermintHistTx `json:"txs"` // a
	} `json:"result"`
}

type TendermintAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type TendermintEvent struct {
	Type       string `json:"type"`
	Attributes []TendermintAttribute
}

type TendermintLog struct {
	MsgIndex uint              `json:"msg_index"`
	Events   []TendermintEvent `json:"events"`
}

type TxAction interface {
	GetType() string
}
type BaseTxAction struct {
	Type string `json:"type"`
}

func (action BaseTxAction) GetType() string {
	return action.Type
}

type TransferTxAction struct {
	BaseTxAction
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Asset       string `json:"asset"`
	Amount      string `json:"amount"`
}

type DelegateTxAction struct {
	BaseTxAction
	Validator string `json:"validator"`
	Amount    string `json:"amount"`
}

type ReDelegateTxAction struct {
	BaseTxAction
	SourceValidator string `json:"sourceValidator"`
	DestValidator   string `json:"destValidator"`
	Amount          string `json:"amount"`
}

type WithdrawRewardsTxAction struct {
	BaseTxAction
	Validator string `json:"validator"`
	Asset     string `json:"asset"`
	Amount    string `json:"amount"`
}

type IBCTransferTxAction struct {
	BaseTxAction
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
}

// Unchained API Historical Tx Model
type HistTx struct {
	Type          string     `json:"type"`
	TxID          string     `json:"txid"`
	Status        string     `json:"status"`
	From          string     `json:"from,omitempty"`
	To            string     `json:"to,omitempty"`
	BlockHeight   int64      `json:"blockHeight"`
	Index         uint       `json:"index"`
	Confirmations int64      `json:"confirmations"`
	Timestamp     int64      `json:"timestamp"`
	Value         string     `json:"value,omitempty"`
	Memo          string     `json:"memo,omitempty"`
	Fee           string     `json:"fee"`
	FeeAsset      string     `json:"feeAsset"`
	GasWanted     string     `json:"gasWanted"`
	GasUsed       string     `json:"gasUsed"`
	Actions       []TxAction `json:"actions"`
}

type TxHistory struct {
	// TODO - pagination
	Pubkey string   `json:"pubkey"`
	Txs    []HistTx `json:"txs"`
}

func (tx *HistTx) applyLogs(rawLogs string, chainConfig ChainConfig) error {
	var (
		err error
	)
	if rawLogs == "" {
		return fmt.Errorf("no logs present")
	}

	logs := make([]TendermintLog, 0, 100)
	if err = json.Unmarshal([]byte(rawLogs), &logs); err != nil {
		return fmt.Errorf("error unmarshaling logs %s", err)
	}

	for _, sourceLog := range logs {
		var (
			action    string
			module    string
			sender    string
			recipient string
			validator string
			amount    string
			denom     string
		)
		for _, evt := range sourceLog.Events {
			if chainConfig.EventHandler != nil {
				txActions := chainConfig.EventHandler(tx, evt)
				if txActions != nil && len(txActions) > 0 {
					tx.Actions = append(tx.Actions, txActions...)
					continue
				}
			}
			attribs := ToAttributeMap(evt.Attributes)
			switch evt.Type {
			case "message":
				action = attribs["action"]
				module = attribs["module"]
				tx.Type = action
				log.Debugf("%s %s message", module, action)
			case "transfer":
				sender = attribs["sender"]
				recipient = attribs["recipient"]
				amount, denom, err = ParseAmountDenom(attribs["amount"])
				if err != nil {
					log.Errorf("error parsing transfer denomAmount %s: %s", attribs["amount"], err)
					amount = "0"
				}
				txAction := TransferTxAction{
					BaseTxAction: BaseTxAction{Type: evt.Type},
					FromAddress:  sender,
					ToAddress:    recipient,
					Asset:        denom,
					Amount:       amount,
				}
				tx.Actions = append(tx.Actions, txAction)
				log.Debugf("%s %s -> %s", sender, amount, recipient)
			case "delegate":
				amount = attribs["amount"]
				validator = attribs["validator"]
				txAction := DelegateTxAction{
					BaseTxAction: BaseTxAction{Type: evt.Type},
					Validator:    validator,
					Amount:       amount,
				}
				tx.Actions = append(tx.Actions, txAction)
			case "redelegate":
				sourceValidator := attribs["source_validator"]
				destValidator := attribs["destination_validator"]
				amount := attribs["amount"]
				txAction := ReDelegateTxAction{
					BaseTxAction:    BaseTxAction{Type: evt.Type},
					SourceValidator: sourceValidator,
					DestValidator:   destValidator,
					Amount:          amount,
				}
				tx.Actions = append(tx.Actions, txAction)
			case "withdraw_rewards":
				amount, denom, err = ParseAmountDenom(attribs["amount"])
				if err != nil {
					log.Errorf("error parsing withdraw_rewards denomAmount %s: %s", attribs["amount"], err)
					amount = "0"
				}

				validator = attribs["validator"]
				txAction := WithdrawRewardsTxAction{
					BaseTxAction: BaseTxAction{Type: evt.Type},
					Validator:    validator,
					Asset:        denom,
					Amount:       amount,
				}
				tx.Actions = append(tx.Actions, txAction)
			case "ibc_transfer":
				sender := attribs["sender"]
				receiver := attribs["receiver"]
				txAction := IBCTransferTxAction{
					BaseTxAction: BaseTxAction{Type: evt.Type},
					FromAddress:  sender,
					ToAddress:    receiver,
				}
				tx.Actions = append(tx.Actions, txAction)
			case "send_packet":
				log.Debug("")
			case "fungible_token_packet":
				log.Debug("")
			case "recv_packet":
				log.Debug("")
				packetData := attribs["packet_data"]
				// "{\"amount\":\"9129485\",\"denom\":\"transfer/channel-0/uatom\",\"receiver\":\"cosmos1ry8xqeyl7v02k6tlugjkzgvpz9ccpj6s4wt3nc\",\"sender\":\"osmo1ry8xqeyl7v02k6tlugjkzgvpz9ccpj6sa4cp92\"}"
				_ = packetData
			case "write_acknowledgement":
				log.Debug("")
			case "update_client":
				log.Debug("")
			default:
				log.Warnf("unhandled event type %s", evt.Type)
			}
		}
	}

	return nil
}

func (tx *HistTx) applyMsg(msg sdk.Msg) {
	// keep txType derived from first msg in tx
	if tx.Type == "" {
		tx.Type = "unknown"
	}

	switch v := msg.(type) {
	case *banktypes.MsgSend:
		tx.Type = "transfer"
		amount, _, err := ParseAmountDenom(v.Amount.String())
		if err != nil {
			amount = "0"
			log.Errorf("tx %s error parsing amount %s: %s", v.Amount.String(), tx.TxID, err)
		}
		tx.From = v.FromAddress
		tx.To = v.ToAddress
		tx.Value = amount
	case *stakingtypes.MsgDelegate:
		tx.Type = "delegate"
		tx.From = v.DelegatorAddress
	case *stakingtypes.MsgBeginRedelegate:
		tx.Type = "redelegate"
		tx.From = v.DelegatorAddress
	case *disttypes.MsgWithdrawDelegatorReward:
		tx.Type = "withdrawRewards"
		tx.From = v.DelegatorAddress
	case *ibctypes.MsgTransfer:
		tx.Type = "ibcTransfer"
		amount := v.Token.Amount.String()
		tx.From = v.Sender
		tx.To = v.Receiver
		tx.Value = amount
	case *ibccoretypes.MsgUpdateClient:
		tx.Type = "ibcUpdateClient"
	case *ibcchanneltypes.MsgRecvPacket:
		tx.Type = "ibcRecvPacket"
	default:
		log.Warnf("tx %s unsupported msg type %s", tx.TxID, msg.Type())
	}
}

func (c *CosmosService) decodeTx(raw string) (sdk.Tx, error) {
	protoEncoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("error decoding b64: %s", err)
	}

	sdkTx, err := c.encodingConfig.TxConfig.TxDecoder()(protoEncoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding proto tx: %s", err)
	}

	return sdkTx, nil
}

func (c *CosmosService) queryTxs(query string) (*TendermintTxSearchResponse, error) {
	txsResponse := TendermintTxSearchResponse{}
	params := map[string][]string{"query": {fmt.Sprintf("\"%s\"", query)}}
	if err := doGet(c.tendermintRPCClient, "/tx_search", params, &txsResponse); err != nil {
		return nil, fmt.Errorf("error querying (%s): %s", query, err)
	}
	return &txsResponse, nil
}

func (c *CosmosService) queryTxsBySender(address string) (*TendermintTxSearchResponse, error) {
	txsResponse, err := c.queryTxs(fmt.Sprintf("message.sender='%s'", address))
	if err != nil {
		return nil, fmt.Errorf("error reading txhistory for sender %s: %s", address, err)
	}
	return txsResponse, nil
}

func (c *CosmosService) queryTxsByRecipient(address string) (*TendermintTxSearchResponse, error) {
	txsResponse, err := c.queryTxs(fmt.Sprintf("transfer.recipient='%s'", address))
	if err != nil {
		return nil, fmt.Errorf("error reading txhistory for recipient %s: %s", address, err)
	}
	return txsResponse, nil
}

func (c *CosmosService) queryTxsByID(txid string) (*TendermintTxSearchResponse, error) {
	txsResponse, err := c.queryTxs(fmt.Sprintf("tx.hash='%s'", txid))
	if err != nil {
		return nil, fmt.Errorf("error reading tx for id %s: %s", txid, err)
	}
	return txsResponse, nil
}

func (c *CosmosService) GetTxByID(txid string) (*HistTx, error) {
	txsResponse, err := c.queryTxsByID(txid)
	if err != nil {
		log.Errorf("error reading tx for id %s: %s", txid, err)
	}
	if len(txsResponse.Result.Txs) < 1 {
		return nil, fmt.Errorf("no tx found for id %s", txid)
	}

	return c.normalizeTx(txsResponse.Result.Txs[0])
}

func (c *CosmosService) normalizeTx(sourceTx TendermintHistTx) (*HistTx, error) {
	height := int64(-1)
	height, err := strconv.ParseInt(sourceTx.Height, 10, 64)
	if err != nil {
		log.Warnf("tx %s: unable to parse block height to int64: %s", sourceTx.Hash, err)
		height = -1
	}

	status := "pending"
	confirms := int64(0)
	if height > 0 {
		status = "confirmed"
		confirms = c.GetCurrentHeight() - height
	}

	histTx := HistTx{
		TxID:          sourceTx.Hash,
		Status:        status,
		BlockHeight:   height,
		Index:         sourceTx.Index,
		Confirmations: confirms,
		Actions:       make([]TxAction, 0, 10),
		GasWanted:     sourceTx.TxResult.GasWanted,
		GasUsed:       sourceTx.TxResult.GasUsed,
		Timestamp:     time.Now().UnixMilli(), // TODO blockTime?
	}

	sdkTx, err := c.decodeTx(sourceTx.Tx)
	if err != nil {
		log.Errorf("error decoding tx %s: %s", sourceTx.Hash, err)
	} else {
		// typedTx, ok := sdkTx.(*txtypes.Tx)
		builder, err := c.encodingConfig.TxConfig.WrapTxBuilder(sdkTx)
		if err != nil {
			return nil, fmt.Errorf("error wrapping it: %s", err)
		}
		stx := builder.GetTx()
		histTx.Memo = stx.GetMemo()

		fee := stx.GetFee()
		amt, denom, err := ParseAmountDenom(fee.String())
		if err == nil {
			histTx.Fee = amt
			histTx.FeeAsset = denom // this is uatom cosmos? atom? caip?
		} else {
			log.Errorf("tx %s error parsing amount %s: %s", sourceTx.Hash, fee.String(), err)
		}

		msgs := sdkTx.GetMsgs()
		log.Debugf("tx %s has %d msgs", sourceTx.Hash, len(msgs))
		for _, msg := range msgs {
			histTx.applyMsg(msg)
		}
	}

	if err = histTx.applyLogs(sourceTx.TxResult.Log, c.chainConfig); err != nil {
		log.Errorf("error applying logs for tx %s: %s", sourceTx.Hash, err)
	}
	return &histTx, nil
}

func (c *CosmosService) GetTxHistory(address string) (*TxHistory, error) {
	const unit = "GetTxHistory:"
	var (
		err error
	)
	txHistory := TxHistory{
		Pubkey: address,
		Txs:    make([]HistTx, 0, 10),
	}

	allTxs := make([]TendermintHistTx, 0, 100)
	// find txs this address created
	txsResponse, err := c.queryTxsBySender(address)
	if err != nil {
		log.Errorf("error reading txhistory for sender %s: %s", address, err)
	} else {
		allTxs = append(allTxs, txsResponse.Result.Txs...)
	}

	// find txs this address was a recipient of in some way
	txsResponse, err = c.queryTxsByRecipient(address)
	if err != nil {
		log.Errorf("error reading txhistory for recipient %s: %s", address, err)
	} else {
		allTxs = append(allTxs, txsResponse.Result.Txs...)
	}

	// de-duplicate txs across sender and recipient queries
	seen := make(map[string]interface{}, len(allTxs))
	uniqueTxs := make([]TendermintHistTx, 0, len(allTxs))
	for _, tx := range allTxs {
		if _, ok := seen[tx.Hash]; !ok {
			uniqueTxs = append(uniqueTxs, tx)
			seen[tx.Hash] = nil
		}
	}

	log.Debugf("have %d total txs for %s", len(uniqueTxs), address)
	for _, sourceTx := range uniqueTxs {
		histTx, err := c.normalizeTx(sourceTx)
		if err != nil {
			log.Errorf("error normalizing tx %s: %s", sourceTx.Hash, err)
			continue
		}

		txHistory.Txs = append(txHistory.Txs, *histTx)
	}
	// sort descending by block height
	sort.Slice(txHistory.Txs, func(i, j int) bool {
		if txHistory.Txs[i].BlockHeight == txHistory.Txs[j].BlockHeight {
			return txHistory.Txs[i].Index > txHistory.Txs[j].Index
		}
		return txHistory.Txs[i].BlockHeight > txHistory.Txs[j].BlockHeight
	})

	return &txHistory, nil
}

func ToAttributeMap(attributes []TendermintAttribute) map[string]string {
	result := make(map[string]string, len(attributes))
	for _, a := range attributes {
		result[a.Key] = a.Value
	}
	return result
}

// splits a cosmos denom such as 1234567uatom into amount and denom
func ParseAmountDenom(cosmosAmt string) (amount string, denom string, err error) {
	var (
		i     int
		c     rune
		found bool
	)
	for i = 0; i < len(cosmosAmt) && !found; i++ {
		c = rune(cosmosAmt[i])
		if c < '0' || c > '9' {
			found = true
		}
	}
	if !found {
		i++
	}

	amount = cosmosAmt[:i-1]
	denom = cosmosAmt[i-1:]
	return
}
