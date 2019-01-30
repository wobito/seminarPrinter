package main

type JobData struct {
	File    string `json:"file"`
	Printer string `json:"printer"`
}

type Job struct {
	Job  string  `json:"job"`
	Data JobData `json:"data"`
}
