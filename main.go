package main

import (
	"runtime"

	"github.com/josephbateh/senior-project-server/server"
)

func main() {
	runtime.GOMAXPROCS(4)
	server.Start()
}
