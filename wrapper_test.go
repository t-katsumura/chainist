package chainist

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func test1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("t1"))
}

func test2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("t2"))
}

func TestHandlerFuncWrapper(t *testing.T) {
	{
		h := &HandlerFuncWrapper{}
		assert.Nil(t, h.f)
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		f := http.HandlerFunc(test1)
		assert.Equal(t, reflect.ValueOf(f), reflect.ValueOf(h.f))
	}
}

func TestPreMiddleware(t *testing.T) {
	{
		h := &HandlerFuncWrapper{
			f: nil,
		}
		s := httptest.NewServer(h.preMiddleware(http.HandlerFunc(test1)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.preMiddleware(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: nil,
		}
		s := httptest.NewServer(h.preMiddleware(http.HandlerFunc(nil)))
		defer s.Close()

		_, err := http.Get(s.URL)
		assert.Error(t, err)
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.preMiddleware(http.HandlerFunc(test2)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1t2", string(body))
	}
}

func TestPostMiddleware(t *testing.T) {
	{
		h := &HandlerFuncWrapper{
			f: nil,
		}
		s := httptest.NewServer(h.postMiddleware(http.HandlerFunc(test1)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.postMiddleware(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: nil,
		}
		s := httptest.NewServer(h.postMiddleware(http.HandlerFunc(nil)))
		defer s.Close()

		_, err := http.Get(s.URL)
		assert.Error(t, err)
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.postMiddleware(http.HandlerFunc(test2)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t2t1", string(body))
	}
}

func TestMiddleware(t *testing.T) {
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.middleware(nil))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: nil,
		}
		s := httptest.NewServer(h.middleware(http.HandlerFunc(test1)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1", string(body))
	}
	{
		h := &HandlerFuncWrapper{
			f: test1,
		}
		s := httptest.NewServer(h.middleware(http.HandlerFunc(test2)))
		defer s.Close()

		res, err := http.Get(s.URL)
		assert.NoError(t, err)

		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		assert.NoError(t, err)

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, "t1t2", string(body))
	}
}
