package domain

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	ID             string     `gorm:"primaryKey;type:VARCHAR(50)" json:"id"`
	Name           string     `gorm:"type:VARCHAR(50);not null" valid:"required" form:"name" json:"name" example:"Social Media"`
	SocialMediaUrl string     `gorm:"not null" valid:"required" form:"social_media_url" json:"social_media_url" example:"https://www.example.com/social-media"`
	UserID         string     `gorm:"type:VARCHAR(50);not null" json:"user_id"`
	CreatedAt      *time.Time `gorm:"not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `gorm:"not null;autoCreateTime" json:"updated_at,omitempty"`
	User           *User      `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}

func (s *SocialMedia) BeforeCreate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return err
	}

	return
}

func (s *SocialMedia) BeforeUpdate(db *gorm.DB) (err error) {
	if _, err := govalidator.ValidateStruct(s); err != nil {
		return err
	}
	return
}

type SocialMediaUseCase interface {
	Save(context.Context, *SocialMedia) error
	Update(context.Context, SocialMedia, string) (SocialMedia, error)
	DeleteById(context.Context, string) error
	FindAllByUser(context.Context, *[]SocialMedia, string) error
	FindById(context.Context, *SocialMedia, string) error
}

type SocialMediaRepository interface {
	Save(context.Context, *SocialMedia) error
	Update(context.Context, SocialMedia, string) (SocialMedia, error)
	DeleteById(context.Context, string) error
	FindAllByUser(context.Context, *[]SocialMedia, string) error
	FindById(context.Context, *SocialMedia, string) error
}

type AddSocialMedia struct {
	Name           string `json:"name" example:"Example"`
	SocialMediaUrl string `json:"social_media_url" example:"https://www.example.com/johndoe"`
}

type AddedDataSocialMedia struct {
	ID             string     `json:"id" example:"the social media id generated here"`
	Name           string     `json:"name" example:"Example"`
	SocialMediaUrl string     `json:"social_media_url" example:"https://www.example.com/johndoe"`
	UserID         string     `json:"user_id" example:"here is the generated user id"`
	CreatedAt      *time.Time `json:"created_at" example:"the created at generated here"`
}

type AddedSocialMedia struct {
	Status  string               `json:"status" example:"success"`
	Message string               `json:"message" example:"message you if the process has been successful"`
	Data    AddedDataSocialMedia `json:"data"`
}

type UpdateSocialMedia struct {
	Name           string `json:"name" example:"New Example"`
	SocialMediaUrl string `json:"social_media_url" example:"https://www.newexample.com/johndoe"`
}

type UpdatedDataSocialMedia struct {
	ID             string     `json:"id" example:"here is the generated social media id"`
	Name           string     `json:"name" example:"New Example"`
	SocialMediaUrl string     `json:"social_media_url" example:"https://www.newexample.com/johndoe"`
	UserID         string     `json:"user_id" example:"here is the generated user id"`
	UpdatedAt      *time.Time `json:"updated_at" example:"the updated at generated here"`
}

type UpdatedSocialMedia struct {
	Status  string                 `json:"status" example:"success"`
	Message string                 `json:"message" example:"message you if the process has been successful"`
	Data    UpdatedDataSocialMedia `json:"data"`
}

type DeletedSocialMedia struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"your social media has been successfully deleted"`
}

type SocialMedias struct {
	SocialMedias []SocialMedia `json:"social_medias"`
}

type GetDataSocialMedia struct {
	SocialMedias interface{} `json:"social_medias"`
}

type ResponseDataFetchedSocialMedia struct {
	Status  string       `json:"status" example:"success"`
	Message string       `json:"message" example:"message you if the process has been successful"`
	Data    SocialMedias `json:"data"`
}
