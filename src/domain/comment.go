package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	ID        string     `gorm:"primaryKey;type:VARCHAR(50)"`
	UserID    string     `gorm:"type:VARCHAR(50);not null"`
	PhotoID   string     `gorm:"type:VARCHAR(50);not null"`
	Message   string     `gorm:"not null" valid:"required"`
	CreatedAt *time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"not null;autoCreateTime"`
	User      *User      `gorm:"foreignKey:UserID;constraint:opUpdate:CASCADE,onDelete:CASCADE"`
	Photo     *Photo     `gorm:"foreignKey:PhotoID;constraint:opUpdate:CASCADE,onDelete:CASCADE"`
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
	Save(context.Context, *Comment) error
	Update(context.Context, Comment, string) (Comment, error)
	DeleteById(context.Context, string) error
	FindById(context.Context, *Comment, string) error
	FindByPhoto(context.Context, *[]Comment, string) error
}

type CommentUseCase interface {
	Save(context.Context, *Comment) error
	Update(context.Context, Comment, string) (Comment, error)
	DeleteById(context.Context, string) error
	FindById(context.Context, *Comment, string) error
	FindByPhoto(context.Context, *[]Comment, string) error
}

// Represents for request add comment
type AddComment struct {
	PhotoID string `json:"photo_id"  form:"photoId" example:"photo-123"`
	Message string `json:"message" form:"message" example:"A comment"`
}

// Represents for added comment
type AddedComment struct {
	ID        string     `json:"id" example:"here is the generated comment id"`
	Message   string     `json:"message" form:"message" example:"A comment"`
	UserID    string     `json:"user_id" form:"user_id" example:"here is the generated user id"`
	PhotoID   string     `json:"photo_id" form:"photo_id" example:"here is the generated photo id"`
	CreatedAt *time.Time `json:"created_at" example:"the created at generated here"`
}

// Represents for response added comment
type ResponseAddedComment struct {
	Status string       `json:"status" example:"success"`
	Data   AddedComment `json:"data"`
}

// Represents for request update comment
type UpdateComment struct {
	Message string `json:"message" example:"A new comment"`
}

// Represents for updated comment
type UpdatedComment struct {
	ID        string     `json:"id"`
	Message   string     `json:"message" form:"message" example:"A comment"`
	PhotoID   string     `json:"photo_id" form:"photo_id" example:"here is the generated photo id"`
	UserID    string     `json:"user_id" form:"user_id" example:"here is the generated user id"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Represents for response updated comment
type ResponseUpdatedComment struct {
	Status string         `json:"status" example:"success"`
	Data   UpdatedComment `json:"data"`
}

// Represents for getting comment
type GetComment struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PhotoID   string     `json:"photo_id"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	User      *User      `json:"user"`
	Photo     *Photo     `json:"photo"`
}

// Represents for response get comment
type ResponseGetComment struct {
	Status string       `json:"status" example:"success"`
	Data   []GetComment `json:"data"`
}

type ResponseMessageDeletedComment struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your comment has been successfully deleted"`
}

type ResponseMessage struct {
	Status string `json:"status" example:"fail"`
	Data   string `json:"data" example:"the error explained here"`
}
