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
	if _, err := w.Write([]byte("t1")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func test2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("t2")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestHandlerFuncWrapper(t *testing.T) {
	{
		h := &HandlerFuncWrapper{}
		assert.Nil(t, h.HandlerFunc)
	}
	{
		h := &HandlerFuncWrapper{
			HandlerFunc: test1,
		}
		f := http.HandlerFunc(test1)
		assert.Equal(t, reflect.ValueOf(f), reflect.ValueOf(h.HandlerFunc))
	}
}

func TestPreMiddleware(t *testing.T) {
	{
		h := &HandlerFuncWrapper{
			HandlerFunc: nil,
		}
		s := httptest.NewServer(h.PreMiddleware(http.HandlerFunc(test1)))
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
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.PreMiddleware(nil))
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
			HandlerFunc: nil,
		}
		s := httptest.NewServer(h.PreMiddleware(http.HandlerFunc(nil)))
		defer s.Close()

		_, err := http.Get(s.URL)
		assert.Error(t, err)
	}
	{
		h := &HandlerFuncWrapper{
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.PreMiddleware(http.HandlerFunc(test2)))
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
			HandlerFunc: nil,
		}
		s := httptest.NewServer(h.PostMiddleware(http.HandlerFunc(test1)))
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
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.PostMiddleware(nil))
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
			HandlerFunc: nil,
		}
		s := httptest.NewServer(h.PostMiddleware(http.HandlerFunc(nil)))
		defer s.Close()

		_, err := http.Get(s.URL)
		assert.Error(t, err)
	}
	{
		h := &HandlerFuncWrapper{
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.PostMiddleware(http.HandlerFunc(test2)))
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
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.Middleware(nil))
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
			HandlerFunc: nil,
		}
		s := httptest.NewServer(h.Middleware(http.HandlerFunc(test1)))
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
			HandlerFunc: test1,
		}
		s := httptest.NewServer(h.Middleware(http.HandlerFunc(test2)))
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
