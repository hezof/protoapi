package protoapi

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	GRACE_FLAG = "__GRACE__" // 用于优雅重启传递环境变量
	GRACE_GRPC = "__GRPC__"
	GRACE_HTTP = "__HTTP__"
	GRACE_BOTH = "__BOTH__"
	GRACE_NONE = ""
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type keepAliveTCPListener struct {
	*net.TCPListener
	KeepAlivePeriod time.Duration
}

func (ln keepAliveTCPListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if ln.KeepAlivePeriod == 0 {
		_ = tc.SetKeepAlive(false)
		_ = tc.SetKeepAlivePeriod(3 * time.Minute)
	} else {
		_ = tc.SetKeepAlive(true)
		_ = tc.SetKeepAlivePeriod(ln.KeepAlivePeriod)
	}
	return tc, nil
}

func getListenerFile(l net.Listener) *os.File {
	var file *os.File
	switch l := l.(type) {
	case *keepAliveTCPListener:
		file, _ = l.TCPListener.File()
	case *net.TCPListener:
		file, _ = l.File()
	}
	return file
}

func graceGrpcListener(addr string) (net.Listener, error) {
	var flag = os.Getenv(GRACE_FLAG)
	if flag != "" {
		var fd uintptr
		switch flag {
		case GRACE_GRPC:
			fd = 3
		case GRACE_BOTH:
			fd = 3
		default:
			return nil, nil
		}
		var file = os.NewFile(fd, "")
		defer file.Close()

		return net.FileListener(file)
	}
	return net.Listen("tcp", addr)
}

func graceHttpListener(addr string, keepalive time.Duration) (net.Listener, error) {
	var flag = os.Getenv(GRACE_FLAG)
	if flag != "" {
		var fd uintptr
		switch flag {
		case GRACE_HTTP:
			fd = 3
		case GRACE_BOTH:
			fd = 4
		default:
			return nil, nil
		}
		var file = os.NewFile(fd, "")
		defer file.Close()

		return net.FileListener(file)
	}

	var tln, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &keepAliveTCPListener{TCPListener: tln.(*net.TCPListener), KeepAlivePeriod: keepalive}, nil
}

func graceShutdownOrRestart(grpcServer *grpc.Server, grpcListener net.Listener, httpServer *http.Server, httpListener net.Listener, closed *uint32) {
	var sch = make(chan os.Signal, 1)
	defer signal.Stop(sch)

	signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		sig := <-sch

		// HTTP服务器关闭后返回的response添加Connection:closed.防止keep-alive影响! 重启不需要!
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			atomic.StoreUint32(closed, 1)
		}

		switch sig {
		case syscall.SIGHUP:
			var (
				args  []string
				flag  string
				files []*os.File
			)
			// 设置重启标志及参数
			if len(os.Args) > 1 {
				args = os.Args[1:]
			}
			if grpcListener != nil && httpListener != nil {
				flag = GRACE_BOTH
				files = []*os.File{getListenerFile(grpcListener), getListenerFile(httpListener)}
			} else if grpcListener != nil {
				flag = GRACE_GRPC
				files = []*os.File{getListenerFile(grpcListener)}
			} else if httpListener != nil {
				flag = GRACE_HTTP
				files = []*os.File{getListenerFile(httpListener)}
			} else {
				flag = GRACE_NONE
			}

			// 执行重启命令
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), GRACE_FLAG+"="+flag) // 拼加标志
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.ExtraFiles = files
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stdout, "graceful restart server error: %v", err)
			}
			fallthrough // 启动子进程后进入grace shutdown
		case syscall.SIGINT, syscall.SIGTERM:
			wg := new(sync.WaitGroup)
			if httpServer != nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					httpServer.Shutdown(context.Background())
				}()
			}
			if grpcServer != nil {
				wg.Add(1)
				go func() {
					defer wg.Done()
					grpcServer.GracefulStop()
				}()
			}
			wg.Wait()
			return
		}
	}
}
