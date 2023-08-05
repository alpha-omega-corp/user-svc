server:
	go run cmd/main.go

proto:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        pkg/proto/*.proto

db_create:
	docker-compose up -d

db_init:
	go run cmd/migrations/main.go db init

db_reset:
	go run cmd/migrations/main.go db reset
