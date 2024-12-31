package protoapi

/*
默认code错误码. 兼容http/grpc的code范围是[0,math.MaxInt32]
统一约定:
- [0,99]              表示保留错误码! (业务/扩展切勿占用)
- [100,999]           表示请求错误码! (与http status code一致)
- [1000,9999]         表示系统错误码!
- [10000,2147483647]  表示业务错误码!
*/

type StatusResult struct {
	Status  int32         `json:"-"`                 // 状态代码(http).
	Code    int32         `json:"code"`              // 错误代码. 0表示成功
	Name    string        `json:"name,omitempty"`    // 错误名称. OK表示成功
	Message string        `json:"message,omitempty"` // 错误消息.
	Details []interface{} `json:"-"`                 // 错误参数.
}
