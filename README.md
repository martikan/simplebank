[![codecov](https://codecov.io/gh/martikan/simplebank/branch/main/graph/badge.svg?token=50BOAZ3DUQ)](https://codecov.io/gh/martikan/simplebank)

# Demo api for a simple app

## Mocking
*Firstly add* ```~/go/bin``` *to your PATH variable to be able to use mockgen*

## Add new migration version:
Create add_users_table scripts for migration.
```bash
migrate create -ext sql -dir db/migration -seq add_users_table
```