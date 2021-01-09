.PHONY: api

api:
	GOOS=linux GOARCH=amd64 go build -o humor-api cmd/humor_api.go

protoc:
	protoc --go_out=. --go_opt=paths=source_relative api/common/common.proto
	protoc --go_out=. --go_opt=paths=source_relative api/rest/api.proto

lint:
	golangci-lint run
