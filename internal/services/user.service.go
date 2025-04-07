package services

import (
	"errors"

	"github.com/thanhdev1710/flamee_auth/internal/repo"
)

type UserServices struct {
	userRepo *repo.UserRepo
}

func NewUserServices() *UserServices {
	return &UserServices{
		userRepo: repo.NewUserRepo(),
	}
}

func (us *UserServices) ConfirmEmail(email string) error {

	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	user.IsVerified = true

	err = us.userRepo.Save(user)
	if err != nil {
		return err
	}

	return nil
}
