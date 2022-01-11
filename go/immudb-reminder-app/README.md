# immudb-reminder-app
A simple reminder console app that stores all data in immudb

## build
```bash
go mod tidy
go build
```

or with flags

```bash
Usage of ./immudb-reminder-app:
  -addr string
        IP address of immudb server (default "localhost")
  -db string
        Name of the database to use (default "defaultdb")
  -pass string
        Password for authenticating to immudb (default "immudb")
  -port string
        Port number of immudb server (default "3322")
  -user string
        Username for authenticating to immudb (default "immudb")
```

## run

```bash
./immudb-reminder-app
```

## requirements
immudb needs to be accessible https://github.com/codenotary/immudb

## todo
- ~~currently only connecting to localhost immudb, add parameter~~
- improve messaging
- improve output format
- improve CLI menu
