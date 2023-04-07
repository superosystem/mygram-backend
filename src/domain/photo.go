package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	ID        string     `gorm:"primaryKey;type:VARCHAR(50)" json:"id"`
	Title     string     `gorm:"type:VARCHAR(50);not null" valid:"required" form:"title" json:"title" example:"A Photo Title"`
	Caption   string     `form:"caption" json:"caption"`
	PhotoUrl  string     `gorm:"not null" valid:"required" form:"photo_url" json:"photo_url" example:"https://www.example.com/image.jpg"`
	UserID    string     `gorm:"type:VARCHAR(50);not null" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	CreatedAt *time.Time `gorm:"not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"not null;autoCreateTime" json:"updated_at,omitempty"`
	Comment   *Comment   `json:"-"`
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

type PhotoUseCase interface {
	Fetch(context.Context, *[]Photo) error
	Store(context.Context, *Photo) error
	GetByID(context.Context, *Photo, string) error
	Update(context.Context, Photo, string) (Photo, error)
	Delete(context.Context, string) error
}

type PhotoRepository interface {
	Fetch(context.Context, *[]Photo) error
	Store(context.Context, *Photo) error
	GetByID(context.Context, *Photo, string) error
	Update(context.Context, Photo, string) (Photo, error)
	Delete(context.Context, string) error
}
