package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveHandler(t *testing.T) {
	testcases := []struct {
		host               string
		expectedBody       string
		expectedStatusCode int
	}{
		{host: "localhost", expectedBody: "127.0.0.1", expectedStatusCode: http.StatusOK},
		{host: "foo.hoge", expectedBody: "Resolve Error", expectedStatusCode: http.StatusInternalServerError},
	}

	for _, testcase := range testcases {
		req, err := http.NewRequest("GET", "/"+testcase.host, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resolveHandler(w, r)
		})

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != testcase.expectedStatusCode {
			t.Errorf("Expected status code %d, but got %d", testcase.expectedStatusCode, status)
		}

		if rr.Body.String() != testcase.expectedBody {
			t.Errorf("Expected response body '%s', but got '%s'", testcase.expectedBody, rr.Body.String())
		}
	}
}
