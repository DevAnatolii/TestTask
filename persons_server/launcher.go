package persons_server

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"testTask/microservice_comunication"
	"testTask/persons_server/repository"
	"testTask/persons_server/repository/temporary"
)

const (
	StoragePersistent int = iota
	StorageTemporary
)

func Start(port int, storageType int) (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v \n", err)
		return
	}

	storage, err := createStorage(storageType)
	if err != nil {
		log.Fatalf("Error during launch server: %q \n", err)
		return
	}

	s := grpc.NewServer()
	contract.RegisterPersonsServer(s, NewServer(storage))
	log.Printf("Start persons server on : localhost:%d \n", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
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
