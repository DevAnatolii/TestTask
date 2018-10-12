package delegate

import (
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	"strings"
	"testTask/microservice_comunication"
	"testTask/persons_server/model"
	"testTask/persons_server/repository"
	"unicode"
)

const ukPhonePrefix = "(+44)"

type PersonDelegate struct {
	storage repository.PersonRepository
}

func NewPersonDelegate(repository repository.PersonRepository) *PersonDelegate {
	return &PersonDelegate{repository}
}

func (pd *PersonDelegate) AddPersonRecord(request *contract.AddPersonRequest) error {
	if len(strings.TrimSpace(request.Name)) == 0 {
		return errors.New("empty name")
	}

	if err := checkmail.ValidateFormat(request.Email); err != nil {
		return errors.New("invalid email")
	}

	formattedPhone, ok := formatPhone(request.Phone)
	if !ok {
		return errors.New("invalid phone number")
	}

	person := &model.Person{
		Id:    int(request.Id),
		Name:  request.Name,
		Email: request.Email,
		Phone: formattedPhone,
	}
	pd.storage.SavePerson(person)
	return nil
}

func formatPhone(phone string) (string, bool) {
	res := make([]rune, 0)
	for i := 0; i < len(phone); i++ {
		symbol := rune(phone[i])
		switch {
		case symbol == ' ' || symbol == ')' || symbol == '(':
			break
		case unicode.IsDigit(symbol):
			res = append(res, symbol)
			break
		default:
			return "", false
		}
	}
	return fmt.Sprintf("%s%s", ukPhonePrefix, string(res)), true
}
