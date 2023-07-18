package main

import (
	"fmt"
	"log"
	"net"

	"github.com/RohithER12/auth-svc/pkg/config"
	"github.com/RohithER12/auth-svc/pkg/db"
	"github.com/RohithER12/auth-svc/pkg/pb"
	"github.com/RohithER12/auth-svc/pkg/repo"
	"github.com/RohithER12/auth-svc/pkg/repo/user_interface"
	"github.com/RohithER12/auth-svc/pkg/services"
	"github.com/RohithER12/auth-svc/pkg/utils"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	user := InitializeUserImpl(h)
	s := services.Server{
		H:    h,
		Jwt:  jwt,
		User: user,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func InitializeUserImpl(h *db.Handler) user_interface.User {
	wire.Build(user_interface.NewUserImpl)
	return &repo.UserImpl{
		H: *h,
	}
}
