# chainist

[![GoDoc](https://godoc.org/github.com/t-katsumura/chainist?status.svg)](http://godoc.org/github.com/t-katsumura/chainist)
[![Go Report Card](https://goreportcard.com/badge/github.com/t-katsumura/chainist)](https://goreportcard.com/report/github.com/t-katsumura/chainist)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)
[![Test](https://github.com/t-katsumura/chainist/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/t-katsumura/chainist/actions/workflows/test.yml?query=branch%3Amain)
[![Codecov](https://codecov.io/gh/t-katsumura/chainist/branch/main/graph/badge.svg?token=P5J4J1F6RN)](https://codecov.io/gh/t-katsumura/chainist)


`chainist` is a simple go library to create a http handler chain. This is also known as middleware.


## Installation

```
go get github.com/t-katsumura/chainist@latest
```

## Basic Usage

create new chain. http handlers can be added at this time.

```go
chain := chainist.NewChain()
// or
chain := chainist.NewChain(handler1, handler2)
```

add http handler, i.e. `func(next http.Handler) http.Handler`, to the chain.

```go
// add handlers one by one
chain.Append(handler1)
chain.Append(handler2)
// add handlers using method chaining 
chain.Append(handler1).Append(handler2)
// add handlers at a time
chain.Extend(handler1, handler2)
```

add http handlerFunc, i.e. `func(w http.ResponseWriter, r *http.Request)`, to the chain.  
Here, 'Pre' means the handlerFuncs will be executed before calling successing handlers or handlerFuncs. And 'Post' means the handlerFuncs will be executed after calling successing handlers or handlerFuncs.

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

get the http.Handler

```go
chain.Chain()
// or with handlerFunc
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

func handlerFunc1(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hi from handlerFunc 1!\n"))
}

func handlerFunc2(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hi from handlerFunc 2!\n"))
}

func handler1(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hi from handler 1!\n"))
        next.ServeHTTP(w, r)
    })
}

func handler2(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hi from handler 2!\n"))
        next.ServeHTTP(w, r)
    })
}

func main() {

    // create a new chain struct
    chain := chainist.NewChain()

    // add handler functions
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)
    // add handlers
    chain.Extend(handler1, handler2)

    // Run http server which responds
    /*
       Hi from handlerFunc 1!
       Hi from handlerFunc 2!
       Hi from handler 1!
       Hi from handler 2!
    */
    http.Handle("/", chain.Chain())
    http.ListenAndServe(":8080", nil)
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
- `go fmt -x  ./...` - format codes
- `go test -v -cover ./...` - run test and the coverage should always be 100%
