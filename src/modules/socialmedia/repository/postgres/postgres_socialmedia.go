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

func (socialMediaRepository *socialMediaRepository) Fetch(ctx context.Context, socialMedias *[]domain.SocialMedia, userID string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = socialMediaRepository.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("ID", "Email", "Username", "ProfileImageUrl")
	}).Find(&socialMedias).Error; err != nil {
		return err
	}

	return
}

func (socialMediaRepository *socialMediaRepository) Store(ctx context.Context, socialMedia *domain.SocialMedia) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	ID, _ := gonanoid.New(16)

	socialMedia.ID = fmt.Sprintf("socialmedia-%s", ID)

	if err = socialMediaRepository.db.WithContext(ctx).Create(&socialMedia).Error; err != nil {
		return err
	}

	return
}

func (socialMediaRepository *socialMediaRepository) GetByID(ctx context.Context, socialMedia *domain.SocialMedia, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = socialMediaRepository.db.WithContext(ctx).First(&socialMedia, &id).Error; err != nil {
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

func (socialMediaRepository *socialMediaRepository) Delete(ctx context.Context, id string) (err error) {
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
