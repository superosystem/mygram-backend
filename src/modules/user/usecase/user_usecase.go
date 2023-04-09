package usecase

import (
	"context"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(userRepository domain.UserRepository) *userUseCase {
	return &userUseCase{userRepository}
}

func (userUseCase *userUseCase) Register(ctx context.Context, registerUser *domain.RegisterUser) (user domain.User, err error) {
	user.Email = registerUser.Email
	user.Username = registerUser.Username
	user.Password = registerUser.Password
	user.Age = registerUser.Age

	if err = userUseCase.userRepository.Register(ctx, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (userUseCase *userUseCase) Login(ctx context.Context, loginUser *domain.LoginUser) (user domain.User, err error) {
	user.Email = loginUser.Email
	user.Password = loginUser.Password

	if err = userUseCase.userRepository.Login(ctx, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (userUseCase *userUseCase) Update(ctx context.Context, user domain.User) (u domain.User, err error) {
	if u, err = userUseCase.userRepository.Update(ctx, user); err != nil {
		return u, err
	}

	return u, nil
}

func (userUseCase *userUseCase) Delete(ctx context.Context, id string) (err error) {
	if err = userUseCase.userRepository.Delete(ctx, id); err != nil {
		return err
	}

	return
}
