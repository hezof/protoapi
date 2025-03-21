module github.com/hezof/protoapi/demo

go 1.23.0

toolchain go1.23.6

require (
	github.com/hezof/core v0.0.0
	github.com/hezof/log v0.0.0
	github.com/hezof/protoapi v0.0.0
	google.golang.org/grpc v1.67.3
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/hezof/protojson v0.0.0 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
)

replace (
	github.com/hezof/core v0.0.0 => ../../core
	github.com/hezof/log v0.0.0 => ../../log
	github.com/hezof/protoapi v0.0.0 => ../../protoapi
	github.com/hezof/protojson v0.0.0 => ../../protojson
)
