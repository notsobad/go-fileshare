package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
)

func main() {
	port := flag.Int("port", 8080, "port number")

	flag.Parse()

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	obj := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlue, qrcodeTerminal.ConsoleColors.BrightGreen, qrcodeTerminal.QRCodeRecoveryLevels.Medium)
	fmt.Println("Open your browser to the address below, or you can scan the qrcode from your phone.")
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && ipnet.IP.To4() != nil {
			url := fmt.Sprintf("http://%s:%d", ipnet.IP.String(), *port)
			fmt.Println(url)

			obj.Get([]byte(url)).Print()
			fmt.Println()
		}
	}

	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
