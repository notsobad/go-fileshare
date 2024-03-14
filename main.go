package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

var (
	rootDir string
	//go:embed index.html
	indexHTML embed.FS
)

type Auth struct {
	Username string
	Password string
}
type File struct {
	Path    string
	Name    string
	ModTime time.Time
	Size    int64
	IsDir   bool
}

type Directory struct {
	Url   string
	Files []File
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	//Do not cache
	r.Header.Del("If-Modified-Since")
	r.Header.Del("If-None-Match")

	// Disable Client Cache
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// r.URL.Path is already cleaned by the http server, so there is no '../' in url
	// we can safely append it to the rootDir
	path := rootDir + r.URL.Path
	if strings.HasSuffix(path, "/") {
		// if the path is a directory, list the files and directories in it

		files, _ := filepath.Glob(path + "*")
		// write the index page, list file and directorys, show filename, last modified time, size in table

		var dir Directory
		dir.Url = r.URL.Path
		for _, file := range files {
			_, filename := filepath.Split(file)
			if strings.HasPrefix(filename, ".") {
				continue
			}
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			dir.Files = append(dir.Files, File{
				Path:    r.URL.Path + filename,
				Name:    filename,
				ModTime: info.ModTime(),
				Size:    info.Size(),
				IsDir:   info.IsDir(),
			})
		}
		tmpl, _ := template.ParseFS(indexHTML, "index.html")

		err := tmpl.Execute(w, dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if strings.HasSuffix(path, "index.html") {
		// Avoid redirecting '/index.html' to '/'
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, path, time.Now(), f)
	} else {
		// Other files, just serve it
		http.ServeFile(w, r, path)
	}
}

func logHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		h(lrw, r)
		end := time.Now()

		fmt.Printf("%s %s %s %d %s %d\n",
			end.Format("2006/01/02 15:04:05"),
			r.RemoteAddr,
			r.Method,
			lrw.statusCode,
			r.URL.Path,
			end.Sub(start).Microseconds(),
		)
	}
}

func basicAuthHandler(h http.HandlerFunc, auth Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != auth.Username || pass != auth.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {
	ipStr := flag.String("ip", "", "the ip to bind, default is the outbound ip address")
	port := flag.Int("port", 8080, "the port to listen on")
	dir := flag.String("dir", "", "the directroy to serve, default is current directory")
	authStr := flag.String("auth", "", "Basic auth credentials in the format 'username:password'")

	flag.Parse()

	if *dir == "" {
		*dir, _ = os.Getwd()
	}

	var ip net.IP

	if *ipStr != "" {
		ip = net.ParseIP(*ipStr)
		if ip == nil {
			fmt.Printf("%s is not a valid ip address\n", ip)
			os.Exit(1)
		}
	} else {
		ip = getOutboundIP()
	}

	url := fmt.Sprintf("http://%s:%d", ip.String(), *port)
	rootDir, _ = filepath.Abs(*dir)

	fmt.Printf("Visit %s by clicking: %s\n", rootDir, url)
	fmt.Println("Or you can scan the qrcode below:")
	fmt.Println()

	q, _ := qrcode.New(url, qrcode.Highest)
	q.DisableBorder = true
	art := q.ToSmallString(true)
	fmt.Println(art)
	fmt.Println()

	if *authStr != "" {
		parts := strings.SplitN(*authStr, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid auth credentials")
			os.Exit(1)
		}
		auth := Auth{
			Username: parts[0],
			Password: parts[1],
		}
		http.HandleFunc("/", basicAuthHandler(logHandler(fileHandler), auth))
	} else {
		http.HandleFunc("/", logHandler(fileHandler))
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
