postgres:
	docker run --name some-postgres -p 4321:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14.4-alpine


createDB:
	docker exec -it some-postgres createdb --username postgres --owner postgres bank

dropDB:
	docker exec -it some-postgres dropdb --username postgres --owner postgres bank

migrateup:
	migrate -path pkg/db/migrations --database "postgresql://postgres:postgres@localhost:4321/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path pkg/db/migrations --database "postgresql://postgres:postgres@localhost:4321/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination pkg/db/mock/store.go github.com/samuelowad/bank/pkg/db/sqlc Store

.PHONY: postgres createDB dropDB migrateup migratedown sqlc server mock