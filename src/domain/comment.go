package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	ID        string     `gorm:"primaryKey;type:VARCHAR(50)" json:"id"`
	UserID    string     `gorm:"type:VARCHAR(50);not null" json:"user_id"`
	PhotoID   string     `gorm:"type:VARCHAR(50);not null" form:"photo_id" json:"photo_id"`
	Message   string     `gorm:"not null" valid:"required" form:"message" json:"message" example:"A comment"`
	CreatedAt *time.Time `gorm:"not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt *time.Time `gorm:"not null;autoCreateTime" json:"updated_at,omitempty"`
	User      *User      `gorm:"foreignKey:UserID;constraint:opUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	Photo     *Photo     `gorm:"foreignKey:PhotoID;constraint:opUpdate:CASCADE,onDelete:CASCADE" json:"photo"`
}

func (c *Comment) BeforeCreate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}

	return
}

func (c *Comment) BeforeUpdate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}
	return
}

type CommentRepository interface {
	Fetch(context.Context, *[]Comment, string) error
	Store(context.Context, *Comment) error
	GetByID(context.Context, *Comment, string) error
	Update(context.Context, Comment, string) (Photo, error)
	Delete(context.Context, string) error
}

type CommentUseCase interface {
	Fetch(context.Context, *[]Comment, string) error
	Store(context.Context, *Comment) error
	GetByID(context.Context, *Comment, string) error
	Update(context.Context, Comment, string) (Photo, error)
	Delete(context.Context, string) error
}
