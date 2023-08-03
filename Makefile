server:
	go run cmd/main.go

db_init:
	go run cmd/migrations/main.go db init

db_reset:
	go run cmd/migrations/main.go db reset

proto:
	protoc --go_out=. \
		--go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        pkg/pb/*.proto
