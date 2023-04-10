package usecase

import (
	"context"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
)

type photoUseCase struct {
	photoRepository domain.PhotoRepository
}

func NewPhotoUseCase(photoRepository domain.PhotoRepository) *photoUseCase {
	return &photoUseCase{photoRepository}
}

func (photoUseCase *photoUseCase) Save(ctx context.Context, input *domain.AddPhoto, userID string) (photo domain.Photo, err error) {
	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoUrl = input.PhotoUrl
	photo.UserID = userID

	if err = photoUseCase.photoRepository.Save(ctx, &photo); err != nil {
		return photo, err
	}

	if err = photoUseCase.photoRepository.FindById(ctx, &photo, photo.ID); err != nil {
		return photo, err
	}

	return
}

func (photoUseCase *photoUseCase) Update(ctx context.Context, input *domain.UpdatePhoto, id string) (photo domain.Photo, err error) {

	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoUrl = input.PhotoUrl

	if photo, err = photoUseCase.photoRepository.Update(ctx, photo, id); err != nil {
		return photo, err
	}

	if err = photoUseCase.photoRepository.FindById(ctx, &photo, photo.ID); err != nil {
		return photo, err
	}

	return photo, nil
}

func (photoUseCase *photoUseCase) DeleteById(ctx context.Context, id string) (err error) {
	if err = photoUseCase.photoRepository.DeleteById(ctx, id); err != nil {
		return err
	}

	return
}

func (photoUseCase *photoUseCase) FindAll(ctx context.Context, photos *[]domain.Photo) (err error) {
	if err = photoUseCase.photoRepository.FindAll(ctx, photos); err != nil {
		return err
	}

	return
}

func (photoUseCase *photoUseCase) FindById(ctx context.Context, photo *domain.Photo, id string) (err error) {
	if err = photoUseCase.photoRepository.FindById(ctx, photo, id); err != nil {
		return err
	}

	return
}
