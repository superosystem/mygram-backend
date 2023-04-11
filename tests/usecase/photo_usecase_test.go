package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	mocks "github.com/gusrylmubarok/mygram-backend/src/domain/mocks/repository"
	photoUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/photo/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSavePhoto(t *testing.T) {
	now := time.Now()
	mockAddedPhoto := domain.Photo{
		ID:        "photo-123",
		Title:     "A Title",
		Caption:   "A caption",
		PhotoUrl:  "https://www.example.com/image.jpg",
		UserID:    "user-123",
		CreatedAt: &now,
	}

	mockPhotoRepository := new(mocks.PhotoRepository)
	photoUseCase := photoUseCase.NewPhotoUseCase(mockPhotoRepository)

	t.Run("should success add photo", func(t *testing.T) {
		tempMockAddPhoto := domain.Photo{
			Title:    "A Title",
			Caption:  "A caption",
			PhotoUrl: "https://www.example.com/image.jpg",
		}

		tempMockAddPhoto.ID = "photo-123"

		mockPhotoRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Photo")).Return(nil).Once()

		err := photoUseCase.Save(context.Background(), &tempMockAddPhoto)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddPhoto)

		assert.NoError(t, err)
		assert.Equal(t, mockAddedPhoto.ID, tempMockAddPhoto.ID)
		assert.Equal(t, mockAddedPhoto.Title, tempMockAddPhoto.Title)
		assert.Equal(t, mockAddedPhoto.Caption, tempMockAddPhoto.Caption)
		assert.Equal(t, mockAddedPhoto.PhotoUrl, tempMockAddPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should fail add photo with empty title", func(t *testing.T) {
		tempMockAddPhoto := domain.Photo{
			Title:    "",
			Caption:  "A caption",
			PhotoUrl: "https://www.example.com/image.jpg",
		}

		tempMockAddPhoto.ID = "photo-123"

		mockPhotoRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Photo")).Return(nil).Once()

		err := photoUseCase.Save(context.Background(), &tempMockAddPhoto)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddPhoto)

		assert.Error(t, err)
		assert.Equal(t, mockAddedPhoto.ID, tempMockAddPhoto.ID)
		assert.NotEqual(t, mockAddedPhoto.Title, tempMockAddPhoto.Title)
		assert.Equal(t, mockAddedPhoto.Caption, tempMockAddPhoto.Caption)
		assert.Equal(t, mockAddedPhoto.PhotoUrl, tempMockAddPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should fail add photo with empty photo url", func(t *testing.T) {
		tempMockAddPhoto := domain.Photo{
			Title:    "A Title",
			Caption:  "A caption",
			PhotoUrl: "",
		}

		tempMockAddPhoto.ID = "photo-123"

		mockPhotoRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Photo")).Return(nil).Once()

		err := photoUseCase.Save(context.Background(), &tempMockAddPhoto)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddPhoto)

		assert.Error(t, err)
		assert.Equal(t, mockAddedPhoto.ID, tempMockAddPhoto.ID)
		assert.Equal(t, mockAddedPhoto.Title, tempMockAddPhoto.Title)
		assert.Equal(t, mockAddedPhoto.Caption, tempMockAddPhoto.Caption)
		assert.NotEqual(t, mockAddedPhoto.PhotoUrl, tempMockAddPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should faile add photo with not contain needed property", func(t *testing.T) {
		tempMockAddPhoto := domain.Photo{
			Title:   "A Title",
			Caption: "A caption",
		}

		tempMockAddPhoto.ID = "photo-123"

		mockPhotoRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Photo")).Return(nil).Once()

		err := photoUseCase.Save(context.Background(), &tempMockAddPhoto)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddPhoto)

		assert.Error(t, err)
		assert.Equal(t, mockAddedPhoto.ID, tempMockAddPhoto.ID)
		assert.Equal(t, mockAddedPhoto.Title, tempMockAddPhoto.Title)
		assert.Equal(t, mockAddedPhoto.Caption, tempMockAddPhoto.Caption)
		assert.NotEqual(t, mockAddedPhoto.PhotoUrl, tempMockAddPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})
}

func TestUpdatePhoto(t *testing.T) {
	now := time.Now()
	mockUpdatedPhoto := domain.Photo{
		ID:        "photo-123",
		Title:     "A New Title",
		Caption:   "A new caption",
		PhotoUrl:  "https://www.example.com/new-image.jpg",
		UserID:    "user-123",
		UpdatedAt: &now,
	}

	mockPhotoRepository := new(mocks.PhotoRepository)
	photoUseCase := photoUseCase.NewPhotoUseCase(mockPhotoRepository)

	t.Run("should success update photo", func(t *testing.T) {
		tempMockPhotoID := "photo-123"
		tempMockUpdatePhoto := domain.Photo{
			Title:    "A New Title",
			Caption:  "A new caption",
			PhotoUrl: "https://www.example.com/new-image.jpg",
		}

		mockPhotoRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Photo"), mock.AnythingOfType("string")).Return(mockUpdatedPhoto, nil).Once()

		photo, err := photoUseCase.Update(context.Background(), tempMockUpdatePhoto, tempMockPhotoID)

		assert.NoError(t, err)

		tempMockUpdatedPhoto := domain.Photo{
			ID:        tempMockPhotoID,
			Title:     tempMockUpdatePhoto.Title,
			Caption:   tempMockUpdatePhoto.Caption,
			PhotoUrl:  tempMockUpdatePhoto.PhotoUrl,
			UserID:    "user-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdatedPhoto)

		assert.NoError(t, err)
		assert.Equal(t, photo, tempMockUpdatedPhoto)
		assert.Equal(t, mockUpdatedPhoto.Title, tempMockUpdatePhoto.Title)
		assert.Equal(t, mockUpdatedPhoto.Caption, tempMockUpdatePhoto.Caption)
		assert.Equal(t, mockUpdatedPhoto.PhotoUrl, tempMockUpdatedPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should success update photo with empty title", func(t *testing.T) {
		tempMockPhotoID := "photo-123"
		tempMockUpdatePhoto := domain.Photo{
			Title:    "",
			Caption:  "A new caption",
			PhotoUrl: "https://www.example.com/new-image.jpg",
		}

		mockPhotoRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Photo"), mock.AnythingOfType("string")).Return(mockUpdatedPhoto, nil).Once()

		photo, err := photoUseCase.Update(context.Background(), tempMockUpdatePhoto, tempMockPhotoID)

		assert.NoError(t, err)

		tempMockUpdatedPhoto := domain.Photo{
			ID:        tempMockPhotoID,
			Title:     tempMockUpdatePhoto.Title,
			Caption:   tempMockUpdatePhoto.Caption,
			PhotoUrl:  tempMockUpdatePhoto.PhotoUrl,
			UserID:    "user-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdatedPhoto)

		assert.Error(t, err)
		assert.NotEqual(t, photo, tempMockUpdatedPhoto)
		assert.NotEqual(t, mockUpdatedPhoto.Title, tempMockUpdatePhoto.Title)
		assert.Equal(t, mockUpdatedPhoto.Caption, tempMockUpdatePhoto.Caption)
		assert.Equal(t, mockUpdatedPhoto.PhotoUrl, tempMockUpdatedPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should success update photo with empty photo url", func(t *testing.T) {
		tempMockPhotoID := "photo-123"
		tempMockUpdatePhoto := domain.Photo{
			Title:    "A New Title",
			Caption:  "A new caption",
			PhotoUrl: "",
		}

		mockPhotoRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Photo"), mock.AnythingOfType("string")).Return(mockUpdatedPhoto, nil).Once()

		photo, err := photoUseCase.Update(context.Background(), tempMockUpdatePhoto, tempMockPhotoID)

		assert.NoError(t, err)

		tempMockUpdatedPhoto := domain.Photo{
			ID:        tempMockPhotoID,
			Title:     tempMockUpdatePhoto.Title,
			Caption:   tempMockUpdatePhoto.Caption,
			PhotoUrl:  tempMockUpdatePhoto.PhotoUrl,
			UserID:    "user-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdatedPhoto)

		assert.Error(t, err)
		assert.NotEqual(t, photo, tempMockUpdatedPhoto)
		assert.Equal(t, mockUpdatedPhoto.Title, tempMockUpdatePhoto.Title)
		assert.Equal(t, mockUpdatedPhoto.Caption, tempMockUpdatePhoto.Caption)
		assert.NotEqual(t, mockUpdatedPhoto.PhotoUrl, tempMockUpdatedPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should success update photo with empty title and photo url", func(t *testing.T) {
		tempMockPhotoID := "photo-123"
		tempMockUpdatePhoto := domain.Photo{
			Title:    "",
			Caption:  "A new caption",
			PhotoUrl: "",
		}

		mockPhotoRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Photo"), mock.AnythingOfType("string")).Return(mockUpdatedPhoto, nil).Once()

		photo, err := photoUseCase.Update(context.Background(), tempMockUpdatePhoto, tempMockPhotoID)

		assert.NoError(t, err)

		tempMockUpdatedPhoto := domain.Photo{
			ID:        tempMockPhotoID,
			Title:     tempMockUpdatePhoto.Title,
			Caption:   tempMockUpdatePhoto.Caption,
			PhotoUrl:  tempMockUpdatePhoto.PhotoUrl,
			UserID:    "user-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdatedPhoto)

		assert.Error(t, err)
		assert.NotEqual(t, photo, tempMockUpdatedPhoto)
		assert.NotEqual(t, mockUpdatedPhoto.Title, tempMockUpdatePhoto.Title)
		assert.Equal(t, mockUpdatedPhoto.Caption, tempMockUpdatePhoto.Caption)
		assert.NotEqual(t, mockUpdatedPhoto.PhotoUrl, tempMockUpdatedPhoto.PhotoUrl)
		mockPhotoRepository.AssertExpectations(t)
	})
}

func TestDeletePhoto(t *testing.T) {
	mockPhoto := domain.Photo{
		ID:       "photo-123",
		Title:    "A Title",
		Caption:  "A caption",
		PhotoUrl: "https://www.example.com/image.jpg",
		UserID:   "user-123",
	}

	mockPhotoRepository := new(mocks.PhotoRepository)
	photoUseCase := photoUseCase.NewPhotoUseCase(mockPhotoRepository)

	t.Run("should success delete photo", func(t *testing.T) {
		mockPhotoRepository.On("DeleteById", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()

		err := photoUseCase.DeleteById(context.Background(), mockPhoto.ID)

		assert.NoError(t, err)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should fail delete photo with not found photo", func(t *testing.T) {
		mockPhotoRepository.On("DeleteById", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := photoUseCase.DeleteById(context.Background(), "photo-234")

		assert.Error(t, err)
		mockPhotoRepository.AssertExpectations(t)
	})
}

func TestFindAllPhoto(t *testing.T) {
	mockPhoto := domain.Photo{
		ID:       "photo-123",
		Title:    "A Title",
		Caption:  "A caption",
		PhotoUrl: "https://www.example.com/image.jpg",
		UserID:   "user-123",
	}

	mockPhotos := make([]domain.Photo, 0)

	mockPhotos = append(mockPhotos, mockPhoto)

	mockPhotoRepository := new(mocks.PhotoRepository)
	photoUseCase := photoUseCase.NewPhotoUseCase(mockPhotoRepository)

	t.Run("should success find all photos", func(t *testing.T) {
		mockPhotoRepository.On("FindAll", mock.Anything, mock.AnythingOfType("*[]domain.Photo")).Return(nil).Once()

		err := photoUseCase.FindAll(context.Background(), &mockPhotos)

		assert.NoError(t, err)
	})
}

func TestFindPhotoById(t *testing.T) {
	mockPhotoId := "photo-123"
	mockPhoto := domain.Photo{
		ID:       "photo-123",
		Title:    "A Title",
		Caption:  "A caption",
		PhotoUrl: "https://www.example.com/image.jpg",
		UserID:   "user-123",
	}

	mockPhotoRepository := new(mocks.PhotoRepository)
	photoUseCase := photoUseCase.NewPhotoUseCase(mockPhotoRepository)

	t.Run("should success find a photo", func(t *testing.T) {
		mockPhotoRepository.On("FindById", mock.Anything, mock.AnythingOfType("*domain.Photo"), mock.AnythingOfType("string")).Return(nil).Once()

		err := photoUseCase.FindById(context.Background(), &mockPhoto, mockPhotoId)

		assert.NoError(t, err)
		mockPhotoRepository.AssertExpectations(t)
	})

	t.Run("should fail find a photo with id not found", func(t *testing.T) {
		mockPhotoRepository.On("FindById", mock.Anything, mock.AnythingOfType("*domain.Photo"), mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := photoUseCase.FindById(context.Background(), &mockPhoto, "photo-345")

		assert.Error(t, err)
		mockPhotoRepository.AssertExpectations(t)
	})
}
