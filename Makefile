proto:
	protoc pkg/pb/*.proto --go-grpc_out=.
server:
	go run cmd/main.go