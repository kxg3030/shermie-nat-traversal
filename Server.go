package main

import (
	"github.com/kxg3030/shermie-nat-traversal/Log"
	"github.com/kxg3030/shermie-nat-traversal/ServerEndpoint"
)

func init() {
	Log.NewLogger().Initialize()
}

func main() {
	s := ServerEndpoint.NewServer(":9090", "123456")
	s.Run()
}
