package main

import (
	"userService/startup"
)

func main() {
	config1 := startup.NewConfig()
	//fmt.Printf("Starting server with config: %+v\n", config1) // Print the configuration
	server := startup.NewServer(config1)
	server.Start()
}
