# Migadu API in Go

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/MrLemur/migadu-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/MrLemur/migadu-go)](https://goreportcard.com/report/github.com/MrLemur/migadu-go)

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

A client is required for all methods of the library. The client is account-level and can manage multiple domains.

You will need an admin email address and an API key to create a client (API keys can be made [here](https://admin.migadu.com/account/api/keys)).

Use `migadu.New("admin_email", "api_key")` to create a new client.

Example:

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/MrLemur/migadu-go"
)

func main() {
    client, err := migadu.New("admin_email@example.com", "xxxxxxxxxxxxxxxxxx")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    ctx := context.Background()

    // Create a new domain
    domain, err := client.NewDomain(ctx, &migadu.Domain{
        Name: "example.com",
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Create a new mailbox
    mailbox, err := client.NewMailbox(ctx, domain, &migadu.Mailbox{
        LocalPart:             "hello",
        Name:                  "Hello User",
        PasswordRecoveryEmail: "recovery@example.com",
        PasswordMethod:        "invitation",
    })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // List mailboxes
    mailboxes, err := client.ListMailboxes(ctx, domain)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
```

## Operations

All operations follow a consistent pattern:

### Pattern

- **List** operations return `[]Type` (slice of structs)
- **Get/Update/Delete** operations accept `*Type` (pointer to struct)
- **New** operations accept `*Type` (user constructs the struct with required fields)

### Account-Level Operations

Domain management operations:

- `ListDomains(ctx)` → `[]Domain`
- `GetDomain(ctx, *Domain)` → `*Domain`
- `NewDomain(ctx, *Domain)` → `*Domain`
- `UpdateDomain(ctx, *Domain)` → `*Domain`
- `GetDomainRecords(ctx, *Domain)` → `[]DNSRecord`
- `GetDomainRecordsDetailed(ctx, *Domain)` → `*DomainRecords`
- `GetDomainDiagnostics(ctx, *Domain)` → `*DomainDiagnostics`
- `ActivateDomain(ctx, *Domain)` → `*Domain`

### Domain-Scoped Operations

All require `*Domain` as first parameter after context:

**Mailboxes:**

- `ListMailboxes(ctx, *Domain)` → `[]Mailbox`
- `GetMailbox(ctx, *Domain, localPart)` → `*Mailbox`
- `NewMailbox(ctx, *Domain, *Mailbox)` → `*Mailbox`
- `UpdateMailbox(ctx, *Domain, *Mailbox)` → `*Mailbox`
- `DeleteMailbox(ctx, *Domain, *Mailbox)` → `error`

**Forwardings:**

- `ListForwardings(ctx, *Domain, mailbox)` → `[]Forwarding`
- `GetForwarding(ctx, *Domain, mailbox, address)` → `*Forwarding`
- `NewForwarding(ctx, *Domain, mailbox, *Forwarding)` → `*Forwarding`
- `UpdateForwarding(ctx, *Domain, mailbox, *Forwarding)` → `*Forwarding`
- `DeleteForwarding(ctx, *Domain, mailbox, *Forwarding)` → `error`

**Aliases:**

- `ListAliases(ctx, *Domain)` → `[]Alias`
- `GetAlias(ctx, *Domain, localPart)` → `*Alias`
- `NewAlias(ctx, *Domain, *Alias)` → `*Alias`
- `UpdateAlias(ctx, *Domain, *Alias)` → `*Alias`
- `DeleteAlias(ctx, *Domain, *Alias)` → `error`

**Rewrites:**

- `ListRewrites(ctx, *Domain)` → `[]Rewrite`
- `GetRewrite(ctx, *Domain, name)` → `*Rewrite`
- `NewRewrite(ctx, *Domain, *Rewrite)` → `*Rewrite`
- `UpdateRewrite(ctx, *Domain, *Rewrite)` → `*Rewrite`
- `DeleteRewrite(ctx, *Domain, *Rewrite)` → `error`

**Identities:**

- `ListIdentities(ctx, *Domain, mailbox)` → `[]Identity`
- `GetIdentity(ctx, *Domain, mailbox, localPart)` → `*Identity`
- `NewIdentity(ctx, *Domain, mailbox, *Identity)` → `*Identity`
- `UpdateIdentity(ctx, *Domain, mailbox, *Identity)` → `*Identity`
- `DeleteIdentity(ctx, *Domain, mailbox, *Identity)` → `error`

Each method requires a context - `ctx := context.Background()` will suffice for most operations.

The Migadu API doesn't return very useful status codes when errors occur. Generally anything other than status code `200` indicates failure. This can include things like a mailbox already existing.

Each method returns an error if an operation fails.
