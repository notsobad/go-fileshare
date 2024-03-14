# Go-FileShare: Your Go-To File Sharing Tool

Go-FileShare is a lightweight, user-friendly file sharing program written in Go. It enables you to swiftly share files from any directory on your system by initiating a web server and generating a QR code for effortless access.

## 🌟 Features

- **Ease of Use**: Simply navigate to the directory you wish to share and execute the command.
- **Portability**: Operates on any system with Go installed.
- **QR Code Generation**: Generates a QR code for seamless access from mobile devices.
- **Flexibility**: Allows you to specify the IP and port to bind to.

## 🚀 Installation

Install Go-FileShare with the following command:

```bash
go install github.com/notsobad/go-fileshare@latest
```

## 📖 Usage

Navigate to the directory you wish to share and run the `go-fileshare` command:

```bash
cd /path/to/directory
go-fileshare
```

This will initiate a web server and print a QR code to the terminal:

```bash
Visit /home/wang by clicking: http://10.0.1.68:8080
Or you can scan the qrcode below:
[[QRCODE HERE]]
```

You can now access the shared directory from any device by scanning the QR code or entering the URL into a web browser.

## ⚙️ Options

You can customize the behavior of Go-FileShare with the following options:

- `-dir`: The directory to serve. Default is the current directory.
- `-ip`: The IP to bind to.
- `-port`: The port to listen on. Default is 8080.
- `-auth`: Use basic auth to protect the service. like 'username:password'

For example, to serve a specific directory on a specific IP and port, you could run:

```bash
go-fileshare -dir /path/to/directory -ip 192.168.1.100 -port 8000 -auth admin:123456
```

## 🎉 Conclusion

Go-FileShare is a simple, portable solution for sharing files from any directory on your system. With its QR code generation feature, accessing shared files has never been easier.