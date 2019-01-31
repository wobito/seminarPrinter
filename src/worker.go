package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Worker struct {
	MasterBeanstalk
	tube       string
	storageDir string
}

func (worker *Worker) ProcessJob() {
	id, body, err := worker.MasterBeanstalk.tubeSet.Reserve(5 * time.Hour)

	if err != nil {
		panic(err)
	}

	fmt.Println("Received Job:", id)

	fmt.Println("-- Decoding Job Data...")
	prtData := Job{}
	json.Unmarshal([]byte(string(body)), &prtData)
	fmt.Println("-- Decoded Data:", prtData)

	err = worker.DownloadFile(prtData.Data)
	if err != nil {
		panic(err)
	}
	fmt.Println("-- Download Complete")
	fmt.Println("-- Sending To Printer: ", prtData.Data.Printer)

	worker.PrintJob(prtData.Data)

	fmt.Println("-- Printing... ", prtData.Data.Printer)

	worker.MasterBeanstalk.serverConnection.Delete(id)
	fmt.Println("Job Released:", id)
}

func (worker *Worker) PrintJob(data JobData) {
	printCmd := "/usr/bin/lp"
	params := "-d"

	cmd := exec.Command(printCmd, params, data.Printer, worker.storageDir+data.File)
	stdout, err := cmd.Output()
	fmt.Println(printCmd, params, data.Printer, worker.storageDir+data.File)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}

func (worker *Worker) DownloadFile(data JobData) error {
	worker.MakeStorageDir()
	fmt.Println("-- Downloading File...")
	storagePath := worker.storageDir + data.File
	fileUrl := os.Getenv("REMOTE_FILE") + data.File

	out, err := os.Create(storagePath)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (worker *Worker) Watch() {
	worker.MasterBeanstalk.WatchTube(worker.tube)
	fmt.Println("Waiting For Jobs....")
}

func (worker *Worker) Connect() {
	worker.MasterBeanstalk.Connect()
	worker.Watch()
}

func (worker *Worker) MakeStorageDir() {
	if _, err := os.Stat(worker.storageDir); os.IsNotExist(err) {
		fmt.Println("Making Storage Directory")
		err = os.MkdirAll(worker.storageDir, 0777)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Storage Directory Already Made")
	}
}

func makeNewWorker(serverAddress string, tube string) *Worker {
	worker := Worker{tube: tube}
	worker.ServerAddress = serverAddress
	worker.storageDir = "storage/"

	return &worker
}
