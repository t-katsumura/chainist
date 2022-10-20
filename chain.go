package chainist

import (
	"net/http"
)

/*
Define the "middleware" type.
Note, this is the commonly used definition of middleware.

An example of basic definition of middleware is

    func YourMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // you can write your code here
            if next != nil {
                next.ServeHTTP(w, r)
            }
            // you can write your code here
        })
    }
*/
type Middleware func(h http.Handler) http.Handler

/*
Chain is the struct for middleware chain.
This holds the set of some middleware (i.e. func(h http.Handler) http.Handler)
and a http handler function (i.e. `func(http.ResponseWriter, *http.Request)`)
*/
type Chain struct {
	// Middleware is the list of middleware.
	// This can contain nil values but they are ignored
	// when getting middleware chains with Chain() or ChainFunc().
	Middleware []Middleware

	// HandlerFunc is the http handler function at the edge of the chain.
	// If it is not set before calling Chain() or ChainFunc(),
	// an empty handler function is automatically used.
	HandlerFunc http.HandlerFunc
}

/*
NewChain creates a new middleware chain.

    chain := chainist.NewChain()

Middleware can be added at the time of creation of a chain.
Middleware is a function with the signature of `func(h http.Handler) http.Handler`.
If nil is contained in the given arguments, nil is returned.

    // handler1 is called at first, and handler3 at last
    chain := chainist.NewChain(handler1, handler2, handler3)
*/
func NewChain(ms ...Middleware) *Chain {
	for _, m := range ms {
		if m == nil {
			// it might be better to return an error
			return nil
		}
	}
	c := &Chain{
		Middleware: ms,
	}
	return c
}

/*
Append middleware at the last of the chain.
If nil is given, then the chain will be returned without adding it.

    chain := chainist.NewChain()
    chain.Append(handler1)
         .Append(handler2)
         .Append(handler3)
*/
func (c *Chain) Append(m Middleware) *Chain {
	if m == nil {
		return c
	}
	c.Middleware = append(c.Middleware, m)
	return c
}

/*
AppendPreFunc appends a http handler function as pre-executable function which is executed before invoking succeeding middleware.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If nil is given as handler function, then the chain will be returned as it is.

If you pass `YourPreHandlerFunc(w http.ResponseWriter, r *http.Request)` as an argument, it is treaded as the middleware of

    func handler(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            YourPreHandlerFunc(w, r)
            if next != nil {
                // invoke the next handler function after YourPreHandlerFunc
                next.ServeHTTP(w, r)
            }
        })
    }

Usage example

    chain := chainist.NewChain()

    chain.AppendPreFunc(handlerFunc1)
         .AppendPreFunc(handlerFunc2)
         .AppendPreFunc(handlerFunc3)
*/
func (c *Chain) AppendPreFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{HandlerFunc: f}
	return c.Append(h.PreMiddleware)
}

/*
AppendPostFunc appends a http handler function as post-executable function which is executed after invoked succeeding middleware.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If nil is given as handler function, then the chain will be returned as it is.

If you pass `YourPostHandlerFunc(w http.ResponseWriter, r *http.Request)` as an argument, it is treaded as the middleware of

    func handler(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if next != nil {
                // invoke the next handler function before YourPostHandlerFunc
                next.ServeHTTP(w, r)
            }
            YourPostHandlerFunc(w, r)
        })
    }

Usage:

    chain := chainist.NewChain()
    chain.AppendPostFunc(handlerFunc1)
         .AppendPostFunc(handlerFunc2)
         .AppendPostFunc(handlerFunc3)
*/
func (c *Chain) AppendPostFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{HandlerFunc: f}
	return c.Append(h.PostMiddleware)
}

/*
Insert inserts middleware at designated position of the chain.
Middleware must have the signature of `func(h http.Handler) http.Handler`.
If nil is given as middleware, then the chain will be returned without adding it.

If the number less than 0 is given as the position, given handler is added at the first of the chain.
If the number grater than the length of middleware is given as the position, given handler is added at the last of the chain.

Usage:

    // create a middleware chain with the order of handler1,handler2,handler3
    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handler4 between handler1 and handler2
    // chain.Middleware[1] will be handler4
    chain.Insert(handler4, 1)

    // insert handler5 at the first of the chain
    // chain.Middleware[0] will be handler4
    chain.Insert(handler5, 0)
*/
func (c *Chain) Insert(m Middleware, i int) *Chain {
	if m == nil {
		return c
	}
	if len(c.Middleware) == 0 || i >= len(c.Middleware) {
		c.Middleware = append(c.Middleware, m)
	} else {
		if i < 0 {
			i = 0
		}
		c.Middleware = append(c.Middleware[:i+1], c.Middleware[i:]...)
		c.Middleware[i] = m
	}
	return c
}

/*
Insert a pre-executable handler function at designated number of chain.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If nil is given as handler function, the chain will be returned without inserting it.

If the number less than 0 is given as the position, given http handler function is added at the first of the chain.
If the number grater than the length of middleware is given as the position, the given http handler function is added at the last of the chain.

    // create a middleware chain with the order of handler1,handler2,handler3
    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handlerFunc4 between handler1 and handler2
    // chain.Middleware[1] will be a middleware created with handlerFunc4
    chain.InsertPreFunc(handlerFunc4, 1)

    // insert handlerFunc5 at the first of the chain
    // chain.Middleware[0] will be a middleware created with handlerFunc5
    chain.InsertPreFunc(handlerFunc5, 0)
*/
func (c *Chain) InsertPreFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{HandlerFunc: f}
	return c.Insert(h.PreMiddleware, i)
}

/*
Insert a post-executable handler function at designated number of chain.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If nil is given as handler function, the chain will be returned without inserting it.

If the number less than 0 is given as the position, given http handler function is added at the first of the chain.
If the number grater than the length of middleware is given as the position, the given http handler function is added at the last of the chain.

    // create a middleware chain with the order of handler1,handler2,handler3
    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handlerFunc4 between handler1 and handler2
    // chain.Middleware[1] will be a middleware created with handlerFunc4
    chain.InsertPostFunc(handlerFunc4, 1)

    // insert handlerFunc5 at the first of the chain
    // chain.Middleware[0] will be a middleware created with handlerFunc5
    chain.InsertPostFunc(handlerFunc5, 0)
*/
func (c *Chain) InsertPostFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{HandlerFunc: f}
	return c.Insert(h.PostMiddleware, i)
}

/*
Extend appends multiple middleware at a time.
This function append multiple middleware at the end of the chain.
If nil is contained in the arguments, they are ignored.

    // create a middleware chain
    chain := chainist.NewChain(handler1, handler2)

    // append handler3 and handler4 with this order
    // chain will have handler1,handler2,handler3,handler4 with this order
    chain.Extend(handler3, handler4)
*/
func (c *Chain) Extend(ms ...Middleware) *Chain {
	for _, m := range ms {
		if m == nil {
			continue
		}
		c.Append(m)
	}
	return c
}

/*
ExtendPreFunc appends multiple pre-executable handler functions at a time.
This function append multiple handler function at the end of the chain.
nil is ignored if contained in the arguments.

    // create a middleware chain
    chain := chainist.NewChain(handler1, handler2)

    // append middleware created with handlerFunc3 and handlerFunc4 with this order
    // chain will have handler1,handler2,handler3,handler4 with this order
    chain.ExtendPreFunc(handlerFunc3, handlerFunc4)
*/
func (c *Chain) ExtendPreFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{HandlerFunc: f}
		c.Append(h.PreMiddleware)
	}
	return c
}

/*
ExtendPostFunc appends multiple post-executable handler functions at a time.
This function append multiple handler function at the end of the chain.
nil is ignored if contained in the arguments.

    // create a middleware chain
    chain := chainist.NewChain(handler1, handler2)

    // append middleware created with handlerFunc3 and handlerFunc4 with this order
    // chain will have handler1,handler2,handler3,handler4 with this order
    chain.ExtendPostFunc(handlerFunc3, handlerFunc4)
*/
func (c *Chain) ExtendPostFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{HandlerFunc: f}
		c.Append(h.PostMiddleware)
	}
	return c
}

/*
SetHandlerFunc sets the handler function which will be invoked at the edge of the chain.
If nil is given as the argument, it is ignored.

    // create a middleware chain with handler1 and handler2
    chain := chainist.NewChain(handler1, handler2)

    // set the handler function of handlerFuncAtEdge
    chain.SetHandlerFunc(handlerFuncAtEdge)
*/
func (c *Chain) SetHandlerFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	c.HandlerFunc = f
	return c
}

/*
Join joins two chains.

    // create two chains with handlers.
    chain1 := chainist.NewChain(handler1, handler2)
    chain2 := chainist.NewChain(handler3, handler4)

    // join two chains
    // this operation extends chain1 with chain2
    // chain1 will have handler1,handler2,handler3,handler4 with this order
    chain1.Join(chain2)
*/
func (c *Chain) Join(o *Chain) *Chain {
	if o == nil {
		return c
	}
	c.Middleware = append(c.Middleware, o.Middleware...)
	return c
}

/*
Len gets the length of middleware chain.
This contains the length of Middleware and the Handler Function if it's already set.

    chain := chainist.NewChain(handler1, handler2)
    chain.SetHandlerFunc(handlerFuncAtEdge)

    // this shows 3
    println(chain.len())
*/
func (c *Chain) Len() int {
	length := len(c.Middleware)
	if c.HandlerFunc != nil {
		length += 1
	}
	return length
}

/*
Chain returns a new middleware chain.
This function returns nil if there is no middleware and no handler function.

Usage:

    // defining a middleware chain
    // this create the chain of [handlerFunc1, handlerFunc2, handler1, handler2]
    chain := chainist.NewChain()
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)
    chain.Extend(handler1, handler2)

    // get a new handler from chain
    handler := chain.Chain()
    if handler == nil {
        panic("handler is nil")
    }

    // run http server
    http.ListenAndServe(":8080", handler)
*/
func (c *Chain) Chain() http.Handler {

	n := c.Len()

	// no handlers or handler functions are set in the chain
	if n < 1 {
		return nil
	}

	var h http.Handler
	if c.HandlerFunc != nil {
		c.Append((&HandlerFuncWrapper{HandlerFunc: c.HandlerFunc}).Middleware)
	}

	h = c.Middleware[n-1](nil)
	for i := range c.Middleware[:n-1] {
		h = c.Middleware[n-2-i](h)
	}

	return h
}

/*
ChainFunc returns a new middleware chain with a handler function.
If the given handler function is not nil, it is used instead of the handler function set with `SetHandlerFunc()`.

Usage:

    // defining a middleware chain
    chain := chainist.NewChain()
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)
    chain.Extend(handler1, handler2)

    // get a new handler from chain
    handler := chain.ChainFunc(handlerFuncAtEdge)
    if handler == nil {
        panic("handler is nil")
    }

    // run http server
    http.ListenAndServe(":8080", handler)
*/
func (c *Chain) ChainFunc(f http.HandlerFunc) http.Handler {
	if f != nil {
		c.HandlerFunc = f
	}
	return c.Chain()
}
