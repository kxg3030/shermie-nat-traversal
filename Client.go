package main

import (
	"flag"
	"github.com/kxg3030/shermie-nat-traversal/ClientEndPoint"
	"github.com/kxg3030/shermie-nat-traversal/Log"
)

func init() {
	Log.NewLogger().Initialize()
}

func main() {
	host := flag.String("host", "127.0.0.1", "server host")
	port := flag.String("port", "9090", "listen port")
	pass := flag.String("pass", "123456", "auth password")
	bind := flag.String("bind", "127.0.0.1:80", "your local server host and port")
	flag.Parse()
	c := ClientEndPoint.NewClient(*host, *port, *pass,*bind)
	c.Run()
}
