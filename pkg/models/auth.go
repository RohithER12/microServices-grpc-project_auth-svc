package models

type User struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	Email    string `json:"email"`
	Password string `json:"password"`
	MobileNo string `json:"phone_number"`
	Verified bool   `json:"verified" gorm:"default:false"`
	Blocked  bool   `json:"blocked" gorm:"default:false"`
}

type RegisterOTPValidation struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	MobileNo string `json:"phone_number"`
	Key      string `json:"key"`
}
