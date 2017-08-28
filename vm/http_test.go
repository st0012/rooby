package vm

import (
	"testing"
	//"net/http/httptest"
	//"net/http"
	"net/http"
	"io/ioutil"
)

func TestHTTPObject(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		//test get request
		{`
		require "net/http"

		Net::HTTP.get("http://127.0.0.1:3000")
		`, "GET Hello World"},
		{`
		require "net/http"

		Net::HTTP.post("http://127.0.0.1:3000", "text/plain", "Hi Again")
		`, "POST Hi Again"},
	}

	c := make(chan bool, 1)

	//server to test off of
	go func() {
		m := http.NewServeMux()

		m.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if r.Method == http.MethodPost {
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.Write([]byte(r.Method + " " + string(b)))
			} else {
				w.Write([]byte(r.Method + " Hello World"))
			}

		})

		c <- true

		http.ListenAndServe(":3000", m)
	}()

	//block until server is ready
	<- c

	for i, tt := range tests {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkExpected(t, i, evaluated, tt.expected)
		v.checkCFP(t, i, 0)
		v.checkSP(t, i, 1)
	}
}
