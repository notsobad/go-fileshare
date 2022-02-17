package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"

	qrcode "github.com/skip2/go-qrcode"
)

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	port := flag.Int("port", 8080, "port number")
	dir := flag.String("dir", ".", "the directroy to serve")
	flag.Parse()

	ip := getOutboundIP()

	path, _ := filepath.Abs(*dir)
	http.Handle("/", http.FileServer(http.Dir(path)))

	url := fmt.Sprintf("http://%s:%d", ip.String(), *port)

	fmt.Printf("Will serve %s on %s\n", path, url)
	fmt.Println()

	fmt.Println("Or you can scan the qrcode below on your phone:")

	q, _ := qrcode.New(url, qrcode.Highest)
	q.DisableBorder = true
	art := q.ToSmallString(true)
	fmt.Println(art)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
