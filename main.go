package main

import (
	"fmt"
	"testTask/persons_server"
	"testTask/upload_server"
)

const (
	PersonsServerAddress = "localhost:8000"
	UploadServerAddress  = "localhost:8080"
)

func main() {
	err := persons_server.Start(PersonsServerAddress, persons_server.StorageTemporary)
	if err != nil {
		fmt.Println("Could not start server")
		return
	}
	upload_server.Start(UploadServerAddress, PersonsServerAddress)
}
