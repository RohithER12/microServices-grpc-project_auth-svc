package user_interface

import (
	"github.com/RohithER12/auth-svc/pkg/models"
	"github.com/RohithER12/auth-svc/pkg/repo"
)

type User interface {
	Register(models.User) error
	Login(models.User) error
}

func NewUserImpl() User {
	return &repo.UserImpl{}
}
