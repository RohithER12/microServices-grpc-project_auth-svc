package services

import (
	"context"
	"net/http"

	"github.com/RohithER12/auth-svc/pkg/db"
	"github.com/RohithER12/auth-svc/pkg/models"
	"github.com/RohithER12/auth-svc/pkg/pb"
	"github.com/RohithER12/auth-svc/pkg/repo/user_interface"
	"github.com/RohithER12/auth-svc/pkg/utils"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	H    *db.Handler
	Jwt  utils.JwtWrapper
	User user_interface.User
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	// if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
	// 	return &pb.RegisterResponse{
	// 		Status: http.StatusConflict,
	// 		Error:  "E-Mail already exists",
	// 	}, nil
	// }
	_, err := s.User.FindByEmail(req.Email)
	if err == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "E-Mail already exists",
		}, nil
	}

	_, err = s.User.FindByPhoneNumber(req.PhoneNumber)
	if err == nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "Mobile Number already exists",
		}, nil
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)
	user.MobileNo = req.PhoneNumber

	// s.H.DB.Create(&user)
	err = s.User.Register(user)
	if err != nil {

		return nil, err
	}
	otpValidationKey, err := utils.SendOtp(req.PhoneNumber)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "Otp sending failed",
		}, nil
	}

	creatingOtp := models.RegisterOTPValidation{
		MobileNo: req.PhoneNumber,
		Key:      otpValidationKey,
	}

	err = s.User.RegisterOTPValidation(creatingOtp)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "otp key and mob no saving failed",
		}, nil
	}

	return &pb.RegisterResponse{
		Status:           http.StatusCreated,
		OtpValidationKey: otpValidationKey,
	}, nil
}

func (s *Server) RegisterOTPValidation(ctx context.Context, req *pb.RegisterOTPValidationRequest) (*pb.RegisterOTPValidationResponse, error) {

	fetchingMobNo, err := s.User.FindByMobileNoAndKey(req.Key)
	if err != nil {
		return &pb.RegisterOTPValidationResponse{
			Status: http.StatusConflict,
			Error:  "Invalid Key",
		}, nil
	}
	err = utils.CheckOtp(fetchingMobNo, req.Otp)
	if err != nil {
		return &pb.RegisterOTPValidationResponse{
			Status: http.StatusConflict,
			Error:  "Invalid otp",
		}, nil
	}
	user, err := s.User.FindByPhoneNumber(fetchingMobNo)
	if err != nil {
		return &pb.RegisterOTPValidationResponse{
			Status: http.StatusConflict,
			Error:  "Fetch user data using phone number error in registerOTPvalidation",
		}, nil
	}
	user.Verified = true
	err = s.User.Update(user)
	if err != nil {
		return &pb.RegisterOTPValidationResponse{
			Status: http.StatusConflict,
			Error:  "can't update register verfication",
		}, nil
	}
	return &pb.RegisterOTPValidationResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	email := req.Email
	// if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
	// 	return &pb.LoginResponse{
	// 		Status: http.StatusNotFound,
	// 		Error:  "User not found",
	// 	}, nil
	// }
	result, err := s.User.FindByEmail(email)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, result.Password)

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	if !result.Verified {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Account not verified",
		}, nil
	}

	if result.Blocked {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Account Blocked",
		}, nil
	}

	token, _ := s.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	user, err := s.User.FindByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return &pb.ResetPasswordResponse{
			Status: http.StatusConflict,
			Error:  "Mobile Number not valid",
		}, nil
	}

	otpValidationKey, err := utils.SendOtp(user.MobileNo)
	if err != nil {
		return &pb.ResetPasswordResponse{
			Status: http.StatusConflict,
			Error:  "Otp sending failed",
		}, nil
	}

	creatingOtp := models.RegisterOTPValidation{
		MobileNo: req.PhoneNumber,
		Key:      otpValidationKey,
	}

	err = s.User.RegisterOTPValidation(creatingOtp)
	if err != nil {
		return &pb.ResetPasswordResponse{
			Status: http.StatusConflict,
			Error:  "otp key and mob no saving failed",
		}, nil
	}

	return &pb.ResetPasswordResponse{
		Status:           http.StatusOK,
		OtpValidationKey: otpValidationKey,
	}, nil
}

func (s *Server) ResetPasswordValidation(ctx context.Context, req *pb.ResetPasswordValidationRequest) (
	*pb.ResetPasswordValidationResponse, error) {
	fetchingMobNo, err := s.User.FindByMobileNoAndKey(req.Key)
	if err != nil {
		return &pb.ResetPasswordValidationResponse{
			Status: http.StatusConflict,
			Error:  "Invalid Key",
		}, nil
	}

	err = utils.CheckOtp(fetchingMobNo, req.Otp)
	if err != nil {
		return &pb.ResetPasswordValidationResponse{
			Status: http.StatusConflict,
			Error:  "Invalid otp",
		}, nil
	}

	user, err := s.User.FindByPhoneNumber(fetchingMobNo)
	if err != nil {
		return &pb.ResetPasswordValidationResponse{
			Status: http.StatusConflict,
			Error:  "Mobile Number not valid",
		}, nil
	}

	user.Password = utils.HashPassword(req.Paswword)

	if s.User.Update(user); err != nil {
		return &pb.ResetPasswordValidationResponse{
			Status: http.StatusConflict,
			Error:  "Reset Password Failed",
		}, nil
	}

	return &pb.ResetPasswordValidationResponse{
		Status: http.StatusOK,
	}, nil

}

func (s *Server) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {

	user, err := s.User.FindById(req.UserId)
	if err != nil {
		return &pb.AddAddressResponse{
			Status: http.StatusConflict,
			Error:  "Address creating  Failed",
		}, nil
	}
	address := models.Address{
		DoorNo:     req.DoorNo,
		City:       req.City,
		PostalCode: req.PostalCodev,
		UserId:     req.UserId,
		MobileNo:   user.MobileNo,
	}
	if s.User.CreateAddress(address); err != nil {
		return &pb.AddAddressResponse{
			Status: http.StatusConflict,
			Error:  "Address creating  Failed",
		}, nil
	}

	return &pb.AddAddressResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) UserProfile(ctx context.Context, req *pb.UserProfileRequest) (*pb.UserProfileResponse, error) {
	user, err := s.User.FindById(req.UserId)
	if err != nil {
		return &pb.UserProfileResponse{
			Status: http.StatusConflict,
			Error:  "Fetching User Data Failed",
		}, nil
	}
	return &pb.UserProfileResponse{
		Email:       user.Email,
		PhoneNumber: user.MobileNo,
		Status:      http.StatusCreated,
	}, nil

}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	var user models.User

	if result := s.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}

// // Add an empty implementation of mustEmbedUnimplementedAuthServiceServer to satisfy the interface requirements.
// func (s *Server) mustEmbedUnimplementedAuthServiceServer() {
// 	// Empty implementation
// }
