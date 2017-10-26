package main

import (
	"runtime"

	"github.com/josephbateh/senior-project-server/authentication"
)

func main() {
	runtime.GOMAXPROCS(4)
	authentication.Start()
}
