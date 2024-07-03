# 📄 log

<p align="">
<a href="https://pkg.go.dev/github.com/tsingshaner/go-pkg/log" alt="Go Reference"><img src="https://pkg.go.dev/badge/github.com/tsingshaner/go-pkg/log.svg" /></a>
<a alt="Go Report Card" href="https://goreportcard.com/report/github.com/tsingshaner/go-pkg/log"><img src="https://goreportcard.com/badge/github.com/tsingshaner/go-pkg/log" /></a>
</p>

A ready-to-use stubborn log library encapsulates the log module, providing functions such as log output, log level.

## ✨ Features

- 可选择 `log/slog`, `zap` 作为日志库
- 支持本地日志文件轮转


## 📦 Usage

```shell
go get -u github.com/tsingshaner/go-pkg/log
```

### 🌿 Use log/slog

[example/log/slog](../internal/example/log/slog/main.go)

### ⚡ Use zap

Use zap you need add zap to your project dependencies.

```shell
go get -u go.uber.org/zap
```

[example/log/zap](../internal/example/log/zap/main.go)

## 🎨 Pretty log

If you want a pretty log in the console for development, you can install a cli tool `prettylog`.

```shell
go install github.com/tsingshaner/go-pkg/prettylog/cmd/prettylog
```

Then you can use it like this:

```shell
go run main.go | prettylog
```

## 📄 License

[ISC](../LICENSE) License © 2024-Present [qingshaner](https://gitub.com/tsingshaner)
