package protoapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/hezof/core"
	"io"
	"net"
	"net/http"
	"os"
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

	// Debug 打印调试信息
	Debug bool `json:"debug"`

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
	config   JsonRpcConfig
	endpoint string
	header   JsonRpcHeader
	client   *http.Client
}

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

func (c *JsonRpcClient) GET(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodGet, uri, req, rsp, status...)
}

func (c *JsonRpcClient) POST(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodPost, uri, req, rsp, status...)
}

func (c *JsonRpcClient) PUT(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodPut, uri, req, rsp, status...)
}

func (c *JsonRpcClient) DELETE(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodDelete, uri, req, rsp, status...)
}

func (c *JsonRpcClient) HEAD(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodHead, uri, req, rsp, status...)
}

func (c *JsonRpcClient) PATCH(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodPatch, uri, req, rsp, status...)
}

func (c *JsonRpcClient) OPTIONS(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodOptions, uri, req, rsp, status...)
}

func (c *JsonRpcClient) CONNECT(uri string, req any, rsp any, status ...int) error {
	return c.Do(http.MethodConnect, uri, req, rsp, status...)
}

// Do 远程调用. 可以指定期望的status
func (c *JsonRpcClient) Do(method string, uri string, req any, rsp any, status ...int) error {
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

	if c.config.Debug {
		fmt.Fprintln(os.Stderr, "------JsonRpcClient Request------")
		fmt.Fprintln(os.Stderr, method, hreq.URL.String())
		for k, vs := range hreq.Header {
			fmt.Fprintln(os.Stderr, k, ":", vs)
		}
		os.Stderr.Write(body)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "------JsonRpcClient Response------")
		fmt.Fprintln(os.Stderr, hrsp.Proto, hrsp.Status)
		for k, vs := range hrsp.Header {
			fmt.Fprintln(os.Stderr, k, ":", vs)
		}
		data, _ := io.ReadAll(hrsp.Body)
		hrsp.Body = &BuffBody{
			data: data,
		}
		os.Stderr.Write(data)
		fmt.Fprintln(os.Stderr)
	}

	if len(status) > 0 {
		if !contains(hrsp.StatusCode, status) {
			return fmt.Errorf("unexpected status code: %v, expected %v", hrsp.StatusCode, status)
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
	if c == nil {
		c = new(JsonRpcConfig)
	}
	return &JsonRpcClient{
		config:   *c,
		endpoint: e,
		header:   h,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   core.NvlD(c.DialerTimeout, defaultDialerTimeout),
					KeepAlive: core.NvlD(c.DialerKeepAlive, defaultDialerKeepAlive),
				}).DialContext,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify},
				TLSHandshakeTimeout: core.NvlD(c.TLSHandshakeTimeout, defaultTLSHandshakeTimeout),
				MaxIdleConnsPerHost: core.NvlI(c.MaxIdleConnsPerHost, defaultMaxIdleConnsPerHost),
				MaxConnsPerHost:     core.NvlI(c.MaxConnsPerHost, defaultMaxConnsPerHost),
				IdleConnTimeout:     core.NvlD(c.IdleConnTimeout, defaultIdleConnTimeout),
				WriteBufferSize:     core.NvlI(c.WriteBufferSize, defaultWriteBufferSize),
				ReadBufferSize:      core.NvlI(c.ReadBufferSize, defaultReadBufferSize),
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
