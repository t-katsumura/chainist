# chainist

[![GoDoc](https://godoc.org/github.com/t-katsumura/chainist?status.svg)](http://godoc.org/github.com/t-katsumura/chainist)
[![Go Report Card](https://goreportcard.com/badge/github.com/t-katsumura/chainist)](https://goreportcard.com/report/github.com/t-katsumura/chainist)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Test](https://github.com/t-katsumura/chainist/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/t-katsumura/chainist/actions/workflows/test.yml?query=branch%3Amain)
[![Codecov](https://codecov.io/gh/t-katsumura/chainist/branch/main/graph/badge.svg?token=3ZRzIQTXIw)](https://codecov.io/gh/t-katsumura/chainist)

`chainist` is a simple http handler, or middleware, chaining library for golang.

**Underlying Concepts**

![concepts.png](/docs/images/concepts.png)

## Installation

```
go get github.com/t-katsumura/chainist@latest
```

## Basic Usage

Create a new chain.  
Http handlers or middleware can be added this time.

```go
chain := chainist.NewChain()
// or
chain := chainist.NewChain(handler1, handler2)
```

Add http handlers, i.e. `func(next http.Handler) http.Handler`, to the chain.

```go
// add handlers one by one
chain.Append(handler1)
chain.Append(handler2)

// add handlers using method chaining
chain.Append(handler1).Append(handler2)

// add handlers at a time
chain.Extend(handler1, handler2)
```

Add http handler function, i.e. `func(w http.ResponseWriter, r *http.Request)`, to the chain.  
Here, 'Pre' means the functions will be executed before calling succeeding handlers. 'Post' means the functions will be executed after calling succeeding handlers.

```go
// add handlerFuncss one by one
chain.AppendPreFunc(handlerFunc1)
chain.AppendPreFunc(handlerFunc2)

// add handlerFuncs using method chaining
chain.AppendPreFunc(handlerFunc1).AppendPreFunc(handlerFunc2)

// add handlerFuncs at a time
chain.ExtendPreFunc(handlerFunc1, handlerFunc2)
```

```go
// add handlerFuncs one by one
chain.AppendPostFunc(handlerFunc1)
chain.AppendPostFunc(handlerFunc2)

// add handlerFuncs using method chaining
chain.AppendPostFunc(handlerFunc1).AppendPostFunc(handlerFunc2)

// add handlerFuncs at a time
chain.ExtendPostFunc(handlerFunc1, handlerFunc2)
```

Get the handler chain which type is http.Handler.

```go
chain.Chain()

// or with handler function
chain.ChainFunc(handlerFunc)
```

## Example

This is an example of chainist.
With chainist, you don't need to wrap handlerFunc using http.Handler and http.HandlerFunc. chainist internally wrap your handlerFunc.

```go
package main

import (
    "net/http"

    "github.com/t-katsumura/chainist"
)

// sample http handlerFunc
func handlerFunc1(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hi from handlerFunc 1!\n"))
}

// sample http handlerFunc
func handlerFunc2(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hi from handlerFunc 2!\n"))
}

// sample http handler
func handler1(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hi from handler 1!\n"))
        next.ServeHTTP(w, r)
    })
}

// sample http handler
func handler2(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hi from handler 2!\n"))
        next.ServeHTTP(w, r)
    })
}

func main() {

    // create a new chain
    chain := chainist.NewChain()

    // add handler functions
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)

    // add handlers
    chain.Extend(handler1, handler2)

    // Run http server which responds
    /*
        Hi from handlerFunc 1!
        Hi from handler 1!
        Hi from handler 2!
        Hi from handlerFunc 2!
    */
    http.ListenAndServe(":8080", chain.Chain())
}
```

## Questions and support

All bug reports, questions and suggestions should go though Github Issues.

## Contributing

1. Fork it
1. Create feature branch (`git checkout -b feature/new-feature`)
1. Write codes on feature branch
1. Commit your changes (`git commit -m "Added new feature"`)
1. Push to the branch (`git push -u origin feature/new-feature`)
1. Create new Pull Request on Github

## Development

- Write codes
- `go fmt -x ./...` - format codes
- `go test -v -cover ./...` - run test and the coverage should always be 100%

When generating coverage report, run

1. `go test -cover ./... -coverprofile=cover.out`
2. `go tool cover -html=cover.out -o cover.html`

## Similar packages

[alice](https://github.com/justinas/alice) is a famous middleware chaining library.
This project is inspired by [alice](https://github.com/justinas/alice).
