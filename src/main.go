package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	workerMain()
}

func workerMain() {
	beanstalkServerAddress := os.Getenv("BEANSTALK_SERVER")
	beanstalkTube := os.Getenv("BEANSTALK_TUBE")

	worker := makeNewWorker(beanstalkServerAddress, beanstalkTube)
	worker.Connect()
	defer worker.Close()

	for {
		worker.ProcessJob()
	}
}
