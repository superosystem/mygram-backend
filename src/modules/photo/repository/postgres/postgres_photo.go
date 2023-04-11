package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
	"gorm.io/gorm"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *photoRepository {
	return &photoRepository{db}
}

func (photoRepository *photoRepository) Save(ctx context.Context, photo *domain.Photo) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ID, _ := gonanoid.New(16)

	photo.ID = fmt.Sprintf("photo-%s", ID)

	if err := photoRepository.db.WithContext(ctx).Create(&photo).Error; err != nil {
		return err
	}

	if err = photoRepository.db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email")
	}).First(&photo).Error; err != nil {
		return err
	}

	return
}

func (photoRepository *photoRepository) Update(ctx context.Context, p domain.Photo, id string) (photo domain.Photo, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	photo = domain.Photo{}

	if err = photoRepository.db.WithContext(ctx).First(&photo, &id).Error; err != nil {
		return photo, err
	}

	if err = photoRepository.db.WithContext(ctx).Model(&photo).Updates(p).Error; err != nil {
		return photo, err
	}

	if err = photoRepository.db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email")
	}).First(&photo).Error; err != nil {
		return photo, err
	}

	return photo, nil
}

func (photoRepository *photoRepository) DeleteById(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = photoRepository.db.WithContext(ctx).First(&domain.Photo{}, &id).Error; err != nil {
		return err
	}

	if err = photoRepository.db.WithContext(ctx).Delete(&domain.Photo{}, &id).Error; err != nil {
		return err
	}

	return
}

func (photoRepository *photoRepository) FindAll(ctx context.Context, photos *[]domain.Photo) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = photoRepository.db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email")
	}).Find(&photos).Error; err != nil {
		return err
	}

	return
}

func (photoRepository *photoRepository) FindById(ctx context.Context, photo *domain.Photo, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = photoRepository.db.WithContext(ctx).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email")
	}).First(&photo, &id).Error; err != nil {
		return err
	}

	return
}
