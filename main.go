package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Faizal-Asep/crud-app/service"
	"github.com/joho/godotenv"
)

func init() {
	if len(os.Args) > 1 {
		os.Setenv("ENV_FILE", filepath.Join("./config", os.Args[1]))
	} else {
		os.Setenv("ENV_FILE", filepath.Join("./config", ".env"))
	}

	err := godotenv.Load(os.Getenv("ENV_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	a := service.App{}
	a.Initialize()

	a.Run("8081")
	os.Exit(0)
}
