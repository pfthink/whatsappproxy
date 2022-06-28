package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pfthink/whatsappproxy/src/structs"
	"github.com/pfthink/whatsappproxy/src/utils"
)

func ValidateUserInfo(request structs.UserInfoRequest) {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, is.E164, validation.Length(10, 15)),
	)

	if err != nil {
		panic(utils.ValidationError{
			Message: err.Error(),
		})
	}
}
func ValidateUserAvatar(request structs.UserAvatarRequest) {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.Phone, validation.Required, is.E164, validation.Length(10, 15)),
	)

	if err != nil {
		panic(utils.ValidationError{
			Message: err.Error(),
		})
	}
}
