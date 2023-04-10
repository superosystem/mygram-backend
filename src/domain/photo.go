package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	ID        string     `gorm:"primaryKey;type:VARCHAR(50)"`
	Title     string     `gorm:"type:VARCHAR(50);not null" valid:"required" `
	Caption   string     `gorm:"type:VARCHAR(50)"`
	PhotoUrl  string     `gorm:"not null" valid:"required"`
	UserID    string     `gorm:"type:VARCHAR(50);not null"`
	User      *User      `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
	CreatedAt *time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"not null;autoCreateTime"`
	Comment   *Comment
}

func (photo *Photo) BeforeCreate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(photo); err != nil {
		return err
	}

	return
}

func (photo *Photo) BeforeUpdate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(photo); err != nil {
		return err
	}
	return
}

type PhotoRepository interface {
	Save(context.Context, *Photo) error
	Update(context.Context, Photo, string) (Photo, error)
	DeleteById(context.Context, string) error
	FindAll(context.Context, *[]Photo) error
	FindById(context.Context, *Photo, string) error
}

type PhotoUseCase interface {
	Save(context.Context, *AddPhoto, string) (Photo, error)
	Update(context.Context, *UpdatePhoto, string) (Photo, error)
	DeleteById(context.Context, string) error
	FindAll(context.Context, *[]Photo) error
	FindById(context.Context, *Photo, string) error
}

// Represents for request add photo
type AddPhoto struct {
	Title    string `json:"title" form:"title" example:"A Photo Title"`
	Caption  string `json:"caption" form:"caption" example:"A caption"`
	PhotoUrl string `json:"photo_url" form:"photo_url" example:"https://www.example.com/image.jpg"`
}

// Represents for added photo
type AddedDataPhoto struct {
	ID        string     `json:"id"`
	Title     string     `json:"title" form:"title" example:"A Photo Title"`
	Caption   string     `json:"caption" form:"caption" example:"A caption"`
	PhotoUrl  string     `json:"photo_url"  form:"photo_url" example:"https://www.example.com/image.jpg"`
	User      *GetUser   `json:"user"`
	CreatedAt *time.Time `json:"created_at" example:"create time should be here"`
}

// Represents for response added photo
type AddedPhoto struct {
	Status  string         `json:"status" example:"success"`
	Message string         `json:"message" example:"message you if the process has been successful"`
	Data    AddedDataPhoto `json:"data"`
}

// Represents for request update photo
type UpdatePhoto struct {
	Title    string `json:"title" example:"A new title"`
	Caption  string `json:"caption" example:"A new caption"`
	PhotoUrl string `json:"photo_url" example:"https://www.example.com/new-image.jpg"`
	UserID   string
}

// Represents for updated data photo
type UpdatedDataPhoto struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	User      *GetUser   `json:"user"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Represents for updated photo
type UpdatedPhoto struct {
	Status  string           `json:"status" example:"success"`
	Message string           `json:"message" example:"message you if the process has been successful"`
	Data    UpdatedDataPhoto `json:"data"`
}

// Represents for response deleted user
type DeletedPhoto struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your account has been successfully deleted"`
}

// RepresentGetPhoto
type GetPhoto struct {
	ID        string     `json:"id"`
	Title     string     `json:"title,"`
	Caption   string     `json:"caption"`
	PhotoUrl  string     `json:"photo_url"`
	User      *GetUser   `json:"user"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type GetAllPhotos struct {
	Status string     `json:"status" example:"success"`
	Data   []GetPhoto `json:"data"`
}

type GetByIdPhoto struct {
	Status string   `json:"status" example:"success"`
	Data   GetPhoto `json:"data"`
}
