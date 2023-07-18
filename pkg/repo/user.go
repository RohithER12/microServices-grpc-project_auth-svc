package repo

import (
	"fmt"

	"github.com/RohithER12/auth-svc/pkg/db"
	"github.com/RohithER12/auth-svc/pkg/models"
	"github.com/RohithER12/auth-svc/pkg/pb"
)

type UserImpl struct {
	H db.Handler
	pb.UnimplementedAuthServiceServer
}

func (u *UserImpl) Register(user models.User) error {
	fmt.Println(
		"\n here",
		"\nuser", user,
		"\nh\n", u.H,
	)
	u.H.DB.Create(&user)
	return nil
}

func (u *UserImpl) Login(user models.User) error {
	// Implement the Register method
	return nil
}
