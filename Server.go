package main

import (
	"flag"
	"github.com/kxg3030/shermie-nat-traversal/Log"
	"github.com/kxg3030/shermie-nat-traversal/ServerEndpoint"
)

func init() {
	Log.NewLogger().Initialize()
}

func main() {
	port := flag.String("port", ":9090", "listen port")
	pass := flag.String("pass", "123456", "auth password")
	flag.Parse()
	s := ServerEndpoint.NewServer(*port, *pass)
	s.Run()
}
