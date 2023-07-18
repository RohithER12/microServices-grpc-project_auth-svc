package repo

import (
	"fmt"

	"github.com/RohithER12/auth-svc/pkg/db"
	"github.com/RohithER12/auth-svc/pkg/models"
	"gorm.io/gorm"
)

type UserImpl struct {
	H db.Handler
}

func (u *UserImpl) Register(user models.User) error {
	fmt.Println(
		"\n here",
		"\nuser", user,
		"\nh\n", u.H,
	)
	result := u.H.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// func (u *UserImpl) Login(email string) (models.User, error) {
// 	var user models.User
// 	var result *gorm.DB
// 	if result = u.H.DB.Where(&models.User{Email: email}).First(&user); result.Error != nil {
// 		return models.User{}, result.Error
// 	}
// 	return user, nil
// }

func (u *UserImpl) FindByEmail(email string) (models.User, error) {
	var user models.User
	var result *gorm.DB
	if result = u.H.DB.Where(&models.User{Email: email}).First(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}
