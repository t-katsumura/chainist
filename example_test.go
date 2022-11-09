package chainist

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func myHandlerFunc0(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hi from myHandlerFunc0!\n")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func myHandlerFunc1(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hi from myHandlerFunc1!\n")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func myHandlerFunc2(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hi from myHandlerFunc2!\n")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func myHandler1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hi from myHandler1!\n")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if next != nil {
			next.ServeHTTP(w, r)

		}
	})
}

func myHandler2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hi from myHandler2!\n")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}

func ExampleChain() {

	// create a new chain struct
	chain := NewChain()

	// add handler functions
	chain.AppendPreFunc(myHandlerFunc1)
	chain.AppendPostFunc(myHandlerFunc2)
	// add handlers
	chain.Extend(myHandler1, myHandler2)
	chain.SetHandlerFunc(myHandlerFunc0)

	// Run http server
	ts := httptest.NewServer(chain.Chain())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		// handle error
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		// handle error
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	resp.Body.Close()
	s := string(body)
	fmt.Println(s)

	// Output:
	// Hi from myHandlerFunc1!
	// Hi from myHandler1!
	// Hi from myHandler2!
	// Hi from myHandlerFunc0!
	// Hi from myHandlerFunc2!
}
