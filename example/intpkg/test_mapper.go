package intpkg

import (
	"fmt"
	"github.com/YReshetko/go-annotation/example/intpkg/model/api"
	"github.com/YReshetko/go-annotation/example/intpkg/model/db"
	"github.com/google/uuid"
)

// UserDTOMapper my test DTO mapper
// Mapper(name="UserDTOMapperImplementation")
// @Mapper
type UserMapper interface {

	// ToUserDTO my test DTO mapper function
	// @Mapping(target="FirstName", source="Name")
	// @Mapping(target="Contact.Email", source="Email")
	DBToAPI(user db.User) api.User

	//APIToDB(user api.User) db.User
}

func fdfsf(user db.User) api.User {
	res_0 := api.User{}
	res_0.ID = uuid.UUID{}
	res_0.FirstName = &user.Name
	if user.Age != nil {
		res_0.Age = fmt.Sprintf("%f", *user.Age)
	}
	res_0.Address = &api.Address{}
	res_0.Contact = &api.Contact{}
	res_0.Contact.Email = user.Email

	return res_0
}
