package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {

	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ipMap := make(map[string][]string)

	for _, iface := range ifaces {

		addrs, err := iface.Addrs()
		if err != nil {
			// Skip error interface
			continue
		}

		ips := make([]string, 0)

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}

		if len(ips) > 0 {
			ipMap[iface.Name] = ips
		}
	}

	res, _ := json.Marshal(ipMap)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func main() {
	portPtr := flag.Int("p", 8080, "Port")
	hostPtr := flag.String("h", "0.0.0.0", "Host")
	flag.Parse()

	http.HandleFunc("/", handler)

	listenAddr := fmt.Sprintf("%v:%v", *hostPtr, *portPtr)

	fmt.Printf("Listening on %v...\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))

}
