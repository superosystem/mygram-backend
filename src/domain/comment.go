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
	User      *User      `gorm:"foreignKey:UserID;constraint:opUpdate:CASCADE,onDelete:CASCADE" json:"user,omitempty"`
	Photo     *Photo     `gorm:"foreignKey:PhotoID;constraint:opUpdate:CASCADE,onDelete:CASCADE" json:"photo,omitempty"`
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
	FindAllByUser(context.Context, *[]Comment, string) error
	FindAllByPhoto(context.Context, *[]Comment, string) error
	FindById(context.Context, *Comment, string) error
}

type CommentUseCase interface {
	Save(context.Context, *Comment) error
	Update(context.Context, Comment, string) (Comment, error)
	DeleteById(context.Context, string) error
	FindAllByUser(context.Context, *[]Comment, string) error
	FindAllByPhoto(context.Context, *[]Comment, string) error
	FindById(context.Context, *Comment, string) error
}

// Represents for request add comment
type AddComment struct {
	Message string `json:"message" form:"message" example:"A comment"`
	PhotoID string `json:"photo_id"  form:"photoId" example:"photo-123"`
	UserID  string
}

// Represents for added comment
type AddedDataComment struct {
	ID        string     `json:"id" example:"here is the generated comment id"`
	Message   string     `json:"message" form:"message" example:"A comment"`
	User      *GetUser   `json:"user"`
	Photo     *GetPhoto  `json:"photo"`
	CreatedAt *time.Time `json:"created_at" example:"the created at generated here"`
}

// Represents for response added comment
type AddedComment struct {
	Status  string           `json:"status" example:"success"`
	Message string           `json:"message" example:"message you if the process has been successful"`
	Data    AddedDataComment `json:"data"`
}

// Represents for request update comment
type UpdateComment struct {
	Message string `json:"message" example:"A new comment"`
	UserID  string
}

// Represents for updated comment
type UpdatedDataComment struct {
	ID        string     `json:"id"`
	Message   string     `json:"message" form:"message" example:"A comment"`
	User      *GetUser   `json:"user"`
	Photo     *GetPhoto  `json:"photo"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Represents for response updated comment
type UpdatedComment struct {
	Status string             `json:"status" example:"success"`
	Data   UpdatedDataComment `json:"data"`
}

// Represents for response delted comment
type DeletedComment struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your comment has been successfully deleted"`
}

// Represents for getting comment
type GetComment struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	User      *User      `json:"user"`
	Photo     *Photo     `json:"photo"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Represents for response get comment
type GetAllComments struct {
	Status  string    `json:"status" example:"success"`
	Message string    `json:"message" example:"message you if the process has been successful"`
	Data    []Comment `json:"data"`
}

// Represents for response get comment
type GetAComment struct {
	Status  string  `json:"status" example:"success"`
	Message string  `json:"message" example:"message you if the process has been successful"`
	Data    Comment `json:"data"`
}
