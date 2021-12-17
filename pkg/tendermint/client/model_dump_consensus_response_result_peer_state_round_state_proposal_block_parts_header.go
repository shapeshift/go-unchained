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

// DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader struct for DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader
type DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader struct {
	Total int32 `json:"total"`
	Hash string `json:"hash"`
}

// NewDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader instantiates a new DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader(total int32, hash string) *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader {
	this := DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader{}
	this.Total = total
	this.Hash = hash
	return &this
}

// NewDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeaderWithDefaults instantiates a new DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeaderWithDefaults() *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader {
	this := DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader{}
	return &this
}

// GetTotal returns the Total field value
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) GetTotal() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Total
}

// GetTotalOk returns a tuple with the Total field value
// and a boolean to check if the value has been set.
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) GetTotalOk() (*int32, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Total, true
}

// SetTotal sets field value
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) SetTotal(v int32) {
	o.Total = v
}

// GetHash returns the Hash field value
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) GetHash() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hash
}

// GetHashOk returns a tuple with the Hash field value
// and a boolean to check if the value has been set.
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) GetHashOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Hash, true
}

// SetHash sets field value
func (o *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) SetHash(v string) {
	o.Hash = v
}

func (o DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["total"] = o.Total
	}
	if true {
		toSerialize["hash"] = o.Hash
	}
	return json.Marshal(toSerialize)
}

type NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader struct {
	value *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader
	isSet bool
}

func (v NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) Get() *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader {
	return v.value
}

func (v *NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) Set(val *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) {
	v.value = val
	v.isSet = true
}

func (v NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) IsSet() bool {
	return v.isSet
}

func (v *NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader(val *DumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) *NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader {
	return &NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader{value: val, isSet: true}
}

func (v NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDumpConsensusResponseResultPeerStateRoundStateProposalBlockPartsHeader) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


