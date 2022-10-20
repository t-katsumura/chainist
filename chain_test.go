package chainist

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func handlerFunc1(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("f1")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handlerFunc2(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("f2")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

var handler1 = &HandlerFuncWrapper{
	HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("h1")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	},
}

var handler2 = &HandlerFuncWrapper{
	HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("h2")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	},
}

func funcPointer(f any) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func TestChainStruct(t *testing.T) {
	{
		c := &Chain{}
		e := &Chain{
			Middleware:  nil,
			HandlerFunc: nil,
		}
		assert.Equal(t, e, c)
	}
	{
		c := &Chain{
			Middleware:  []Middleware{handler1.PreMiddleware},
			HandlerFunc: handler2.HandlerFunc,
		}
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.PreMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.HandlerFunc), funcPointer(c.HandlerFunc))
	}
}

func TestNewChain(t *testing.T) {
	{
		c := NewChain()
		assert.Equal(t, 0, len(c.Middleware))
		assert.Nil(t, c.HandlerFunc)
	}
	{
		c := NewChain(nil)
		assert.Nil(t, c)
	}
	{
		c := NewChain(handler1.Middleware)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Nil(t, c.HandlerFunc)
	}
	{
		c := NewChain(handler1.Middleware, handler2.Middleware)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[1]))
		assert.Nil(t, c.HandlerFunc)
	}
}

func TestAppend(t *testing.T) {
	{
		c := NewChain()
		c.Append(nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.Append(handler1.Middleware)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.Append(handler1.Middleware)
		c.Append(handler2.Middleware)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[1]))
	}
}

func TestAppendPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.AppendPreFunc(nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.AppendPreFunc(handlerFunc1)
		e := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e.PreMiddleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.AppendPreFunc(handlerFunc1)
		c.AppendPreFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PreMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PreMiddleware), funcPointer(c.Middleware[1]))
	}
}

func TestAppendPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.AppendPostFunc(nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.AppendPostFunc(handlerFunc1)
		e := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e.PostMiddleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.AppendPostFunc(handlerFunc1)
		c.AppendPostFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PostMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PostMiddleware), funcPointer(c.Middleware[1]))
	}
}

func TestInsert(t *testing.T) {
	{
		c := NewChain()
		c.Insert(nil, 0)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.Insert(handler1.Middleware, 0)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.Middleware, -99)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain().Append(handler2.Middleware)
		c.Insert(handler1.Middleware, -99)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.Middleware, 99)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain().Append(handler2.Middleware)
		c.Insert(handler1.Middleware, 99)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.Middleware, 0)
		c.Insert(handler2.Middleware, 1)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[1]))
	}
}

func TestInsertPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.InsertPreFunc(nil, 0)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 0)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PreMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, -99)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PreMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 99)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PreMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 0)
		c.InsertPreFunc(handlerFunc2, 0)
		e1 := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PreMiddleware
		e2 := (&HandlerFuncWrapper{HandlerFunc: handlerFunc2}).PreMiddleware
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2), funcPointer(c.Middleware[1]))
	}
}

func TestInsertPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.InsertPostFunc(nil, 0)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 0)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PostMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, -99)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PostMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 99)
		e := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PostMiddleware
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 0)
		c.InsertPostFunc(handlerFunc2, 0)
		e1 := (&HandlerFuncWrapper{HandlerFunc: handlerFunc1}).PostMiddleware
		e2 := (&HandlerFuncWrapper{HandlerFunc: handlerFunc2}).PostMiddleware
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2), funcPointer(c.Middleware[1]))
	}
}

func TestExtend(t *testing.T) {
	{
		c := NewChain()
		c.Extend(nil, nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.Extend(handler1.Middleware)
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.Extend(handler1.Middleware)
		c.Extend(handler2.Middleware)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[1]))
	}
	{
		c := NewChain()
		c.Extend(handler1.Middleware, handler2.Middleware)
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c.Middleware[1]))
	}
}

func TestExtendPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.ExtendPreFunc(nil, nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1)
		e := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e.PreMiddleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1)
		c.ExtendPreFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PreMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PreMiddleware), funcPointer(c.Middleware[1]))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1, handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PreMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PreMiddleware), funcPointer(c.Middleware[1]))
	}
}

func TestExtendPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.ExtendPostFunc(nil, nil)
		assert.Equal(t, 0, len(c.Middleware))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1)
		e := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		assert.Equal(t, 1, len(c.Middleware))
		assert.Equal(t, funcPointer(e.PostMiddleware), funcPointer(c.Middleware[0]))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1)
		c.ExtendPostFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PostMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PostMiddleware), funcPointer(c.Middleware[1]))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1, handlerFunc2)
		e1 := &HandlerFuncWrapper{HandlerFunc: handlerFunc1}
		e2 := &HandlerFuncWrapper{HandlerFunc: handlerFunc2}
		assert.Equal(t, 2, len(c.Middleware))
		assert.Equal(t, funcPointer(e1.PostMiddleware), funcPointer(c.Middleware[0]))
		assert.Equal(t, funcPointer(e2.PostMiddleware), funcPointer(c.Middleware[1]))
	}
}

func TestSetHandlerFunc(t *testing.T) {
	{
		c := NewChain()
		c.SetHandlerFunc(nil)
		assert.Nil(t, c.HandlerFunc)
	}
	{
		c := NewChain()
		c.SetHandlerFunc(handlerFunc1)
		assert.Equal(t, funcPointer(handlerFunc1), funcPointer(c.HandlerFunc))
	}
}

func TestJoin(t *testing.T) {
	{
		c1 := NewChain()
		c1.Join(nil)
		assert.Equal(t, 0, len(c1.Middleware))
	}
	{
		c1 := NewChain()
		c2 := NewChain()
		c1.Join(c2)
		assert.Equal(t, 0, len(c1.Middleware))
	}
	{
		c1 := NewChain().Append(handler1.Middleware)
		c1.Join(nil)
		assert.Equal(t, 1, len(c1.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c1.Middleware[0]))
	}
	{
		c1 := NewChain().Append(handler1.Middleware)
		c2 := NewChain().Append(handler2.Middleware)
		c1.Join(c2)
		assert.Equal(t, 2, len(c1.Middleware))
		assert.Equal(t, funcPointer(handler1.Middleware), funcPointer(c1.Middleware[0]))
		assert.Equal(t, funcPointer(handler2.Middleware), funcPointer(c1.Middleware[1]))
	}
}

func TestLen(t *testing.T) {
	{
		c := NewChain()
		assert.Equal(t, 0, c.Len())
	}
	{
		c := NewChain().Append(handler1.Middleware).Append(handler2.Middleware)
		assert.Equal(t, 2, c.Len())
	}
}

func TestChain(t *testing.T) {
	{
		c := NewChain()
		s := httptest.NewServer(c.Chain())
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 404, res.StatusCode)
		assert.NotEmpty(t, body)
	}
	{
		c := NewChain().Append(handler1.Middleware)
		s := httptest.NewServer(c.Chain())
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "h1", string(body))
	}
	{
		c := NewChain().Append(handler1.Middleware).Append(handler2.Middleware)
		s := httptest.NewServer(c.Chain())
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "h1h2", string(body))
	}
	{
		c := NewChain().AppendPreFunc(handlerFunc1).AppendPreFunc(handlerFunc2)
		s := httptest.NewServer(c.Chain())
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f1f2", string(body))
	}
	{
		c := NewChain().AppendPostFunc(handlerFunc1).AppendPostFunc(handlerFunc2)
		s := httptest.NewServer(c.Chain())
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f2f1", string(body))
	}
}

func TestChainFunc(t *testing.T) {
	{
		c := NewChain()
		s := httptest.NewServer(c.ChainFunc(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 404, res.StatusCode)
		assert.NotEmpty(t, body)
	}
	{
		c := NewChain()
		s := httptest.NewServer(c.ChainFunc(handlerFunc1))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f1", string(body))
	}
	{
		c := NewChain().Append(handler1.Middleware).Append(handler2.Middleware)
		s := httptest.NewServer(c.ChainFunc(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "h1h2", string(body))
	}
	{
		c := NewChain().Append(handler1.Middleware).Append(handler2.Middleware)
		s := httptest.NewServer(c.ChainFunc(handlerFunc1))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "h1h2f1", string(body))
	}
	{
		c := NewChain().AppendPreFunc(handlerFunc1).AppendPreFunc(handlerFunc2)
		s := httptest.NewServer(c.ChainFunc(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f1f2", string(body))
	}
	{
		c := NewChain().AppendPreFunc(handlerFunc2).AppendPreFunc(handlerFunc2)
		s := httptest.NewServer(c.ChainFunc(handlerFunc1))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f2f2f1", string(body))
	}
	{
		c := NewChain().AppendPostFunc(handlerFunc1).AppendPostFunc(handlerFunc2)
		s := httptest.NewServer(c.ChainFunc(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f2f1", string(body))
	}
	{
		c := NewChain().AppendPostFunc(handlerFunc2).AppendPostFunc(handlerFunc2)
		s := httptest.NewServer(c.ChainFunc(handlerFunc1))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "f1f2f2", string(body))
	}
}
