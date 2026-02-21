# Migadu API in Go

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/MrLemur/migadu-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/MrLemur/migadu-go)](https://goreportcard.com/report/github.com/MrLemur/migadu-go)

`migadu-go` is a Go library for interfacing with the [Migadu API](https://www.migadu.com/api/).

## Installing

```go
go get github.com/MrLemur/migadu-go
```

## Usage

See the [pkg.go.dev documentation](http://pkg.go.dev/github.com/MrLemur/migadu-go) for full API reference.

Create a client with your admin email and API key ([generate one here](https://admin.migadu.com/account/api/keys)):

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/MrLemur/migadu-go"
)

func main() {
    client, err := migadu.New("admin@example.com", "api_key")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    ctx := context.Background()

    domain, err := client.NewDomain(ctx, &migadu.Domain{Name: "example.com"})
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    mailboxes, err := client.ListMailboxes(ctx, domain)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(mailboxes)
}
```

All methods accept a `context.Context` as the first argument and return an error if the operation fails.
