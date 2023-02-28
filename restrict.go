package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func ipRestrictMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Parse Error")
			log.Printf("ip: %s, block", r.RemoteAddr)
			return
		}

		if !isAllowedIp(ip) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Fobidden")
			log.Printf("ip: %s, block", r.RemoteAddr)
			return
		}

		log.Printf("ip: %s, allow", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func isAllowedIp(ip string) bool {
	envIpRanges := os.Getenv("ALLOW_IP_RANGES")
	targetIp := net.ParseIP(ip)

	if envIpRanges == "" {
		return false
	}

	for _, ipCidr := range strings.Split(envIpRanges, ",") {
		_, cidrNet, err := net.ParseCIDR(ipCidr)

		if err != nil {
			log.Println(err)
			return false
		}

		return cidrNet.Contains(targetIp)
	}

	return false
}
