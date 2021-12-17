/*
Tendermint RPC

Tendermint supports the following RPC protocols:  * URI over HTTP * JSONRPC over HTTP * JSONRPC over websockets  ## Configuration  RPC can be configured by tuning parameters under `[rpc]` table in the `$TMHOME/config/config.toml` file or by using the `--rpc.X` command-line flags.  Default rpc listen address is `tcp://0.0.0.0:26657`. To set another address, set the `laddr` config parameter to desired value. CORS (Cross-Origin Resource Sharing) can be enabled by setting `cors_allowed_origins`, `cors_allowed_methods`, `cors_allowed_headers` config parameters.  ## Arguments  Arguments which expect strings or byte arrays may be passed as quoted strings, like `\"abc\"` or as `0x`-prefixed strings, like `0x616263`.  ## URI/HTTP  A REST like interface.      curl localhost:26657/block?height=5  ## JSONRPC/HTTP  JSONRPC requests can be POST'd to the root RPC endpoint via HTTP.      curl --header \"Content-Type: application/json\" --request POST --data '{\"method\": \"block\", \"params\": [\"5\"], \"id\": 1}' localhost:26657  ## JSONRPC/websockets  JSONRPC requests can be also made via websocket. The websocket endpoint is at `/websocket`, e.g. `localhost:26657/websocket`. Asynchronous RPC functions like event `subscribe` and `unsubscribe` are only available via websockets.  Example using https://github.com/hashrocket/ws:      ws ws://localhost:26657/websocket     > { \"jsonrpc\": \"2.0\", \"method\": \"subscribe\", \"params\": [\"tm.event='NewBlock'\"], \"id\": 1 } 

API version: Master
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// UnconfirmedTransactionsResponseResult struct for UnconfirmedTransactionsResponseResult
type UnconfirmedTransactionsResponseResult struct {
	NTxs string `json:"n_txs"`
	Total string `json:"total"`
	TotalBytes string `json:"total_bytes"`
	Txs []*string `json:"txs"`
}

// NewUnconfirmedTransactionsResponseResult instantiates a new UnconfirmedTransactionsResponseResult object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUnconfirmedTransactionsResponseResult(nTxs string, total string, totalBytes string, txs []*string) *UnconfirmedTransactionsResponseResult {
	this := UnconfirmedTransactionsResponseResult{}
	this.NTxs = nTxs
	this.Total = total
	this.TotalBytes = totalBytes
	this.Txs = txs
	return &this
}

// NewUnconfirmedTransactionsResponseResultWithDefaults instantiates a new UnconfirmedTransactionsResponseResult object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUnconfirmedTransactionsResponseResultWithDefaults() *UnconfirmedTransactionsResponseResult {
	this := UnconfirmedTransactionsResponseResult{}
	return &this
}

// GetNTxs returns the NTxs field value
func (o *UnconfirmedTransactionsResponseResult) GetNTxs() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NTxs
}

// GetNTxsOk returns a tuple with the NTxs field value
// and a boolean to check if the value has been set.
func (o *UnconfirmedTransactionsResponseResult) GetNTxsOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.NTxs, true
}

// SetNTxs sets field value
func (o *UnconfirmedTransactionsResponseResult) SetNTxs(v string) {
	o.NTxs = v
}

// GetTotal returns the Total field value
func (o *UnconfirmedTransactionsResponseResult) GetTotal() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Total
}

// GetTotalOk returns a tuple with the Total field value
// and a boolean to check if the value has been set.
func (o *UnconfirmedTransactionsResponseResult) GetTotalOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Total, true
}

// SetTotal sets field value
func (o *UnconfirmedTransactionsResponseResult) SetTotal(v string) {
	o.Total = v
}

// GetTotalBytes returns the TotalBytes field value
func (o *UnconfirmedTransactionsResponseResult) GetTotalBytes() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TotalBytes
}

// GetTotalBytesOk returns a tuple with the TotalBytes field value
// and a boolean to check if the value has been set.
func (o *UnconfirmedTransactionsResponseResult) GetTotalBytesOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.TotalBytes, true
}

// SetTotalBytes sets field value
func (o *UnconfirmedTransactionsResponseResult) SetTotalBytes(v string) {
	o.TotalBytes = v
}

// GetTxs returns the Txs field value
// If the value is explicit nil, the zero value for []*string will be returned
func (o *UnconfirmedTransactionsResponseResult) GetTxs() []*string {
	if o == nil {
		var ret []*string
		return ret
	}

	return o.Txs
}

// GetTxsOk returns a tuple with the Txs field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UnconfirmedTransactionsResponseResult) GetTxsOk() (*[]*string, bool) {
	if o == nil || o.Txs == nil {
		return nil, false
	}
	return &o.Txs, true
}

// SetTxs sets field value
func (o *UnconfirmedTransactionsResponseResult) SetTxs(v []*string) {
	o.Txs = v
}

func (o UnconfirmedTransactionsResponseResult) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["n_txs"] = o.NTxs
	}
	if true {
		toSerialize["total"] = o.Total
	}
	if true {
		toSerialize["total_bytes"] = o.TotalBytes
	}
	if o.Txs != nil {
		toSerialize["txs"] = o.Txs
	}
	return json.Marshal(toSerialize)
}

type NullableUnconfirmedTransactionsResponseResult struct {
	value *UnconfirmedTransactionsResponseResult
	isSet bool
}

func (v NullableUnconfirmedTransactionsResponseResult) Get() *UnconfirmedTransactionsResponseResult {
	return v.value
}

func (v *NullableUnconfirmedTransactionsResponseResult) Set(val *UnconfirmedTransactionsResponseResult) {
	v.value = val
	v.isSet = true
}

func (v NullableUnconfirmedTransactionsResponseResult) IsSet() bool {
	return v.isSet
}

func (v *NullableUnconfirmedTransactionsResponseResult) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUnconfirmedTransactionsResponseResult(val *UnconfirmedTransactionsResponseResult) *NullableUnconfirmedTransactionsResponseResult {
	return &NullableUnconfirmedTransactionsResponseResult{value: val, isSet: true}
}

func (v NullableUnconfirmedTransactionsResponseResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUnconfirmedTransactionsResponseResult) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


