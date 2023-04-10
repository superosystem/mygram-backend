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
	Username     string         `gorm:"type:VARCHAR(50);uniqueIndex;not null" valid:"required" form:"username" json:"username" example:"johndoe"`
	Email        string         `gorm:"type:VARCHAR(50);uniqueIndex;not null" valid:"email,required" form:"email" json:"email" example:"johndoe@example.com"`
	Password     string         `gorm:"not null" valid:"required,minstringlength(6)" form:"password" json:"password,omitempty" example:"secret"`
	Age          uint           `gorm:"not null" valid:"required,range(8|63)" form:"age" json:"age,omitempty" example:"8"`
	CreatedAt    *time.Time     `gorm:"not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `gorm:"not null;autocreateTime" json:"updated_at,omitempty"`
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
	Update(context.Context, UpdateUser, string) (User, error)
	Delete(context.Context, string) error
}

type UserRepository interface {
	Register(context.Context, *User) error
	Login(context.Context, *User) error
	Update(context.Context, User, string) (User, error)
	DeleteById(context.Context, string) error
	FindByEmail(context.Context, User) (User, error)
	FindByUsername(context.Context, User) (User, error)
}

// Represents for register user
type RegisterUser struct {
	Email    string `json:"email" form:"email" example:"johndoe@example.com"`
	Username string `json:"username" form:"username" example:"johndoe"`
	Password string `json:"password,omitempty" form:"password" example:"secret"`
	Age      uint   `json:"age" form:"age" example:"8"`
}

// Reprensents for registered user
type RegisteredUser struct {
	Status  string  `json:"status" example:"success"`
	Message string  `json:"message" example:"message you if the process has been successful"`
	Data    GetUser `json:"data"`
}

// Represents for login user
type LoginUser struct {
	Email    string `json:"email" form:"email" example:"johndoe@example.com"`
	Password string `json:"password,omitempty" form:"password" example:"secret"`
}

// Represents for jwt user
type Token struct {
	Token string `json:"token" example:"the token generated here"`
}

// Represents for loggedin user
type LoggedInUser struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"message you if the process has been successful"`
	Data    Token  `json:"data"`
}

// Represents for update user
type UpdateUser struct {
	Email    string `json:"email" example:"newjohndoe@example.com"`
	Username string `json:"username" example:"newjohndoe"`
	Age      uint   `json:"age" example:"8"`
}

// Represents for updated user
type UpdatedDataUser struct {
	ID        string     `json:"id" example:"here is the generated user id"`
	Email     string     `json:"email" example:"newjohndoe@example.com"`
	Username  string     `json:"username" example:"newjohndoe"`
	Age       uint       `json:"age" example:"8"`
	UpdatedAt *time.Time `json:"updated_at" example:"update time should be here"`
}

// Represents for response updated user
type UpdatedUser struct {
	Status  string          `json:"status" example:"success"`
	Message string          `json:"message" example:"message you if the process has been successful"`
	Data    UpdatedDataUser `json:"data"`
}

// Represents for response deleted user
type DeletedUser struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your account has been successfully deleted"`
}

// Represents for Get User
type GetUser struct {
	ID       string `json:"id"`
	Username string `json:"username" example:"newjohndoe"`
	Email    string `json:"email" example:"newjohndoe@example.com"`
}
