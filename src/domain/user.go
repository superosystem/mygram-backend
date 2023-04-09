package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gusrylmubarok/mygram-backend/src/helpers"
	"gorm.io/gorm"
)

// User represents entity for a user
type User struct {
	ID           string         `gorm:"primaryKey;type:VARCHAR(50)" json:"id"`
	Username     string         `gorm:"column:username;type:VARCHAR(50);uniqueIndex;not null" valid:"required"`
	Email        string         `gorm:"type:VARCHAR(50);uniqueIndex;not null" valid:"email,required"`
	Password     string         `gorm:"not null" valid:"required,minstringlength(6)" json:"-"`
	Age          uint           `gorm:"not null" valid:"required,range(8|63)"`
	CreatedAt    *time.Time     `gorm:"not null;autoCreateTime"`
	UpdatedAt    *time.Time     `gorm:"not null;autocreateTime"`
	Photos       *[]Photo       `json:"-"`
	SocialMedias *[]SocialMedia `json:"-"`
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
	Register(context.Context, *RegisterUser) (User, error)
	Login(context.Context, *LoginUser) (User, error)
	Update(context.Context, User) (User, error)
	DeleteById(context.Context, string) error
}

type UserRepository interface {
	Register(context.Context, *User) error
	Login(context.Context, *User) error
	Update(context.Context, User) (User, error)
	DeleteById(context.Context, string) error
}

// Represents for register user
type RegisterUser struct {
	Email    string `json:"email" form:"email" example:"johndoe@example.com"`
	Username string `json:"username" form:"username" example:"johndoe"`
	Password string `json:"password,omitempty" form:"password" example:"secret"`
	Age      uint   `json:"age" form:"age" example:"8"`
}

// Represents for registered use
type RegisteredUser struct {
	ID       string `json:"id" example:"the user id generated here"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Username string `json:"username" example:"johndoe"`
	Age      uint   `json:"age" example:"8"`
}

// Reprensents for response success register user
type ResponseRegisteredUser struct {
	Status string         `json:"status" example:"success"`
	Data   RegisteredUser `json:"data"`
}

// Represents for login user
type LoginUser struct {
	Email    string `json:"email" form:"email" example:"johndoe@example.com"`
	Password string `json:"password,omitempty" form:"password" example:"secret"`
}

// Represents for loggedin user
type LoggedInUser struct {
	Token string `json:"token" example:"the token generated here"`
}

// Represents for response loggedin user
type ResponseLoggedInUser struct {
	Status string       `json:"status" example:"success"`
	Data   LoggedInUser `json:"data"`
}

// Represents for update user
type UpdateUser struct {
	Email    string `json:"email" example:"newjohndoe@example.com"`
	Username string `json:"username" example:"newjohndoe"`
	Age      uint   `json:"age" example:"8"`
}

// Represents for updated user
type UpdatedUser struct {
	ID        string     `json:"id" example:"here is the generated user id"`
	Email     string     `json:"email" example:"newjohndoe@example.com"`
	Username  string     `json:"username" example:"newjohndoe"`
	Age       uint       `json:"age" example:"8"`
	CreatedAt *time.Time `json:"created_at" example:"create time should be here"`
	UpdatedAt *time.Time `json:"updated_at" example:"update time should be here"`
}

// Represents for response updated user
type ResponseUpdatedUser struct {
	Status string      `json:"status" example:"success"`
	Data   UpdatedUser `json:"data"`
}

// Represents for response deleted user
type ResponseMessageDeletedUser struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your account has been successfully deleted"`
}
