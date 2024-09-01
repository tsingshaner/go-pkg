# 🧪 H3Test

<p align="">
<a href="https://pkg.go.dev/github.com/tsingshaner/go-pkg/h3test" alt="Go Reference"><img src="https://pkg.go.dev/badge/github.com/tsingshaner/go-pkg/h3test.svg" /></a>
<a alt="Go Report Card" href="https://goreportcard.com/report/github.com/tsingshaner/go-pkg/h3test"><img src="https://goreportcard.com/badge/github.com/tsingshaner/go-pkg/h3test" /></a>
</p>

A helper func for testing http requests.

## 📦 Usage

```bash
go get -u github.com/tsingshaner/go-pkg/h3test
```

```go
package e2e

// import ...

type Data struct {
    Code string `json:"code"`
}

var handler http.Handler

func TestMain(m *testing.M) {
    server = &http.Server{
      // ...
    }

    handler = server.Handler

    os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
    request := h3test.New("/api/auth/register").POST()
    payload := map[string]string{
        "username": "test",
        "password": "123456",
    }

    response := request.JSON(payload).Send(handler)

    assert.Equal(t, http.StatusCreated, res.Code)

    data := h3test.ExtractJSON[Data](res)

    assert.Equal(t, "Custom Success Code", data.Code)
}

```
