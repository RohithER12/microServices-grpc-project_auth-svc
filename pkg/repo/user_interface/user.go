package user_interface

import (
	"github.com/RohithER12/auth-svc/pkg/models"
	"github.com/RohithER12/auth-svc/pkg/repo"
)

type User interface {
	Register(models.User) error
	FindById(id int64) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByPhoneNumber(mob string) (models.User, error)
	Update(user models.User) error
	RegisterOTPValidation(user models.RegisterOTPValidation) error
	FindByMobileNoAndKey(key string) (string, error)
	CreateAddress(address models.Address) error
}

func NewUserImpl() User {
	return &repo.UserImpl{}
}
