# Pretty Log

pretty log is a cli tool that format `JSON` logs to a more readable format.

## Installation

```bash
go install github.com/tsingshaner/go-pkg/prettylog/cmd/prettylog@latest
```

## Usage

```bash
go run server.go | prettylog
```

this will format the logs from `server.go` to a more readable format.

You can also specify the log level to filter the logs.

```bash
go run server.go -c ./config/app.yml | prettylog -level=info
```
