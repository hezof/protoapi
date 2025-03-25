module github.com/hezof/protoapi

go 1.21

require (
	github.com/hezof/log v0.0.0
	github.com/hezof/core v0.0.0
	github.com/hezof/protojson v0.0.0
	google.golang.org/grpc v1.67.3
	google.golang.org/protobuf v1.36.5
)


replace (
	github.com/hezof/log v0.0.0 => ../log
	github.com/hezof/core v0.0.0 => ../core
	github.com/hezof/protojson v0.0.0 => ../protojson
)
