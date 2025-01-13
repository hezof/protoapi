package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CustomOptions 用户选项
type CustomOptions struct {
	Help      bool   // 打印帮助
	Debug     bool   // 打印调试
	Update    bool   // 更新插件
	Clean     bool   // 清理文件(*.pb.go, *_grpc.pb.go, *_protoapi.pb.go, *_protoapi.json)
	Config    string // 配置变量, 例如: "VERSION=0.5.1;GOPROXY=https://goproxy.cn;GOPRIVATE=*.net,*.cn"
	GoOut     string // GO输出目录
	ProtoBase string // PB基准目录
	ProtoPath string // PB查找目录,多值逗号分隔
	GrpcV2    bool   // 生成GRPCv2代码[require_unimplemented_servers=true]
}

// SystemOptions 系统选项
type SystemOptions struct {
	HomeDir                 string // .protogen目录
	TempDir                 string // .protogen/tmp目录
	IncludeDir              string // .protogen/include目录
	ProtocFile              string // .protogen/protoc文件
	ProtocGenGoFile         string // .protogen/protoc-gen-go文件
	ProtocGenGoGrpcFile     string // .protogen/protoc-gen-go-grpc文件
	ProtocGenGoProtoapiFile string // .protogen/protoc-gen-go-protoapi文件
	GoModFile               string // .protogen/go.mod文件
	GoSumFile               string // .protogen/go.sum文件

	GO111MODULE string
	GOSUMDB     string
	GOEXE       string
	GOBIN       string
	GOMODCACHE  string
	GOCACHE     string
	GOTMPDIR    string
}

func initCustomOptions(ops *Context) {
	ops.flagset.BoolVar(&ops.Help, `help`, false, `打印帮助`)
	ops.flagset.BoolVar(&ops.Debug, `debug`, false, `打印调试`)
	ops.flagset.BoolVar(&ops.Update, `update`, false, `更新插件`)
	ops.flagset.BoolVar(&ops.Clean, `clean`, false, `清理文件[*.pb.go, *_grpc.pb.go, *_protoapi.pb.go, *_protoapi.json]`)
	ops.flagset.StringVar(&ops.Config, `config`, ``, fmt.Sprintf(`配置变量.默认"VERSION=%v;GOPROXY=%v;GOPRIVATE=%v;MAVEN_CENTRAL=%v"`, VERSION, `https://goproxy.cn`, `*.net,*.cn`, `https://maven.aliyun.com/repository/central`))
	ops.flagset.StringVar(&ops.GoOut, `go_out`, work(), `GO输出目录,默认当前目录`)
	ops.flagset.StringVar(&ops.ProtoBase, `proto_base`, work(), `PB基准目录,默认当前目录`)
	ops.flagset.StringVar(&ops.ProtoPath, `proto_path`, ``, `PB查找目录[逗号分隔]`)
	ops.flagset.BoolVar(&ops.GrpcV2, `grpc_v2`, false, `生成GRPC代码[require_unimplemented_servers=true]`)
}

func initSystemOptions(ops *Context) {
	ops.HomeDir = home()
	ops.TempDir = filepath.Join(ops.HomeDir, `tmp`)
	ops.IncludeDir = filepath.Join(ops.HomeDir, include())
	ops.ProtocFile = filepath.Join(ops.HomeDir, protoc())
	ops.ProtocGenGoFile = filepath.Join(ops.HomeDir, protocGenGo())
	ops.ProtocGenGoGrpcFile = filepath.Join(ops.HomeDir, protocGenGoGrpc())
	ops.ProtocGenGoProtoapiFile = filepath.Join(ops.HomeDir, protocGenGoProtoapi())
	ops.GoModFile = filepath.Join(ops.HomeDir, `go.mod`)
	ops.GoSumFile = filepath.Join(ops.HomeDir, `go.sum`)

	ops.GO111MODULE = `on`
	ops.GOSUMDB = `off`
	ops.GOEXE = goexe()
	ops.GOBIN = ops.TempDir
	ops.GOMODCACHE = ops.TempDir
	ops.GOCACHE = ops.TempDir
	ops.GOTMPDIR = ops.TempDir

}

type Config struct {
	VERSION       string // protogen版本, 默认: VERSION
	GOPROXY       string // go代理仓库, 默认: https://goproxy.cn
	GOPRIVATE     string // go私有代理, 默认: *.net,*.cn
	MAVEN_CENTRAL string // maven中央仓库, 默认: https://maven.aliyun.com/repository/central
}

func parseConfig(s string) *Config {

	// 将windows环境变量替换为linux
	s = strings.ReplaceAll(s, `:`, `;`)

	c := new(Config)
	// 默认值
	c.VERSION = VERSION
	c.GOPROXY = `https://goproxy.cn`
	c.GOPRIVATE = `*.net,*.cn`
	c.MAVEN_CENTRAL = `https://maven.aliyun.com/repository/central`
	// 参数值
	for _, env := range strings.Split(s, ";") {
		kvs := strings.SplitN(strings.TrimSpace(env), "=", 3)
		if len(kvs) > 1 {
			k := strings.TrimSpace(kvs[0])
			v := strings.TrimSpace(kvs[1])
			switch k {
			case `VERSION`:
				c.VERSION = v
			case `GOPROXY`:
				c.GOPROXY = v
			case `GOPRIVATE`:
				c.GOPRIVATE = v
			case `MAVEN_CENTRAL`:
				c.MAVEN_CENTRAL = v
			}
		}
	}
	return c
}

func work() string {
	cwd, _ := os.Getwd()
	if cwd == "" {
		cwd = "./"
	}
	cwd, _ = filepath.Abs(cwd)
	return cwd
}

const PROTOGEN_HOME_ENV = `PROTOGEN_HOME`

func home() string {
	val := os.Getenv(PROTOGEN_HOME_ENV)
	if val == `` {
		loc, err := exec.LookPath(os.Args[0])
		if err != nil {
			loc = os.Args[0]
		}
		abs, err := filepath.Abs(loc)
		if err != nil {
			abs = loc
		}
		val = filepath.Join(filepath.Dir(abs), `.protogen`)
	}
	return val
}

func goexe() string {
	switch runtime.GOOS {
	case `windows`:
		return `.exe`
	default:
		return ``
	}
}

func root(path string) string {
	if path == `` {
		if cwd, _ := os.Getwd(); cwd != `` {
			path = cwd
		} else {
			path = `./`
		}
	}
	ret, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return ret
}

func protoc() string {
	for _, p := range Plugins {
		if p.Name == `protoc` {
			return filepath.Base(p.Module) + `_` + p.Version[1:] + goexe()
		}
	}
	return ``
}

func include() string {
	for _, p := range Plugins {
		if p.Name == `include` {
			return filepath.Base(p.Module) + `_` + p.Version[1:]
		}
	}
	return ``
}

func protocGenGo() string {
	for _, p := range Plugins {
		if p.Name == `protoc-gen-go` {
			return filepath.Base(p.Module) + `_` + p.Version[1:] + goexe()
		}
	}
	return ``
}

func protocGenGoGrpc() string {
	for _, p := range Plugins {
		if p.Name == `protoc-gen-go-grpc` {
			return filepath.Base(p.Module) + `_` + p.Version[1:] + goexe()
		}
	}
	return ``
}

func protocGenGoProtoapi() string {
	for _, p := range Plugins {
		if p.Name == `protoc-gen-go-protoapi` {
			return filepath.Base(p.Module) + `_` + p.Version[1:] + goexe()
		}
	}
	return ``
}
