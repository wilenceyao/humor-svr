.PHONY: api

api:
	GOOS=linux GOARCH=amd64 go build -o humor-svr cmd/humor_svr.go

protoc:
	protoc --go_out=. --go_opt=paths=source_relative api/common/common.proto
	protoc --go_out=. --go_opt=paths=source_relative api/svr/rest/api.proto

lint:
	golangci-lint run
