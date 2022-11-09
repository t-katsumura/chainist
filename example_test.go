package chainist

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func handlerFunc11(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi from handlerFunc 1!\n"))
}

func handlerFunc22(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi from handlerFunc 2!\n"))
}

func handler11(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte(r.URL.Path + "\n"))
		// w.Write([]byte(r.URL.RawQuery + "\n"))
		// v := r.URL.Query()
		// for key, vs := range v {
		// 	fmt.Fprintf(w, "%s = %s\n", key, vs[0])
		// }
		w.Write([]byte("Hi from handler 1!\n"))
		if next != nil {
			next.ServeHTTP(w, r)

		}
	})
}

func handler22(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi from handler 2!\n"))
		next.ServeHTTP(w, r)
	})
}

func ExampleChain() {

	// create a new chain struct
	chain := NewChain()

	// add handler functions
	// chain.AppendPreFunc(handlerFunc1)
	// chain.AppendPostFunc(handlerFunc2)
	// add handlers
	// chain.Extend(handler11, handler22)
	chain.SetHandlerFunc(handlerFunc11)

	// Run http server which responds
	/*
		Hi from handlerFunc 1!
		Hi from handler 1!
		Hi from handler 2!
		Hi from handlerFunc 2!
	*/
	// http.Handle("/", chain.Chain())
	// println(chain.Len())
	ts := httptest.NewServer(chain.Chain())
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, _ := (&http.Client{}).Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	s := string(body)
	fmt.Println(s)

	// Output: Hi from handlerFunc 1!
}

func ExampleChain_AppendPreFunc() {
	// create a new chain struct
	chain := NewChain()

	// add handler functions
	chain.AppendPreFunc(handlerFunc1)
	chain.AppendPostFunc(handlerFunc2)
	// add handlers
	chain.Extend(handler11, handler22)
	chain.SetHandlerFunc(handlerFunc1)

	// Run http server which responds
	/*
		Hi from handlerFunc 1!
		Hi from handler 1!
		Hi from handler 2!
		Hi from handlerFunc 2!
	*/
	// http.Handle("/", chain.Chain())
	println(chain.Len())
	http.ListenAndServe(":8080", chain.Chain())
}
