package protoapi

import (
	"github.com/hezof/protoapi/internal/websocket"
	"io"
	"ksogit.kingsoft.net/kgo/log"
)

// Handler 处理句柄. 注意事项: handler存储在前缀树(httprouter), 必须保证所有属性在运行时只读!
type Handler struct {
	Meta          *MethodSetting // 自动创建的handler才有meta-info. 手动创建的handler无meta-info
	Method        string         // http请求方法
	Path          string         // http请求路径
	BodyMaxBytes  int64          // http请求体最大字节数. 如果设置则用http.MaxBytesReader()限制读入字节数! 如果是Websocket则自动忽略此参数
	FormMaxMemory int64          // 文件内部部分最大字节数. 默认 32 << 20
	HandleChain   []HandleFunc   // 处理链表, 最长不超过HandleChainCapacity设置
}

// RestfulHandleFunc 使用call生成restful的HandleFunc
func RestfulHandleFunc(fun Call) HandleFunc {
	return func(ctx *Context) {

		// 业务处理(前置校验).
		rsp, err := fun(ctx, ctx.Request.Body)

		// 使用DownFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
		if err != nil {
			if cer := ctx.WriteErrorResult(err); cer != nil {
				log.Error("write error result %v %v: %+v", ctx.Request.Method, ctx.Request.RequestURI, cer)
			}
		} else if ctx.ResponseWriter.statusCode == 0 {
			if cer := ctx.WriteApplyResult(rsp); cer != nil {
				log.Error("write apply result %v %v: %+v", ctx.Request.Method, ctx.Request.RequestURI, cer)
			}
		}
	}
}

/*
WebsocketHandleFunc 使用call生成websocket的HandleFunc
*/
func WebsocketHandleFunc(fun Call, upgrader *websocket.Upgrader) HandleFunc {

	return func(ctx *Context) {

		if ctx.mux.closed != 0 {
			ctx.ResponseWriter.Header()["Connection"] = closeConnection
		}

		// websocket必须用回原生的ResponseWriter, 因为它实现的http.Hijack接口
		conn, err := upgrader.Upgrade(ctx.ResponseWriter.ResponseWriter, ctx.Request, nil)
		if err != nil {
			log.Error("websocket upgrade request error: %v", err)
			return // 马上结束当前ws链接
		}
		defer conn.Close()

		var (
			in     io.Reader
			out    io.WriteCloser
			rsp    interface{}
			resMap map[uint32]*resource
		)
		for {
			// 获取读入流
			_, in, err = conn.NextReader()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); !ok {
					log.Error("websocket next reader error: %v", err)
				}
				return // 结束当前ws链接
			}
			// 获取写出流
			out, err = conn.NextWriter(websocket.TextMessage)
			if err != nil {
				if cer, ok := err.(*websocket.CloseError); !ok {
					log.Error("websocket next writer error: %v", cer)
				}
				return // 结束当前ws链接
			}
			// 业务逻辑(前置校验)
			rsp, err = fun(ctx, in)
			// 错误处理
			if err != nil {
				// 统一转换错误为StatusError
				result := StatusErrorFrom(err, profile.DefaultErrorStatus)
				// 国际化错误消息(延后初始化)
				if lenResMap > 0 {
					if resMap == nil {
						resMap = fastGetResMapByAcceptLanguage(ctx.getAcceptLanguage())
					}
					ctx.i18nErrorResult(result, resMap)
				}
				err = enc.Encode(result)
			} else {
				if ctx.Handler.Unwrap {
					err = enc.Encode(rsp)
				} else {
					// 重用wrapper输出,减少临时对象GC!
					ctx.wrapper.Data = rsp
					err = enc.Encode(&ctx.wrapper)
					ctx.wrapper.Data = nil // 及时清理
				}
			}
			if err != nil {
				if _, ok := err.(*websocket.CloseError); !ok {
					log.Error("websocket write result error: %v", err)
				}
				return // 结束当前ws链接
			}
			// 务必关闭输出刷新缓存
			if err = out.Close(); err != nil {
				if cer, ok := err.(*websocket.CloseError); !ok {
					log.Error("websocket write result error: %v", cer)
				}
				return // 结束当前ws链接
			}
		}
	}
}
