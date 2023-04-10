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

func (commentRepository *commentRepository) Save(ctx context.Context, comment *domain.Comment) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ID, _ := gonanoid.New(16)

	comment.ID = fmt.Sprintf("comment-%s", ID)

	if err = commentRepository.db.WithContext(ctx).Create(&comment).Error; err != nil {
		return err
	}

	return
}

func (commentRepository *commentRepository) Update(ctx context.Context, c domain.Comment, id string) (comment domain.Comment, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = commentRepository.db.WithContext(ctx).First(&comment, &id).Error; err != nil {
		return comment, err
	}

	if err = commentRepository.db.WithContext(ctx).Model(&comment).Updates(comment).Error; err != nil {
		return comment, err
	}

	fmt.Println(comment)

	return comment, nil
}

func (commentRepository *commentRepository) DeleteById(ctx context.Context, id string) (err error) {
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

func (commentRepository *commentRepository) FindById(ctx context.Context, comment *domain.Comment, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = commentRepository.db.WithContext(ctx).First(&comment, &id).Error; err != nil {
		return err
	}

	return
}

func (commentRepository *commentRepository) FindByPhoto(ctx context.Context, comments *[]domain.Comment, photoId string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = commentRepository.db.WithContext(ctx).Where("photo_id = ?", photoId).Preload("Photo", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_id", "title", "photo_url", "caption")
	}).Find(&comments).Error; err != nil {
		return err
	}

	return
}

/*
func (commentRepository *commentRepository) FindAllByUser(ctx context.Context, comments *[]domain.Comment, userID string) (err error) {
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
*/
