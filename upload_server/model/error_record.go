package model

import (
	"fmt"
)

type ErrorRecord struct {
	*Record
	Error error
}

func (er *ErrorRecord) String() string {
	return fmt.Sprintf("\"id = %d, name = %s, email = %s, phone = %s\" - %s",
		er.Id, er.Name, er.Email, er.Phone, er.Error)
}
