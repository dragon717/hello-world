module github.com/hello-world/logic-server

go 1.22.0

require (
	github.com/gorilla/websocket v1.5.1
	github.com/hello-world v0.0.0
	github.com/redis/go-redis/v9 v9.7.3
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

replace github.com/hello-world => ../

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
