module github.com/hezof/protoapi/demo

go 1.23.0

toolchain go1.23.6

require (
	github.com/hezof/log v0.0.4
	github.com/hezof/protoapi v1.0.0-beta
	golang.org/x/net v0.36.0
	google.golang.org/grpc v1.67.3
	google.golang.org/protobuf v1.36.5
)

replace (
	github.com/hezof/core => ./../../base
	github.com/hezof/protoapi => ../../protoapi
)
