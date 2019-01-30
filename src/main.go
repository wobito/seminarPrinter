package main

func main() {
	workerMain()
}

func workerMain() {
	worker := makeNewWorker("localhost:11300", "printer")
	worker.Connect()
	defer worker.Close()

	for {
		worker.ProcessJob()
	}
}
