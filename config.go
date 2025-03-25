package protoapi

import (
	"time"
)

type Config struct {
	Name                string        // component Name
	GrpcAddr            string        // grpc server listening address
	GrpcKeepAlive       time.Duration // grpc server keep alive
	GrpcI18nError       bool          // grpc i18n error
	GrpcKeepAlivePolicy time.Duration // grpc keep alive Enforcement policy min time
	HttpAddr            string        // http server listening address
	HttpKeepAlive       time.Duration // http server keep alive
	HttpCertFile        string        // http server TLS cert file
	HttpKeyFile         string        // http server TLS key file
	HttpBodyMaxBytes    int64         // http body最多字节, 负或0表示默认(不限制),正表示具体字节数
	HttpFormMaxMemory   int64         // http file最多内存,超出部分存在硬盘. 负或0表示默认(32M), 正表示具体字节数
	HttpReadTimeout     time.Duration // http ReadTimeout is the maximum duration for reading the entire request, including the body
	HttpWriteTimeout    time.Duration // http WriteTimeout is the maximum duration before timing out writes of the response
	HttpIdleTimeout     time.Duration // http IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled.
	WbskReadBuffer      int           // websocket read buffer size
	WbskWriteBuffer     int           // websocket write buffer size
	WbskOriginDisable   bool          // websocket would not check origin
}

// mergeConfig 合并默认配置
func mergeConfig(c *Config) *Config {
	if c == nil {
		c = new(Config)
	}
	if c.GrpcKeepAlive == 0 {
		c.GrpcKeepAlive = 5 * time.Minute
	}
	if c.GrpcKeepAlivePolicy == 0 {
		c.GrpcKeepAlivePolicy = 5 * time.Minute
	}
	if c.HttpKeepAlive == 0 {
		c.HttpKeepAlive = 3 * time.Minute
	}
	if c.HttpBodyMaxBytes == 0 {
		c.HttpBodyMaxBytes = DefMaxMem
	}
	if c.HttpFormMaxMemory == 0 {
		c.HttpFormMaxMemory = DefMaxMem
	}
	return c
}
