package main

import (
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
