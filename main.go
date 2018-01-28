package main

import (
	"flag"

	"github.com/msoedov/tcp-file/server"
)

var portPointer = flag.String("p", "3333", "Port")
var sourceFile = flag.String("f", "main.go", "Served file")

func main() {
	flag.Parse()
	hostAndPort := ":" + *portPointer
	server.Run(hostAndPort, *sourceFile)
}
