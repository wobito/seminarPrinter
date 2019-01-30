package main

func main() {
	workerMain()
}

func workerMain() {
	worker := makeNewWorker("45.55.53.144:11300", "printer")
	worker.Connect()
	defer worker.Close()

	for {
		worker.ProcessJob()
	}
}
