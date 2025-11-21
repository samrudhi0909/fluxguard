PROTO_FILE = api/limiter.proto
GO_OUT = .

.PHONY: proto clean

proto:
	protoc --go_out=$(GO_OUT) --go_opt=paths=source_relative \
	--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
	$(PROTO_FILE)

clean:
	rm -f api/*.pb.go