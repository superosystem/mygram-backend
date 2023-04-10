package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"gorm.io/gorm"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type socialMediaRepository struct {
	db *gorm.DB
}

func NewSocialMediaRepository(db *gorm.DB) *socialMediaRepository {
	return &socialMediaRepository{db}
}

func (socialMediaRepository *socialMediaRepository) Save(ctx context.Context, socialMedia *domain.SocialMedia) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	ID, _ := gonanoid.New(16)

	socialMedia.ID = fmt.Sprintf("socialmedia-%s", ID)

	if err = socialMediaRepository.db.WithContext(ctx).Create(&socialMedia).Error; err != nil {
		return err
	}

	return
}

func (socialMediaRepository *socialMediaRepository) Update(ctx context.Context, socialMedia domain.SocialMedia, id string) (socmed domain.SocialMedia, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	socmed = domain.SocialMedia{}

	if err = socialMediaRepository.db.WithContext(ctx).First(&socmed, &id).Error; err != nil {
		return socmed, err
	}

	if err = socialMediaRepository.db.WithContext(ctx).Model(&socmed).Updates(socialMedia).Error; err != nil {
		return socmed, err
	}

	return socmed, nil
}

func (socialMediaRepository *socialMediaRepository) DeleteById(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = socialMediaRepository.db.WithContext(ctx).First(&domain.SocialMedia{}, &id).Error; err != nil {
		return err
	}

	if err = socialMediaRepository.db.WithContext(ctx).Delete(&domain.SocialMedia{}, &id).Error; err != nil {
		return err
	}

	return
}

func (socialMediaRepository *socialMediaRepository) FindAllByUser(ctx context.Context, socialMedias *[]domain.SocialMedia, userID string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = socialMediaRepository.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email", "Username")
	}).Find(&socialMedias).Error; err != nil {
		return err
	}

	return
}

func (socialMediaRepository *socialMediaRepository) FindById(ctx context.Context, socialMedia *domain.SocialMedia, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = socialMediaRepository.db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email")
	}).First(&socialMedia, &id).Error; err != nil {
		return err
	}

	return
}
