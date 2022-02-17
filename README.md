# A simple file share program

Just like `python -m SimpleHTTPServer`, it opens a simple web server and generate a QR code of the url in the terminal, you can get files from another computer of a phone.


```

$ go-fileshare -h
Usage of ./go-fileshare:
  -dir string
        the directroy to serve (default ".")
  -port int
        port number (default 8080)

$ go-fileshare -port 8083 -dir ../
Will serve /home/xxxx/scripts on http://10.0.2.1:8083

Or you can scan the QR code below on your phone:
█▀▀▀▀▀█ █ ▀▀█▀▀▀██▀ ▀ █▀▀▀▀▀█
█ ███ █ ▄▀▀██▀ ▄ ▀▀▄▄ █ ███ █
█ ▀▀▀ █ █▄▀   ▄ ▀██▄█ █ ▀▀▀ █
▀▀▀▀▀▀▀ ▀ ▀▄▀ ▀ ▀ ▀▄█ ▀▀▀▀▀▀▀
▄▄▄▀ ▄▀  ▄███ ▄▄▄▀ █ ▄▄█▀▀▄▀█
█▀█▄▄▀▀ █▀   ▄▄ ▀▄▀▄ ▀█▀ ▀▀█
 █▀▄ ▄▀▀▀▄ ▀▄▄█▀▄█ ▄▀▄▀▄ ▄ ▄▄
▄██   ▀▀▀▀▀█▄▄▄█▄▀ ▄  ▀█▄  ▄▀
  ██▀█▀ █▀▄▄ ▀ ▀ █ ▀▄▄█▄ ▄ ▄▄
▀ ▄▀▀▄▀▀█ ▄▀▄▀▀ ▀█▀▀▀▀ ▀ ▀ ██
▀  ▀▀ ▀ ▄█▄▄▀█▄ █▄█ █▀▀▀█ ▀▀
█▀▀▀▀▀█  ▀▀▄ ▀ ██▀▄▄█ ▀ █▀█▄▄
█ ███ █ ▄█ █ ▀ ▄█ ▄▄▀██▀█▀ █▄
█ ▀▀▀ █   █▀▄     ▀▄▄▀ █▀  ▄▀
▀▀▀▀▀▀▀      ▀▀▀ ▀ ▀ ▀  ▀▀▀▀

```
