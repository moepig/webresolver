package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func resolveHandler(w http.ResponseWriter, r *http.Request) {
	host := strings.TrimPrefix(r.URL.Path, "/")
	log.Printf("request ip: %s, host: %s\n", r.RemoteAddr, host)

	addr, err := resolve(host)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Resolve Error")
		return
	}

	if addr == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unknown Error")
		return
	}

	fmt.Fprintf(w, "%s", addr.String())
}
