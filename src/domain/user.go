package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID              string         `gorm:"primaryKey;type:VARCHAR(50)" json:"id"`
	Username        string         `gorm:"type:VARCHAR(50);uniqueIndex;not null" valid:"required" form:"username" json:"username" example:"johndoe"`
	Email           string         `gorm:"type:VARCHAR(50);uniqueIndex;not null" valid:"email,required" form:"email" json:"email" example:"johndoe@example.com"`
	Password        string         `gorm:"not null" valid:"required,minstringlength(6)" form:"password" json:"password,omitempty" example:"secret"`
	Age             uint           `gorm:"not null" valid:"required,range(8|63)" form:"age" json:"age,omitempty" example:"8"`
	ProfileImageUrl string         `json:"profileImageUrl,omitempty" example:"https://www.example.com/image.jpg"`
	CreatedAt       *time.Time     `gorm:"not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt       *time.Time     `gorm:"not null;autocreateTime" json:"updated_at,omitempty"`
	Photos          *[]Photo       `json:"-"`
	SocialMedias    *[]SocialMedia `json:"-"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return err
	}

	user.Password = helpers.Hash(user.Password)

	return
}

func (user *User) BeforeUpdate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return err
	}

	return
}

type UserUseCase interface {
	Register(context.Context, *User) error
	Login(context.Context, *User) error
	Update(context.Context, User) (User, error)
	Delete(context.Context, string) error
}

type UserRepository interface {
	Register(context.Context, *User) error
	Login(context.Context, *User) error
	Update(context.Context, User) (User, error)
	Delete(context.Context, string) error
}
