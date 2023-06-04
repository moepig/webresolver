package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	envIpRanges := os.Getenv("ALLOW_IP_RANGES")
	for _, ipRange := range strings.Split(envIpRanges, ",") {
		log.Printf("allow ip range from env: %s", ipRange)
	}

	http.HandleFunc("/", ipRestrictMiddleware(resolveHandler))

	log.Printf("start server")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
