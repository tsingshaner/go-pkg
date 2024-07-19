# ⚙️ Conf

A configuration library for Go language based on Viper.

<a href="https://pkg.go.dev/github.com/tsingshaner/go-pkg/conf" alt="Go Reference"><img src="https://pkg.go.dev/badge/github.com/tsingshaner/go-pkg/conf.svg" /></a>
<a alt="Go Report Card" href="https://goreportcard.com/report/github.com/tsingshaner/go-pkg/conf"><img src="https://goreportcard.com/badge/github.com/tsingshaner/go-pkg/conf" /></a>

## 📦 Usage

```bash
go get -u github.com/tsingshaner/go-pkg/conf
```

```go
package main

import "github.com/tsingshaner/go-pkg/conf"

type Config struct {
  Name string `json:"name" yaml:"name" toml:"name"`
  Age  int    `json:"age" yaml:"age" toml:"age"`
}

func main() {
  config := conf.New(&Config{}, conf.ParseArgs())

	config.Load()
	config.Watch(func(e conf.Event) {
		println("config file changed reload config file")
		config.Load()
	})
}
```

```bash
go run main.go --config=config.yaml
```

## 📄 License

[ISC](../LICENSE) License © 2024-Present [qingshaner](https://gitub.com/tsingshaner)
