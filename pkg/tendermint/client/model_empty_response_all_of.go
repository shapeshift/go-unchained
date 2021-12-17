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

// EmptyResponseAllOf struct for EmptyResponseAllOf
type EmptyResponseAllOf struct {
	Result *map[string]interface{} `json:"result,omitempty"`
}

// NewEmptyResponseAllOf instantiates a new EmptyResponseAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEmptyResponseAllOf() *EmptyResponseAllOf {
	this := EmptyResponseAllOf{}
	return &this
}

// NewEmptyResponseAllOfWithDefaults instantiates a new EmptyResponseAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEmptyResponseAllOfWithDefaults() *EmptyResponseAllOf {
	this := EmptyResponseAllOf{}
	return &this
}

// GetResult returns the Result field value if set, zero value otherwise.
func (o *EmptyResponseAllOf) GetResult() map[string]interface{} {
	if o == nil || o.Result == nil {
		var ret map[string]interface{}
		return ret
	}
	return *o.Result
}

// GetResultOk returns a tuple with the Result field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EmptyResponseAllOf) GetResultOk() (*map[string]interface{}, bool) {
	if o == nil || o.Result == nil {
		return nil, false
	}
	return o.Result, true
}

// HasResult returns a boolean if a field has been set.
func (o *EmptyResponseAllOf) HasResult() bool {
	if o != nil && o.Result != nil {
		return true
	}

	return false
}

// SetResult gets a reference to the given map[string]interface{} and assigns it to the Result field.
func (o *EmptyResponseAllOf) SetResult(v map[string]interface{}) {
	o.Result = &v
}

func (o EmptyResponseAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Result != nil {
		toSerialize["result"] = o.Result
	}
	return json.Marshal(toSerialize)
}

type NullableEmptyResponseAllOf struct {
	value *EmptyResponseAllOf
	isSet bool
}

func (v NullableEmptyResponseAllOf) Get() *EmptyResponseAllOf {
	return v.value
}

func (v *NullableEmptyResponseAllOf) Set(val *EmptyResponseAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableEmptyResponseAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableEmptyResponseAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEmptyResponseAllOf(val *EmptyResponseAllOf) *NullableEmptyResponseAllOf {
	return &NullableEmptyResponseAllOf{value: val, isSet: true}
}

func (v NullableEmptyResponseAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEmptyResponseAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


