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

func (userUseCase *userUseCase) Register(ctx context.Context, user *domain.User) (err error) {
	if err = userUseCase.userRepository.Register(ctx, user); err != nil {
		return err
	}

	return
}

func (userUseCase *userUseCase) Login(ctx context.Context, user *domain.User) (err error) {
	if err = userUseCase.userRepository.Login(ctx, user); err != nil {
		return err
	}

	return
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
