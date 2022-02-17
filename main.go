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

type CustomResponseWriter struct {
	w      http.ResponseWriter
	Code   int
	Length int
}

func NewCustomResponseWriter(ww http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		w:      ww,
		Code:   0,
		Length: 0,
	}
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if w.Code == 0 {
		w.Code = 200
	}
	n, err := w.w.Write(b)
	w.Length += n
	return n, err
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.w.WriteHeader(statusCode)
}

//See: https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w2 := NewCustomResponseWriter(w)
		handler.ServeHTTP(w2, r)
		log.Printf("%s %s %s %d %s %d\n", r.RemoteAddr, r.Method, r.URL, w2.Code, r.Method, w2.Length)
	})
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

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), logRequest(http.DefaultServeMux)))
}
