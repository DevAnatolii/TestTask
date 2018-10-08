package persons_server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"testTask/persons_server/handler"
	"testTask/persons_server/repository"
	"testTask/persons_server/repository/temporary"
)

const (
	StoragePersistent = 1
	StorageTemporary  = 2
)

func Start(address string, storageType int) (err error) {
	serveMux := http.NewServeMux()

	storage, err := createStorage(storageType)
	if err != nil {
		fmt.Printf("Error during launch server: %q", err)
		return
	}

	studentsHandler := handler.NewPersonHandler(storage)
	serveMux.Handle(handler.HandlePath, studentsHandler)

	log.Fatal(http.ListenAndServe(address, serveMux))
	return
}

func createStorage(storageType int) (storage repository.PersonRepository, err error) {
	switch storageType {
	case StoragePersistent:
		err = errors.New("this type is not supported")
	case StorageTemporary:
		storage = temporary.NewPersonStorage()
	default:
		err = errors.New("unknown type")
	}
	return
}
