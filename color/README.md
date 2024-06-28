# 🎨 color

<p align="center"><img alt="snapshot" src="https://github.com/tsingshaner/go-pkg/assets/170866467/6cbba7b9-df1b-40de-8e67-45c37a0ccaa6" /></p>

<p align="center">A fast and lightweight color library for Go language.</p>

<p align="center">
<a href="https://pkg.go.dev/github.com/tsingshaner/go-pkg/color" alt="Go Reference"><img src="https://pkg.go.dev/badge/github.com/tsingshaner/go-pkg/color.svg" /></a>
<a alt="Go Report Card" href="https://goreportcard.com/report/github.com/tsingshaner/go-pkg/color"><img src="https://goreportcard.com/badge/github.com/tsingshaner/go-pkg/color" /></a>
</p>

## 📦 Usage

```bash
go get -u github.com/tsingshaner/go-pkg/color
```

```go
package main

import "github.com/tsingshaner/go-pkg/color"

func main() {
  print(color.Red("Hello, World!\n"))

  // multiple styles
  print(color.Underline(color.Bold(color.Red("Hello, World!\n"))))
}
```

## 🚨 Unsafe API

The unsafeXXX func is not safe but is faster, if you ensure the text is safe, you can use it.

If you use nesting unsafeXXX func, it may close the color style unexpectedly.

```go
func TestUnsafeWithFail(t *testing.T) {
	mixinString := UnsafeRed(" red " + Green(" green ") + " no color ")

	assert.Equal(t, "\x1b[31m red \x1b[32m green \x1b[39m no color \x1b[39m", mixinString)
}

func TestSafeFormat(t *testing.T) {
	mixinString := Red(" red " + Green(" green ") + " red ")

	assert.Equal(t, "\x1b[31m red \x1b[32m green \x1b[31m red \x1b[39m", mixinString)
}

```

## ⚡ Benchmark

```bash
goos: windows
goarch: amd64
pkg: github.com/tsingshaner/go-pkg/color
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkRed-8                  32791454                34.78 ns/op
BenchmarkUnsafeMagenta-8        76925048                15.84 ns/op
BenchmarkMulti-8                 9167505               120.5 ns/op
BenchmarkUnsafeMulti-8          15919278                69.64 ns/op
PASS
ok      github.com/tsingshaner/go-pkg/color     4.878s
```

## ✨ Inspired By

[`picocolors`](https://github.com/alexeyraspopov/picocolors)

## 📄 License

[ISC](../LICENSE) License © 2024-Present [qingshaner](https://gitub.com/tsingshaner)
