package protoapi

import (
	"time"
)

type Config struct {
	Name              string        `json:"name"`                 // service name
	GrpcAddr          string        `json:"grpc_addr"`            // grpc server listening address
	GrpcKeepAlive     time.Duration `json:"grpc_keep_alive"`      // grpc server keep alive
	GrpcI18nError     bool          `json:"grpc_i18n_error"`      // grpc i18n error
	HttpAddr          string        `json:"http_addr"`            // http server listening address
	HttpKeepAlive     time.Duration `json:"http_keep_alive"`      // http server keep alive
	HttpCertFile      string        `json:"http_cert_file"`       // http server TLS cert file
	HttpKeyFile       string        `json:"http_key_file"`        // http server TLS key file
	HttpBodyMaxBytes  int64         `json:"http_body_max_bytes"`  // http body最多字节, 负或0表示默认(不限制),正表示具体字节数
	HttpFormMaxMemory int64         `json:"http_form_max_memory"` // http file最多内存,超出部分存在硬盘. 负或0表示默认(32M), 正表示具体字节数
	WbskReadBuffer    int           `json:"wbsk_read_buffer"`     // websocket read buffer size
	WbskWriteBuffer   int           `json:"wbsk_write_buffer"`    // websocket write buffer size
	WbskOriginDisable bool          `json:"wbsk_origin_disable"`  // websocket would not check origin
}
