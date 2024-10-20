package main

import (
	server "akshidas/e-com/pkg/server"
)

func main() {
	server := &server.APIServer{Status: "Server is up and running", Port: ":5234"}
	server.Run()
}
