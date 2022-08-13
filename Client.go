package main

import (
	"github.com/kxg3030/shermie-nat-traversal/ClientEndPoint"
	"github.com/kxg3030/shermie-nat-traversal/Log"
)

func init() {
	Log.NewLogger().Initialize()
}



func main() {
	c := ClientEndPoint.NewClient("127.0.0.1","9090")
	c.Run()
}
