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

// BlockSearchResponse struct for BlockSearchResponse
type BlockSearchResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id int32 `json:"id"`
	Result BlockSearchResponseResult `json:"result"`
}

// NewBlockSearchResponse instantiates a new BlockSearchResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBlockSearchResponse(jsonrpc string, id int32, result BlockSearchResponseResult) *BlockSearchResponse {
	this := BlockSearchResponse{}
	this.Jsonrpc = jsonrpc
	this.Id = id
	this.Result = result
	return &this
}

// NewBlockSearchResponseWithDefaults instantiates a new BlockSearchResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBlockSearchResponseWithDefaults() *BlockSearchResponse {
	this := BlockSearchResponse{}
	return &this
}

// GetJsonrpc returns the Jsonrpc field value
func (o *BlockSearchResponse) GetJsonrpc() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Jsonrpc
}

// GetJsonrpcOk returns a tuple with the Jsonrpc field value
// and a boolean to check if the value has been set.
func (o *BlockSearchResponse) GetJsonrpcOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Jsonrpc, true
}

// SetJsonrpc sets field value
func (o *BlockSearchResponse) SetJsonrpc(v string) {
	o.Jsonrpc = v
}

// GetId returns the Id field value
func (o *BlockSearchResponse) GetId() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *BlockSearchResponse) GetIdOk() (*int32, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *BlockSearchResponse) SetId(v int32) {
	o.Id = v
}

// GetResult returns the Result field value
func (o *BlockSearchResponse) GetResult() BlockSearchResponseResult {
	if o == nil {
		var ret BlockSearchResponseResult
		return ret
	}

	return o.Result
}

// GetResultOk returns a tuple with the Result field value
// and a boolean to check if the value has been set.
func (o *BlockSearchResponse) GetResultOk() (*BlockSearchResponseResult, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Result, true
}

// SetResult sets field value
func (o *BlockSearchResponse) SetResult(v BlockSearchResponseResult) {
	o.Result = v
}

func (o BlockSearchResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["jsonrpc"] = o.Jsonrpc
	}
	if true {
		toSerialize["id"] = o.Id
	}
	if true {
		toSerialize["result"] = o.Result
	}
	return json.Marshal(toSerialize)
}

type NullableBlockSearchResponse struct {
	value *BlockSearchResponse
	isSet bool
}

func (v NullableBlockSearchResponse) Get() *BlockSearchResponse {
	return v.value
}

func (v *NullableBlockSearchResponse) Set(val *BlockSearchResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableBlockSearchResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableBlockSearchResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBlockSearchResponse(val *BlockSearchResponse) *NullableBlockSearchResponse {
	return &NullableBlockSearchResponse{value: val, isSet: true}
}

func (v NullableBlockSearchResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBlockSearchResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


