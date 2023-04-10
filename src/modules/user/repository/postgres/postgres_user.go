package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"gorm.io/gorm"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (userRepository *userRepository) Register(ctx context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	ID, _ := gonanoid.New(16)

	user.ID = fmt.Sprintf("user-%s", ID)

	if err = userRepository.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return
}

func (userRepository *userRepository) Login(ctx context.Context, user *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	password := user.Password

	if err = userRepository.db.WithContext(ctx).Where("email = ?", user.Email).Take(&user).Error; err != nil {
		return errors.New("the email you entered are not registered")
	}

	if isValid := helpers.Compare([]byte(user.Password), []byte(password)); !isValid {
		return errors.New("the password you entered are wrong")
	}

	return
}

func (userRepository *userRepository) Update(ctx context.Context, u domain.User) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user = domain.User{}

	if err = userRepository.db.WithContext(ctx).First(&user).Error; err != nil {
		return user, err
	}

	if user.Email == u.Email {
		u.Email = ""
	}

	if user.Username == u.Username {
		u.Username = ""
	}

	if err = userRepository.db.WithContext(ctx).Model(&user).Updates(&u).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (userRepository *userRepository) DeleteById(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).First(&domain.User{}, &id).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Where("user_id = ?", id).Delete(&domain.SocialMedia{}).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Delete(&domain.User{}, &id).Error; err != nil {
		return err
	}

	return
}

func (userRepository *userRepository) FindByEmail(ctx context.Context, u *domain.User) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).First(&user, "email = ?", u.Email).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (userRepository *userRepository) FindByUsername(ctx context.Context, u *domain.User) (user domain.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = userRepository.db.WithContext(ctx).First(&user, "username = ?", u.Username).Error; err != nil {
		return user, err
	}

	return user, nil
}
