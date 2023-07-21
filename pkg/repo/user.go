package repo

import (
	"github.com/RohithER12/auth-svc/pkg/db"
	"github.com/RohithER12/auth-svc/pkg/models"
	"gorm.io/gorm"
)

type UserImpl struct {
	H db.Handler
}

func (u *UserImpl) Register(user models.User) error {

	result := u.H.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserImpl) CreateAddress(address models.Address) error {

	result := u.H.DB.Create(&address)
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

func (u *UserImpl) FindById(id int64) (models.User, error) {
	var user models.User
	var result *gorm.DB
	if result = u.H.DB.Where(&models.User{Id: id}).First(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (u *UserImpl) FindByPhoneNumber(mob string) (models.User, error) {
	var user models.User
	var result *gorm.DB
	if result = u.H.DB.Where(&models.User{MobileNo: mob}).First(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (u *UserImpl) Update(user models.User) error {
	result := u.H.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserImpl) RegisterOTPValidation(user models.RegisterOTPValidation) error {
	result := u.H.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *UserImpl) FindByMobileNoAndKey(key string) (string, error) {
	var otpValidation models.RegisterOTPValidation
	var result *gorm.DB
	if result = u.H.DB.Where("key = ?", key).First(&otpValidation); result.Error != nil {
		return "", result.Error
	}
	return otpValidation.MobileNo, nil
}
