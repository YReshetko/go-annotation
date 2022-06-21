package intpkg

type UserDTO struct {
}

type User struct {
}

// UserDTOMapper my test DTO mapper
// Mapper(name="UserDTOMapperImplementation")
// @Mapper
type UserDTOMapper interface {

	// ToUserDTO my test DTO mapper function
	// @Mapping(target="fn", source="firstName")
	// @Mapping(target="ln", source="lastName")
	ToUserDTO(user *User) *UserDTO

	ToUser(UserDTO) User
}
