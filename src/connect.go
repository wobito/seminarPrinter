package main

import (
	"fmt"
	"os"

	beanstalk "github.com/beanstalkd/go-beanstalk"
)

type MasterBeanstalk struct {
	ServerAddress    string
	serverConnection *beanstalk.Conn
	tubeSet          *beanstalk.TubeSet
}

func (master *MasterBeanstalk) Connect() {
	beanstalkConnection, err := beanstalk.Dial("tcp", master.ServerAddress)
	if err != nil {
		// do retries or whatever you need
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Connected To Server")
	master.serverConnection = beanstalkConnection
}

func (master *MasterBeanstalk) WatchTube(tube string) {
	master.tubeSet = beanstalk.NewTubeSet(master.serverConnection, tube)
	fmt.Println("Connected To Tube")
}

func (master *MasterBeanstalk) Close() {
	if master.serverConnection != nil {
		master.serverConnection.Close()
	}
}
