package protoapi

import "github.com/hezof/base"

func NormalResult[V any](v *V) *StatusResultModel {
	if v == nil {
		panic("NormalResult: unmarshal nil")
	}
	return &StatusResultModel{
		Data: v,
	}
}

func UnwrapResult[V any](v *V) *V {
	if v == nil {
		panic("UnwrapResult: unmarshal nil")
	}
	return v
}

func EventsResult[V any](v *V) *V {
	if v == nil {
		panic("EventsResult: unmarshal nil")
	}
	return v
}

func NewJsonRpcClient(endpoint string, config *base.HttpConfig, header base.HttpHeader) *base.JsonRpcClient {
	return base.NewJsonRpcClient(endpoint, config, header, nil, nil)
}
