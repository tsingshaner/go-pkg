# color

a fast and lightweight color library for GoLanguage

✨ Inspired by [`picocolors`](https://github.com/alexeyraspopov/picocolors)

## Usage

```go
package main

import (
  "fmt"
  "github.com/tsingshaner/go-pkg/color"
)

func main() {
  print(color.Red("Hello, World!\n"))

  // multiple styles
  print(color.Underline(color.Bold(color.Red("Hello, World!\n"))))
}
```

## ⚡ Benchmark

```bash
goos: windows
goarch: amd64
pkg: github.com/tsingshaner/go-pkg/color
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
=== RUN   BenchmarkRed
BenchmarkRed
BenchmarkRed-8
72651098                15.89 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/tsingshaner/go-pkg/color     1.199s
```

```bash
goos: windows
goarch: amd64
pkg: github.com/tsingshaner/gopkg/color
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
=== RUN   BenchmarkMulti
BenchmarkMulti
BenchmarkMulti-8
22164634                52.66 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/tsingshaner/gopkg/color      1.254s
```
