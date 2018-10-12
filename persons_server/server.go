package persons_server

import (
	"io"
	"testTask/microservice_comunication"
	"testTask/persons_server/delegate"
	"testTask/persons_server/repository"
)

type server struct {
	personDelegate *delegate.PersonDelegate
}

func NewServer(storage repository.PersonRepository) *server {
	return &server{
		personDelegate: delegate.NewPersonDelegate(storage),
	}
}

func (s *server) AddRecord(stream contract.Persons_AddRecordServer) error {
	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := s.personDelegate.AddPersonRecord(rec); err != nil {
			stream.Send(&contract.AddPersonResponse{
				Processed:    false,
				ErrorMessage: err.Error(),
			})
		} else {
			stream.Send(&contract.AddPersonResponse{
				Processed: true,
			})
		}
	}
}
