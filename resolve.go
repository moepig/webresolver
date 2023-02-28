package main

import (
	"log"
	"net"
)

func resolve(host string) (*net.IPAddr, error) {
	addr, err := net.ResolveIPAddr("ip", host)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return addr, nil
}
