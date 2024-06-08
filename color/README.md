# 🎨 color

A fast and lightweight color library for GoLanguage

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
=== RUN   BenchmarkRed
BenchmarkRed
BenchmarkRed-8
31155721                37.11 ns/op           16 B/op          1 allocs/op
PASS
ok      github.com/tsingshaner/go-pkg/color     1.231s

=== RUN   BenchmarkMulti
BenchmarkMulti
BenchmarkMulti-8
 9097593               128.6 ns/op            88 B/op          3 allocs/op
PASS
ok      github.com/tsingshaner/go-pkg/color     1.338s

=== RUN   BenchmarkUnsafeMagenta
BenchmarkUnsafeMagenta
BenchmarkUnsafeMagenta-8
70216500                18.10 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/tsingshaner/go-pkg/color     1.325s

=== RUN   BenchmarkUnsafeMulti
BenchmarkUnsafeMulti
BenchmarkUnsafeMulti-8
16142202                73.15 ns/op           48 B/op          1 allocs/op
PASS
ok      github.com/tsingshaner/go-pkg/color     1.289s
```

## ✨ Inspired

 [`picocolors`](https://github.com/alexeyraspopov/picocolors)

## 📄 License
[ISC](github.com/tsingshaner.com/go-pkg/LICENSE) License © 2024-Present [qingshaner](gitub.com/tsingshaner)
