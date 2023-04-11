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

func (userUseCase *userUseCase) Update(ctx context.Context, u domain.User) (user domain.User, err error) {
	if user, err = userUseCase.userRepository.Update(ctx, u); err != nil {
		return user, err
	}

	return user, nil
}

func (userUseCase *userUseCase) DeleteById(ctx context.Context, id string) (err error) {
	if err = userUseCase.userRepository.DeleteById(ctx, id); err != nil {
		return err
	}

	return
}

func (userUseCase *userUseCase) FindByEmail(ctx context.Context, u *domain.User) (user domain.User, err error) {
	if user, err = userUseCase.userRepository.FindByEmail(ctx, u); err != nil {
		return user, err
	}

	return user, nil
}

func (userUseCase *userUseCase) FindByUsername(ctx context.Context, u *domain.User) (user domain.User, err error) {
	if user, err = userUseCase.userRepository.FindByUsername(ctx, u); err != nil {
		return user, err
	}

	return user, nil
}
