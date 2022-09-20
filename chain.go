package chainist

import (
	"net/http"
)

// Define middleware function type.
type middleware func(h http.Handler) http.Handler

// Struct for middleware chain.
type Chain struct {
	// fs : middlewars
	fs []middleware
	// f : handler function at the edge of the chain
	f http.HandlerFunc
}

/*
Create a new chain struct.

    chain := chainist.NewChain()

Middleware can be set at the time of creating a Chain struct.
Middleware must have the signature of `func(h http.Handler) http.Handler`

    chain := chainist.NewChain(handler1, handler2, handler3)
*/
func NewChain(fs ...middleware) *Chain {
	c := &Chain{
		fs: fs,
	}
	return c
}

/*
Append middleware at the last of the chain.
If null is given as middleware, the chain returned as it is.

    chain := chainist.NewChain()
    chain.Append(handler1)
         .Append(handler2)
         .Append(handler3)
*/
func (c *Chain) Append(f middleware) *Chain {
	if f == nil {
		return c
	}
	c.fs = append(c.fs, f)
	return c
}

/*
Append handler function as pre-executable function which is executed before invoking succeeding middleware.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If null is given as handler function, the chain returned as it is.

    chain := chainist.NewChain()
    chain.AppendPreFunc(handlerFunc1)
         .AppendPreFunc(handlerFunc2)
         .AppendPreFunc(handlerFunc3)
*/
func (c *Chain) AppendPreFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Append(h.PreMiddleware)
}

/*
Append handler function as post-executable function which is executed after invoking succeeding middleware.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If null is given as handler function, the chain returned as it is.

    chain := chainist.NewChain()
    chain.AppendPostFunc(handlerFunc1)
         .AppendPostFunc(handlerFunc2)
         .AppendPostFunc(handlerFunc3)
*/
func (c *Chain) AppendPostFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Append(h.PostMiddleware)
}

/*
Insert middleware at designated position of the chain.
Middleware must have the signature of `func(h http.Handler) http.Handler`.
If null is given as middleware, the chain returned as it is.

    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handler4 between handler1 and handler2
    chain.Insert(handler4, 1)

    // insert handler5 at the first of the chain
    chain.Insert(handler5, 0)
*/
func (c *Chain) Insert(f middleware, i int) *Chain {
	if f == nil {
		return c
	}
	if len(c.fs) == 0 || i >= len(c.fs) {
		c.fs = append(c.fs, f)
	} else {
		if i < 0 {
			i = 0
		}
		c.fs = append(c.fs[:i+1], c.fs[i:]...)
		c.fs[i] = f
	}
	return c
}

/*
Insert a pre-executable handler function at designated number of chain.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If null is given as handler function, the chain returned as it is.

    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handlerFunc4 between handler1 and handler2
    chain.InsertPreFunc(handlerFunc4, 1)

    // insert handlerFunc5 at the first of the chain
    chain.InsertPreFunc(handlerFunc5, 0)
*/
func (c *Chain) InsertPreFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Insert(h.PreMiddleware, i)
}

/*
Insert a post-executable handler function at designated number of chain.
Handler function must have the signature of `func(w http.ResponseWriter, r *http.Request)`.
If null is given as handler function, the chain returned as it is.

    chain := chainist.NewChain(handler1, handler2, handler3)

    // insert handlerFunc4 between handler1 and handler2
    chain.InsertPostFunc(handlerFunc4, 1)

    // insert handlerFunc5 at the first of the chain
    chain.InsertPostFunc(handlerFunc5, 0)
*/
func (c *Chain) InsertPostFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Insert(h.PostMiddleware, i)
}

/*
Append multiple middleware at a time.
This function append multiple middleware at the end of the chain.

    chain := chainist.NewChain(handler1, handler2)
    chain.Extend(handler3, handler4)
*/
func (c *Chain) Extend(fs ...middleware) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		c.Append(f)
	}
	return c
}

/*
Append multiple pre-executable handler functions at a time.
This function append multiple handler function at the end of the chain.

    chain := chainist.NewChain(handler1, handler2)
    chain.ExtendPreFunc(handlerFunc3, handlerFunc4)
*/
func (c *Chain) ExtendPreFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{f: f}
		c.Append(h.PreMiddleware)
	}
	return c
}

/*
Append multiple post-executable handler functions at a time.
This function append multiple handler function at the end of the chain.

    chain := chainist.NewChain(handler1, handler2)
    chain.ExtendPostFunc(handlerFunc3, handlerFunc4)
*/
func (c *Chain) ExtendPostFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{f: f}
		c.Append(h.PostMiddleware)
	}
	return c
}

/*
Set the handler function which is invoked at the edge of the chain.

    chain := chainist.NewChain(handler1, handler2)
    chain.SetHandlerFunc(handlerFuncAtEdge)
*/
func (c *Chain) SetHandlerFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	c.f = f
	return c
}

/*
Join two chains.

    chain1 := chainist.NewChain(handler1, handler2)
    chain2 := chainist.NewChain(handler3, handler4)

    // join two chains
    // this operation extends chain1 with the middleware of chain2
    chain1.Join(chain2)
*/
func (c *Chain) Join(o *Chain) *Chain {
	if o == nil {
		return c
	}
	c.fs = append(c.fs, o.fs...)
	return c
}

/*
Get the length of middleware chain.

    chain := chainist.NewChain(handler1, handler2)
    chain.SetHandlerFunc(handlerFuncAtEdge)

    // this shows 3
    println(chain.len())
*/
func (c *Chain) Len() int {
	length := len(c.fs)
	if c.f != nil {
		length += 1
	}
	return length
}

/*
Return new middleware chain.
This function returns nil is no middleware is configured.

    // defining a middleware chain
    // this create the chain of [handlerFunc1, handlerFunc2, handler1, handler2]
    chain := chainist.NewChain()
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)
    chain.Extend(handler1, handler2)

    // run http server
    http.ListenAndServe(":8080", chain.Chain())
*/
func (c *Chain) Chain() http.Handler {

	n := c.Len()

	if n < 1 {
		return nil
	}

	var h http.Handler
	if c.f != nil {
		c.Append((&HandlerFuncWrapper{f: c.f}).Middleware)
	}

	h = c.fs[n-1](nil)
	for i := range c.fs[:n-1] {
		h = c.fs[n-2-i](h)
	}

	return h
}

/*
Return new middleware chain with a handler function.
If no handler function at the edge of the chain is set with `SetHandlerFunc()` method,
handler function is set at the time of generating middleware chain.

    // defining a middleware chain
    chain := chainist.NewChain()
    chain.AppendPreFunc(handlerFunc1)
    chain.AppendPostFunc(handlerFunc2)
    chain.Extend(handler1, handler2)

    // run http server
    http.ListenAndServe(":8080", chain.ChainFunc(handlerFuncAtEdge))
*/
func (c *Chain) ChainFunc(f http.HandlerFunc) http.Handler {
	if f != nil {
		c.f = f
	}
	return c.Chain()
}
