postgres:
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root  -e POSTGRES_PASSWORD=aaa -d postgres:13

createdb:
	docker exec -it postgres13 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres13 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:aaa@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:aaa@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:aaa@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:aaa@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen --destination db/mock/store.go --package=mockdb --build_flags=--mod=mod  github.com/martikan/simplebank/db/sqlc Store

helm-test:
	helm upgrade --install --dry-run --debug -o yaml api helm/simple-bank -n simplebank

helm-deploy:
	helm upgrade --install api helm/simple-bank -n simplebank

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock mock-mac helm-test helm-deploy