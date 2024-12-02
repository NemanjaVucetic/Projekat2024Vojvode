package main

import "notificationService/startup"

func main() {
	config1 := startup.NewConfig()
	server := startup.NewServer(config1)
	server.Start()
}
