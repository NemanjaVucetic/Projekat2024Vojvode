package main

import "userService/startup"

func main() {
	config1 := startup.NewConfig()
	server := startup.NewServer(config1)
	server.Start()
}
