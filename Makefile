include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

local-migration-create:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} create init sql

generate:
	make generate-auth-api

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path=api/auth_v1 --proto_path vendor.protogen\
    	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
    	--plugin=protoc-gen-go=./bin/protoc-gen-go \
    	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
    	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
    	--validate_out lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
    	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
    	api/auth_v1/auth.proto

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi

build:
	GOOS=linux GOARCH=amd64 go build -o service_auth cmd/server/main.go

test:
	go clean -testcache
	go test github.com/semho/chat-microservices/auth/internal/service/... \
			github.com/semho/chat-microservices/auth/internal/api/... -covermode count -count 5


test-coverage:
	go clean -testcache
	go test github.com/semho/chat-microservices/auth/internal/service/... \
            github.com/semho/chat-microservices/auth/internal/api/... -covermode count -coverprofile=coverage.tmp.out -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore