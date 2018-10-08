package temporary

import (
	"fmt"
	"sync"
	"testTask/persons_server/model"
)

type PersonStorage struct {
	storage map[int]model.Person
	mux     sync.Mutex
}

func NewPersonStorage() *PersonStorage {
	return &PersonStorage{
		storage: make(map[int]model.Person),
	}
}

func (p *PersonStorage) GetPerson(id int) (person model.Person, ok bool) {
	p.mux.Lock()
	person, ok = p.storage[id]
	p.mux.Unlock()
	return
}

func (p *PersonStorage) SavePerson(person model.Person) {
	p.mux.Lock()
	p.storage[person.Id] = person
	p.printStorageState()
	p.mux.Unlock()
}

func (p *PersonStorage) printStorageState() {
	fmt.Println("--------------------------------------------------------")
	fmt.Println("Storage state after")
	for k, v := range p.storage {
		fmt.Printf("id = %d, value = %v \n", k, v)
	}
}
