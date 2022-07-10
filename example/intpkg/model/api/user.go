package api

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	FirstName *string
	LastName  string
	Age       string
	Address   *Address
	Contact   *Contact
}
