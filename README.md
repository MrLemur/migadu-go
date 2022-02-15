# Migadu API in Go

`migadu-go` is a Go library for interfacing with the [Migadu API](https://www.migadu.com/api/). It currently supports all endpoints available through the REST API.

- [Installing](#installing)
- [Client](#client)
- [Operations](#operations)

## Installing

**go get**:

```go
go get github.com/MrLemur/migadu-go
```

## Client

A client is required for all methods of the library. A client is scoped to a single domain (e.g. `example.com`).

You will need an admin email address and an API key to create a client (API keys can be made [here](https://admin.migadu.com/account/api/keys)).

Use `miagdu.New("admin_email", "api_key", "domain_name")` to create a new client.

Example:

```go
package main

import (
    "fmt"
    "os"

    "github.com/MrLemur/migadu-go"
)

client, err := migadu.New("admin_email@example.com", "xxxxxxxxxxxxxxxxxx", "example.com")

// Incorrect API details will return an error
if err != nil {
    fmt.Println(err)
    os.Exit(1)
}
...
```

## Operations

For each type of entity, the following operations are available: `List`,`Get`,`New`,`Update`,`Delete`.

Each method requires a context to operate - a single `ctx = context.Background()` will suffice for most operations.

The Migadu doesn't return very useful status codes when errors occur. Generally anything other than status code `200` on a response indicates the operation failed. This can include things like a mailbox existing when trying to create one with the same name.

Each method will return a non `nil` error if an operation fails.
