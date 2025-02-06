package protoapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// 默认值与go/pkg/http相同
const (
	defaultDialerTimeout       = 20 * time.Second
	defaultDialerKeepAlive     = 20 * time.Second
	defaultIdleConnTimeout     = 20 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
	defaultMaxIdleConnsPerHost = 64
	defaultMaxConnsPerHost     = 2048
	defaultWriteBufferSize     = 512 * 1024
	defaultReadBufferSize      = 512 * 1024
)

// JsonRpcConfig http配置
type JsonRpcConfig struct {

	// DialerTimeout 连接超时(默认3分钟)
	DialerTimeout time.Duration `json:"dialer_timeout"`

	// DialerKeepAlive 保持活跃(默认30秒)
	DialerKeepAlive time.Duration `json:"dialer_keep_alive"`

	// TLSHandshakeTimeout TLS握手超时(默认10秒)
	TLSHandshakeTimeout time.Duration

	// MaxIdleConnsPerHost 每个Host最大空闲连接数(默认64)
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host"`

	// MaxConnsPerHost 每个Host最大连接数(默认2048)
	MaxConnsPerHost int `json:"max_conns_per_host"`

	// IdleConnTimeout 空闲连接超时(默认30分)
	IdleConnTimeout time.Duration `json:"idle_conn_timeout"`

	// WriteBufferSize 写缓存区大小(默认64K)
	WriteBufferSize int `json:"write_buffer_size"`

	// ReadBufferSize 读缓存区大小(默认64K)
	ReadBufferSize int `json:"read_buffer_size"`

	// InsecureSkipVerify TLS是否跳过校验(默认false)
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
}

// JsonRpcHeader http头部
type JsonRpcHeader interface {
	InitHeader(furl string, body []byte, header http.Header)
}

// JsonRpcClient rpc客户端
type JsonRpcClient struct {
	endpoint string
	header   JsonRpcHeader
	client   *http.Client
}

// Call 远程调用. 可以指定期望的status
func (c *JsonRpcClient) Call(method string, uri string, req any, rsp any, status ...int) error {
	furl := c.endpoint + uri
	var body []byte
	if req != nil {
		enc := GetEncoder(nil)
		defer PutEncoder(enc)

		EncodeAny(enc, req)
		err := enc.Close()
		if err != nil {
			return err
		}
		body = enc.buff
	}
	hreq, err := http.NewRequest(method, furl, content(body))
	if err != nil {
		return err
	}
	if c.header != nil {
		c.header.InitHeader(furl, body, hreq.Header)
	}
	hreq.Header.Set("Content-Type", "application/json")
	hreq.ContentLength = int64(len(body))

	hrsp, err := c.client.Do(hreq)
	if err != nil {
		return err
	}
	defer discard(hrsp.Body)

	if len(status) > 0 {
		if !contains(hrsp.StatusCode, status) {
			return fmt.Errorf("invalid status code: %v, expected %v", hrsp.StatusCode, status)
		}
	}
	if rsp != nil {
		dec := GetDecoder(hrsp.Body)
		defer PutDecoder(dec)

		DecodeAny(dec, rsp)
		return dec.Close()
	}
	return nil
}

// NewJsonRpcClient 创建rpc客户端
func NewJsonRpcClient(e string, h JsonRpcHeader, c *JsonRpcConfig) *JsonRpcClient {
	return &JsonRpcClient{
		endpoint: e,
		header:   h,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   NvlD(c.DialerTimeout, defaultDialerTimeout),
					KeepAlive: NvlD(c.DialerKeepAlive, defaultDialerKeepAlive),
				}).DialContext,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
				TLSHandshakeTimeout: NvlD(c.TLSHandshakeTimeout, defaultTLSHandshakeTimeout),
				MaxIdleConnsPerHost: NvlI(c.MaxIdleConnsPerHost, defaultMaxIdleConnsPerHost),
				MaxConnsPerHost:     NvlI(c.MaxConnsPerHost, defaultMaxConnsPerHost),
				IdleConnTimeout:     NvlD(c.IdleConnTimeout, defaultIdleConnTimeout),
				WriteBufferSize:     NvlI(c.WriteBufferSize, defaultWriteBufferSize),
				ReadBufferSize:      NvlI(c.ReadBufferSize, defaultReadBufferSize),
				DisableKeepAlives:   true, // 尝试解决UnexpectedEOF
			},
		},
	}
}

func contains(p int, vs []int) bool {
	for _, v := range vs {
		if v == p {
			return true
		}
	}
	return false
}

var buff = make([]byte, 32*1024)

func discard(in io.ReadCloser) {
	_, _ = io.CopyBuffer(io.Discard, in, buff)
	_ = in.Close()
}

func content(data []byte) io.Reader {
	if len(data) == 0 {
		return http.NoBody
	}
	return bytes.NewReader(data)
}
