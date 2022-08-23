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
	w.Write([]byte("f1"))
}

func handlerFunc2(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("f2"))
}

var handler1 = &HandlerFuncWrapper{
	f: func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("h1"))
	},
}

var handler2 = &HandlerFuncWrapper{
	f: func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("h2"))
	},
}

func funcPointer(f any) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func TestChainStruct(t *testing.T) {
	{
		c := &Chain{}
		e := &Chain{
			fs: nil,
			f:  nil,
		}
		assert.Equal(t, e, c)
	}
	{
		c := &Chain{
			fs: []middleware{handler1.preMiddleware},
			f:  handler2.f,
		}
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.preMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.f), funcPointer(c.f))
	}
}

func TestNewChain(t *testing.T) {
	{
		c := NewChain()
		assert.Equal(t, 0, len(c.fs))
		assert.Nil(t, c.f)
	}
	{
		c := NewChain(handler1.middleware)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Nil(t, c.f)
	}
	{
		c := NewChain(handler1.middleware, handler2.middleware)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[1]))
		assert.Nil(t, c.f)
	}
}

func TestAppend(t *testing.T) {
	{
		c := NewChain()
		c.Append(nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.Append(handler1.middleware)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.Append(handler1.middleware)
		c.Append(handler2.middleware)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[1]))
	}
}

func TestAppendPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.AppendPreFunc(nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.AppendPreFunc(handlerFunc1)
		e := &HandlerFuncWrapper{f: handlerFunc1}
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e.preMiddleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.AppendPreFunc(handlerFunc1)
		c.AppendPreFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.preMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.preMiddleware), funcPointer(c.fs[1]))
	}
}

func TestAppendPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.AppendPostFunc(nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.AppendPostFunc(handlerFunc1)
		e := &HandlerFuncWrapper{f: handlerFunc1}
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e.postMiddleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.AppendPostFunc(handlerFunc1)
		c.AppendPostFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.postMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.postMiddleware), funcPointer(c.fs[1]))
	}
}

func TestInsert(t *testing.T) {
	{
		c := NewChain()
		c.Insert(nil, 0)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.Insert(handler1.middleware, 0)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.middleware, -99)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain().Append(handler2.middleware)
		c.Insert(handler1.middleware, -99)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.middleware, 99)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain().Append(handler2.middleware)
		c.Insert(handler1.middleware, 99)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.Insert(handler1.middleware, 0)
		c.Insert(handler2.middleware, 1)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[1]))
	}
}

func TestInsertPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.InsertPreFunc(nil, 0)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 0)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).preMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, -99)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).preMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 99)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).preMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPreFunc(handlerFunc1, 0)
		c.InsertPreFunc(handlerFunc2, 0)
		e1 := (&HandlerFuncWrapper{f: handlerFunc1}).preMiddleware
		e2 := (&HandlerFuncWrapper{f: handlerFunc2}).preMiddleware
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2), funcPointer(c.fs[1]))
	}
}

func TestInsertPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.InsertPostFunc(nil, 0)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 0)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).postMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, -99)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).postMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 99)
		e := (&HandlerFuncWrapper{f: handlerFunc1}).postMiddleware
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.InsertPostFunc(handlerFunc1, 0)
		c.InsertPostFunc(handlerFunc2, 0)
		e1 := (&HandlerFuncWrapper{f: handlerFunc1}).postMiddleware
		e2 := (&HandlerFuncWrapper{f: handlerFunc2}).postMiddleware
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2), funcPointer(c.fs[1]))
	}
}

func TestExtend(t *testing.T) {
	{
		c := NewChain()
		c.Extend(nil, nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.Extend(handler1.middleware)
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.Extend(handler1.middleware)
		c.Extend(handler2.middleware)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[1]))
	}
	{
		c := NewChain()
		c.Extend(handler1.middleware, handler2.middleware)
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c.fs[1]))
	}
}

func TestExtendPreFunc(t *testing.T) {
	{
		c := NewChain()
		c.ExtendPreFunc(nil, nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1)
		e := &HandlerFuncWrapper{f: handlerFunc1}
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e.preMiddleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1)
		c.ExtendPreFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.preMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.preMiddleware), funcPointer(c.fs[1]))
	}
	{
		c := NewChain()
		c.ExtendPreFunc(handlerFunc1, handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.preMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.preMiddleware), funcPointer(c.fs[1]))
	}
}

func TestExtendPostFunc(t *testing.T) {
	{
		c := NewChain()
		c.ExtendPostFunc(nil, nil)
		assert.Equal(t, 0, len(c.fs))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1)
		e := &HandlerFuncWrapper{f: handlerFunc1}
		assert.Equal(t, 1, len(c.fs))
		assert.Equal(t, funcPointer(e.postMiddleware), funcPointer(c.fs[0]))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1)
		c.ExtendPostFunc(handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.postMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.postMiddleware), funcPointer(c.fs[1]))
	}
	{
		c := NewChain()
		c.ExtendPostFunc(handlerFunc1, handlerFunc2)
		e1 := &HandlerFuncWrapper{f: handlerFunc1}
		e2 := &HandlerFuncWrapper{f: handlerFunc2}
		assert.Equal(t, 2, len(c.fs))
		assert.Equal(t, funcPointer(e1.postMiddleware), funcPointer(c.fs[0]))
		assert.Equal(t, funcPointer(e2.postMiddleware), funcPointer(c.fs[1]))
	}
}

func TestSetHandlerFunc(t *testing.T) {
	{
		c := NewChain()
		c.SetHandlerFunc(nil)
		assert.Nil(t, c.f)
	}
	{
		c := NewChain()
		c.SetHandlerFunc(handlerFunc1)
		assert.Equal(t, funcPointer(handlerFunc1), funcPointer(c.f))
	}
}

func TestJoin(t *testing.T) {
	{
		c1 := NewChain()
		c1.Join(nil)
		assert.Equal(t, 0, len(c1.fs))
	}
	{
		c1 := NewChain()
		c2 := NewChain()
		c1.Join(c2)
		assert.Equal(t, 0, len(c1.fs))
	}
	{
		c1 := NewChain().Append(handler1.middleware)
		c1.Join(nil)
		assert.Equal(t, 1, len(c1.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c1.fs[0]))
	}
	{
		c1 := NewChain().Append(handler1.middleware)
		c2 := NewChain().Append(handler2.middleware)
		c1.Join(c2)
		assert.Equal(t, 2, len(c1.fs))
		assert.Equal(t, funcPointer(handler1.middleware), funcPointer(c1.fs[0]))
		assert.Equal(t, funcPointer(handler2.middleware), funcPointer(c1.fs[1]))
	}
}

func TestLen(t *testing.T) {
	{
		c := NewChain()
		assert.Equal(t, 0, c.Len())
	}
	{
		c := NewChain().Append(handler1.middleware).Append(handler2.middleware)
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

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "", string(body))
	}
	{
		c := NewChain().Append(handler1.middleware)
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
		c := NewChain().Append(handler1.middleware).Append(handler2.middleware)
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

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "", string(body))
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
		c := NewChain().Append(handler1.middleware).Append(handler2.middleware)
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
		c := NewChain().Append(handler1.middleware).Append(handler2.middleware)
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
