package repository

import "testTask/persons_server/model"

type PersonRepository interface {
	GetPerson(id int) (model.Person, bool)
	SavePerson(p model.Person)
}
