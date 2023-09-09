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
		requestIpStr, _, err := net.SplitHostPort(r.RemoteAddr)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "IP Parse Error")
			log.Printf("ip: %s, block", r.RemoteAddr)
			return
		}

		requestIp := net.ParseIP(requestIpStr)

		if requestIp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "IP Parse Error")
			log.Printf("ip: %s, block", r.RemoteAddr)
			return
		}

		whitelist := loadAllowIpList()

		if !isAllowedIp(requestIp, whitelist) {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Fobidden")
			log.Printf("ip: %s, block", r.RemoteAddr)
			return
		}

		log.Printf("ip: %s, allow", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func isAllowedIp(ip net.IP, allowRanges []*net.IPNet) bool {
	for _, ipNet := range allowRanges {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

func loadAllowIpList() []*net.IPNet {
	envIpRanges := os.Getenv("ALLOW_IP_RANGES")

	var ipRanges []*net.IPNet

	for _, ipCidr := range strings.Split(envIpRanges, ",") {
		ipCidr = strings.TrimSpace(ipCidr)
		_, cidrNet, err := net.ParseCIDR(ipCidr)

		if err != nil {
			log.Println(err)
			log.Printf("invalid string: %s", ipCidr)
			continue
		}

		ipRanges = append(ipRanges, cidrNet)
	}

	return ipRanges
}
