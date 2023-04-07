package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *commentRepository {
	return &commentRepository{db}
}

func (commentRepository *commentRepository) Fetch(ctx context.Context, comments *[]domain.Comment, userID string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = commentRepository.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "email", "username", "profile_image_url")
	}).Preload("Photo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_id", "title", "photo_url", "caption")
	}).Find(&comments).Error; err != nil {
		return err
	}

	return
}

func (commentRepository *commentRepository) Store(ctx context.Context, comment *domain.Comment) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	ID, _ := gonanoid.New(16)

	comment.ID = fmt.Sprintf("comment-%s", ID)

	if err = commentRepository.db.WithContext(ctx).Create(&comment).Error; err != nil {
		return err
	}

	return
}

func (commentRepository *commentRepository) GetByID(ctx context.Context, comment *domain.Comment, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = commentRepository.db.WithContext(ctx).First(&comment, &id).Error; err != nil {
		return err
	}

	return
}

func (commentRepository *commentRepository) Update(ctx context.Context, comment domain.Comment, id string) (photo domain.Photo, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	c := domain.Comment{}

	photo = domain.Photo{}

	if err = commentRepository.db.WithContext(ctx).First(&c, &id).Error; err != nil {
		return photo, err
	}

	if err = commentRepository.db.WithContext(ctx).Model(&c).Updates(comment).Error; err != nil {
		return photo, err
	}

	if err = commentRepository.db.WithContext(ctx).First(&photo, comment.PhotoID).Error; err != nil {
		return photo, err
	}

	return photo, nil
}

func (commentRepository *commentRepository) Delete(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = commentRepository.db.WithContext(ctx).First(&domain.Comment{}, &id).Error; err != nil {
		return err
	}

	if err = commentRepository.db.WithContext(ctx).Delete(&domain.Comment{}, &id).Error; err != nil {
		return err
	}

	return
}
