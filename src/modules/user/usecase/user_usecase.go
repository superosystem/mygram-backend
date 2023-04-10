package usecase

import (
	"context"
	"errors"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"gorm.io/gorm"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(userRepository domain.UserRepository) *userUseCase {
	return &userUseCase{userRepository}
}

func (userUseCase *userUseCase) Register(ctx context.Context, user *domain.User) (err error) {
	if err = userUseCase.userRepository.Register(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if _, err = userUseCase.userRepository.FindByEmail(ctx, user); err == nil {
				return errors.New("duplicate key on idx_user_email")
			}

			if _, err = userUseCase.userRepository.FindByUsername(ctx, user); err == nil {
				return errors.New("duplicate key on idx_user_username")
			}
		}
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if _, err = userUseCase.userRepository.FindByEmail(ctx, &user); err == nil {
				return user, errors.New("duplicate key on idx_user_email")
			}

			if _, err = userUseCase.userRepository.FindByUsername(ctx, &user); err == nil {
				return user, errors.New("duplicate key on idx_user_username")
			}
		}
		return user, err
	}

	return user, nil
}

func (userUseCase *userUseCase) Delete(ctx context.Context, id string) (err error) {
	if err = userUseCase.userRepository.DeleteById(ctx, id); err != nil {
		return err
	}

	return
}
