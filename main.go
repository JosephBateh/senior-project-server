package main

import (
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/josephbateh/senior-project-server/server"
)

func main() {
	runtime.GOMAXPROCS(4)
	loadEnv()
	server.Start()
}

func loadEnv() {
	if os.Getenv("MLAB_DB") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
