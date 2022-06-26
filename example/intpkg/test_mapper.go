package intpkg

import "github.com/YReshetko/go-annotation/example/intpkg/model"

type UserDTO struct {
}

// UserDTOMapper my test DTO mapper
// Mapper(name="UserDTOMapperImplementation")
// @Mapper
type UserDTOMapper interface {

	// ToUserDTO my test DTO mapper function
	// @Mapping(target="fn", source="firstName")
	// @Mapping(target="ln", source="lastName")
	ToUserDTO(user model.User, i *int, udto UserDTOMapper, fn func(rt int) float32, f **model.ExternalFunction, in interface{}, st struct{}, e error) *UserDTO

	ToUser(UserDTO) *model.User
}

// Comment of the group
type (

	// BlockUserDTOMapper my test DTO mapper
	// Mapper(name="UserDTOMapperImplementation")
	// @Mapper
	BlockUserDTOMapper interface {

		// ToUserDTO my test DTO mapper function
		// @Mapping(target="fn", source="firstName")
		// @Mapping(target="ln", source="lastName")
		ToUserDTO(user model.User, i *int, udto UserDTOMapper, fn func(rt int) float32, f **model.ExternalFunction, in interface{}, st struct{}, e error) *UserDTO

		ToUser(UserDTO) *model.User
	}
)
