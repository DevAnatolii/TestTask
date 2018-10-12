package temporary

import (
	"fmt"
	"sync"
	"testTask/persons_server/model"
)

type PersonStorage struct {
	storage map[int]*model.Person
	mux     sync.Mutex
}

func NewPersonStorage() *PersonStorage {
	return &PersonStorage{
		storage: make(map[int]*model.Person),
	}
}

func (p *PersonStorage) GetPerson(id int) (person *model.Person, ok bool) {
	p.mux.Lock()
	person, ok = p.storage[id]
	p.mux.Unlock()
	return
}

func (p *PersonStorage) SavePerson(person *model.Person) {
	p.mux.Lock()
	if len(p.storage) > 1000 { // need to reset cache, because this collection consumes a lot of RAM memory in case of huge uploading files
		p.storage = make(map[int]*model.Person)
	}
	p.storage[person.Id] = person
	fmt.Printf("Put id in map: %d\n", person.Id)
	p.mux.Unlock()
}
