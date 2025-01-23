package protoapi

import (
	"github.com/hezof/protoapi/internal/websocket"
	"ksogit.kingsoft.net/kgo/log"
)

// HandleFunc 处理逻辑函数
type HandleFunc func(ctx *Context)

// Handler 处理句柄.
type Handler struct {
	Setting       *MethodSetting // protobuf的method的handler
	Method        string         // http请求方法
	Path          string         // http请求路径
	Status        uint32         // http请求状态
	Result        Http_Result    // http请求结果
	HandleChain   []HandleFunc   // 处理链表, 最长不超过HandleChainCapacity设置
	BodyMaxBytes  int64          // http请求体最大字节数. 如果设置则用http.MaxBytesReader()限制读入字节数! 如果是Websocket则自动忽略此参数
	FormMaxMemory int64          // http表单内存部分最大字节数. 默认 32 << 20
}

// RestfulHandleFunc 使用call生成restful的HandleFunc
func RestfulHandleFunc(ctx *Context) {

	// 业务处理(前置校验).
	rsp, err := ctx.Handler.Setting.Call(ctx, ctx.Request.Body)

	// 使用DownFile()/WriterStream()的Service Method实现必须确保返回rsp为nil(即无法用于grpc调用)!
	if err != nil {
		if xrr := ctx.WriteErrorResult(StatusErrorFrom(err)); xrr != nil {
			log.Error("write error result %v %v: %+v", ctx.Request.Method, ctx.Request.RequestURI, xrr)
		}
	} else if ctx.ResponseWriter.statusCode == 0 {
		if xrr := ctx.WriteApplyResult(rsp); xrr != nil {
			log.Error("write apply result %v %v: %+v", ctx.Request.Method, ctx.Request.RequestURI, xrr)
		}
	}
}

/*
WebsocketHandleFunc 使用call生成websocket的HandleFunc
*/
func WebsocketHandleFunc(ctx *Context) {

	// graceful关闭期间断开keepalive连接
	if ctx.mux.closed != 0 {
		ctx.ResponseWriter.Header()["Connection"] = closeConnection
	}

	// websocket必须用回原生的ResponseWriter, 因为它实现的http.Hijack接口
	conn, err := ctx.mux.upgrader.Upgrade(ctx.ResponseWriter.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Error("websocket upgrade request error: %v", err)
		return // 马上结束当前ws链接
	}
	defer conn.Close()

	for {
		// 获取读入流
		_, in, err := conn.NextReader()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				log.Error("websocket next reader error: %v", err)
			}
			return // 结束当前ws链接
		}
		// 获取写出流
		out, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			if cer, ok := err.(*websocket.CloseError); !ok {
				log.Error("websocket next writer error: %v", cer)
			}
			return // 结束当前ws链接
		}
		// 业务逻辑(前置校验)
		rsp, err := ctx.Handler.Setting.Call(ctx, in)

		// 结果处理
		if err != nil {
			if xrr := ctx.writeErrorResult(out, StatusErrorFrom(err)); xrr != nil {
				if _, ok := xrr.(*websocket.CloseError); !ok {
					log.Error("websocket write result error: %v", xrr)
				}
				return // 结束当前ws链接
			}
		} else if ctx.ResponseWriter.statusCode == 0 {
			if xrr := ctx.writeApplyResult(out, rsp); xrr != nil {
				if _, ok := xrr.(*websocket.CloseError); !ok {
					log.Error("websocket write result error: %v", xrr)
				}
				return // 结束当前ws链接
			}
		}

		// 务必关闭输出刷新缓存
		if err = out.Close(); err != nil {
			if xrr, ok := err.(*websocket.CloseError); !ok {
				log.Error("websocket write result error: %v", xrr)
			}
			return // 结束当前ws链接
		}
	}
}
