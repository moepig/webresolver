package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIpRestrictMiddleware(t *testing.T) {
	testcases := []struct {
		allowIPRanges string
		requestAddr   string
		expected      int
	}{
		{allowIPRanges: "192.168.1.0/24, 10.0.0.0/16", requestAddr: "192.168.1.0:12345", expected: http.StatusOK},
		{allowIPRanges: "192.168.1.0/24, 10.0.0.0/16", requestAddr: "10.0.100.0:111", expected: http.StatusOK},
		{allowIPRanges: "192.168.1.0/24, 10.0.0.0/16", requestAddr: "10.100.100.0:111", expected: http.StatusForbidden},
		{allowIPRanges: "192.168.1.0/24, 10.0.0.0/16", requestAddr: "foo:111", expected: http.StatusInternalServerError},
	}

	for _, testcase := range testcases {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		os.Setenv("ALLOW_IP_RANGES", testcase.allowIPRanges)
		req.RemoteAddr = testcase.requestAddr

		rr := httptest.NewRecorder()

		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		middleware := ipRestrictMiddleware(testHandler)

		middleware.ServeHTTP(rr, req)

		if status := rr.Code; status != testcase.expected {
			t.Logf("testcase: %#v", testcase)
			t.Errorf("Expected status code %d, but got %d", testcase.expected, status)
		}
	}
}
