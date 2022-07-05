[![codecov](https://codecov.io/gh/martikan/simplebank/branch/main/graph/badge.svg?token=50BOAZ3DUQ)](https://codecov.io/gh/martikan/simplebank)

# Demo api for a simple app

## Mocking
*Firstly add* ```~/go/bin``` *to your PATH variable to be able to use mockgen*

## Add new migration version:
Create add_users_table scripts for migration.
```bash
migrate create -ext sql -dir db/migration -seq add_users_table
```

## Build docker image locally and run it:
Build the docker image.
```bash
docker build -t simplebank:latest .
```
Run it.
```bash
docker run --name simplebank --network bank-network -p 8085:8085 -e GIN_MODE=release -e DATASOURCE_URL="postgres://root:aaa@postgres13:5432/simple_bank?sslmode=disable" -d simplebank:latest
```