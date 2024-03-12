package main

import (
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
)

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
	fmt.Println(r.URL.Path)
	// r.URL.Path is already cleaned by the http server, so there is no '../' in url
	// we can safely append it to the rootDir
	path := rootDir + r.URL.Path
	if strings.HasSuffix(path, "/") {
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
		tmpl := template.Must(template.New("").Parse(`
		<html>
		<head><style>body{font-size:xx-large}
		table{width:100%}
		</style></head>
		<body>
		<h1>Index of {{.Url}}</h1>
		<hr/>
		<table width="100%">
			{{if ne .Url "/"}}
			<tr>
				<td><a href="../">../</a></td>
				<td></td>
				<td></td>
			</tr>
			{{end}}
			{{range .Files}}
			<tr>
				<td><a href="{{.Path}}">{{.Name}}{{if .IsDir}}/{{end}}</a></td>
				<td>{{.ModTime.Format "2006-01-02 15:03:02"}}</td>
				<td>{{.Size}}</td>
			</tr>
			{{end}}
		</table>
		<hr/>
		</body>
		</html>
		`))

		err := tmpl.Execute(w, dir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		//TODO: do not redirect '/index.html' to '/'
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
	ipStr := flag.String("ip", "", "the ip to bind, leave")
	port := flag.Int("port", 8080, "the port to listen on")
	dir := flag.String("dir", "", "the directroy to serve, default is current directory")
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

	rootDir, _ = filepath.Abs(*dir)

	http.HandleFunc("/", logHandler(fileHandler))

	url := fmt.Sprintf("http://%s:%d", ip.String(), *port)

	fmt.Printf("Visit %s by clicking: %s\n", rootDir, url)
	fmt.Println("Or you can scan the qrcode below:")
	fmt.Println()

	q, _ := qrcode.New(url, qrcode.Highest)
	q.DisableBorder = true
	art := q.ToSmallString(true)
	fmt.Println(art)
	fmt.Println()

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
