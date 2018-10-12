package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"testTask/persons_server"
	"testTask/upload_server"
)

const (
	storagePersistent = "persistent"
	storageTemporary  = "temporary"
)

var (
	personServerPort = flag.Int("personServerPort", 8000, "Port, which will be used for persons server")
	uploadServerPort = flag.Int("uploadServerPort", 8080, "Port, which will be used for upload server")
	storageType      = flag.String("storageType", storageTemporary, "Type of storage")
	errorFileLogging = flag.String("errorFile", "errorLogs.txt", "File for logging error records")
)

func main() {
	flag.Parse()

	if *personServerPort == *uploadServerPort {
		log.Println("Could not start servers with the same ports")
		return
	}

	t, ok := obtainStorageType()
	if !ok {
		log.Println("Could not start persons server: unknown storage type")
		return
	}

	var err error
	go func() {
		err = persons_server.Start(*personServerPort, t)
	}()
	if err != nil {
		fmt.Println("Could not start server")
		return
	}

	var baseUrlPattern = "localhost:%d"
	upload_server.Start(fmt.Sprintf(baseUrlPattern, *uploadServerPort), fmt.Sprintf(baseUrlPattern, *personServerPort),
		*errorFileLogging)
}

func obtainStorageType() (int, bool) {
	flag.Parse()
	switch {
	case strings.ToLower(*storageType) == storagePersistent:
		return persons_server.StoragePersistent, true
	case strings.ToLower(*storageType) == storageTemporary:
		return persons_server.StorageTemporary, true
	default:
		return -1, false
	}
}
