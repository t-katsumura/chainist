package chainist

import (
	"net/http"
)

// define middleware function type
type middleware func(h http.Handler) http.Handler

type Chain struct {
	// fs : middlewares
	fs []middleware
	// f : handler function run at the edge of the chain
	f http.HandlerFunc
}

// Create the new chain struct
func NewChain(fs ...middleware) *Chain {
	c := &Chain{
		fs: fs,
	}
	return c
}

// Append middleware to the chain
func (c *Chain) Append(f middleware) *Chain {
	if f == nil {
		return c
	}
	c.fs = append(c.fs, f)
	return c
}

// Append http handler function as pre-executable function
// which is executed before calling succeeding middlewares.
func (c *Chain) AppendPreFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Append(h.preMiddleware)
}

// Append http handler function as post-executable function
// which is executed after calling succeeding middlewares.
func (c *Chain) AppendPostFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Append(h.postMiddleware)
}

// Insert a middleware at designated number of chain
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

// Insert a pre-executable handler function at designated number of chain
func (c *Chain) InsertPreFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Insert(h.preMiddleware, i)
}

// Insert a post-executable handler function at designated number of chain
func (c *Chain) InsertPostFunc(f http.HandlerFunc, i int) *Chain {
	if f == nil {
		return c
	}
	h := &HandlerFuncWrapper{f: f}
	return c.Insert(h.postMiddleware, i)
}

// append multiple middlewares at a time
func (c *Chain) Extend(fs ...middleware) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		c.Append(f)
	}
	return c
}

// append multiple pre-executable handler functions at a time
func (c *Chain) ExtendPreFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{f: f}
		c.Append(h.preMiddleware)
	}
	return c
}

// append multiple post-executable handler functions at a time
func (c *Chain) ExtendPostFunc(fs ...http.HandlerFunc) *Chain {
	for _, f := range fs {
		if f == nil {
			continue
		}
		h := &HandlerFuncWrapper{f: f}
		c.Append(h.postMiddleware)
	}
	return c
}

// set the htt handler function
func (c *Chain) SetHandlerFunc(f http.HandlerFunc) *Chain {
	if f == nil {
		return c
	}
	c.f = f
	return c
}

// join the two chains
func (c *Chain) Join(o *Chain) *Chain {
	if o == nil {
		return c
	}
	c.fs = append(c.fs, o.fs...)
	return c
}

// get the length of middleware chain
func (c *Chain) Len() int {
	return len(c.fs)
}

// return the new middeware chain
func (c *Chain) Chain() http.Handler {

	var h http.Handler
	if c.f != nil {
		h = c.f
	} else {
		h = (&HandlerFuncWrapper{f: c.f}).middleware(nil)
	}

	n := c.Len()
	for i := range c.fs {
		h = c.fs[n-1-i](h)
	}

	return h
}

// return the new middeware chain
func (c *Chain) ChainFunc(h http.HandlerFunc) http.Handler {
	if h != nil {
		c.f = h
	}
	return c.Chain()
}
