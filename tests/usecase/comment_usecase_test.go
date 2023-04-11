package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gusrylmubarok/mygram-backend/src/domain"
	mocks "github.com/gusrylmubarok/mygram-backend/src/domain/mocks/repository"
	commentUseCase "github.com/gusrylmubarok/mygram-backend/src/modules/comment/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveComment(t *testing.T) {
	now := time.Now()
	mockAddedComment := domain.Comment{
		ID:        "comment-123",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		Message:   "A comment",
		CreatedAt: &now,
	}

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success add comment", func(t *testing.T) {
		tempMockAddComment := domain.Comment{
			Message: "A comment",
			PhotoID: "photo-123",
		}

		tempMockAddComment.ID = "comment-123"

		mockCommentRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil).Once()

		err := commentUseCase.Save(context.Background(), &tempMockAddComment)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddComment)

		assert.NoError(t, err)
		assert.Equal(t, mockAddedComment.ID, tempMockAddComment.ID)
		assert.Equal(t, mockAddedComment.Message, tempMockAddComment.Message)
		assert.Equal(t, mockAddedComment.PhotoID, tempMockAddComment.PhotoID)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail add comment with empty message", func(t *testing.T) {
		tempMockAddComment := domain.Comment{
			Message: "",
			PhotoID: "photo-123",
		}

		tempMockAddComment.ID = "comment-123"

		mockCommentRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil).Once()

		err := commentUseCase.Save(context.Background(), &tempMockAddComment)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddComment)

		assert.Error(t, err)
		assert.Equal(t, mockAddedComment.ID, tempMockAddComment.ID)
		assert.NotEqual(t, mockAddedComment.Message, tempMockAddComment.Message)
		assert.Equal(t, mockAddedComment.PhotoID, tempMockAddComment.PhotoID)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail add comment with empty photo id", func(t *testing.T) {
		tempMockAddComment := domain.Comment{
			Message: "A comment",
			PhotoID: "",
		}

		tempMockAddComment.ID = "comment-123"

		mockCommentRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(errors.New("fail")).Once()

		err := commentUseCase.Save(context.Background(), &tempMockAddComment)

		assert.Error(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddComment)

		assert.NoError(t, err)
		assert.Equal(t, mockAddedComment.ID, tempMockAddComment.ID)
		assert.Equal(t, mockAddedComment.Message, tempMockAddComment.Message)
		assert.NotEqual(t, mockAddedComment.PhotoID, tempMockAddComment.PhotoID)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should add comment with not contain needed property", func(t *testing.T) {
		tempMockAddComment := domain.Comment{
			PhotoID: "photo-123",
		}

		tempMockAddComment.ID = "comment-123"

		mockCommentRepository.On("Save", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil).Once()

		err := commentUseCase.Save(context.Background(), &tempMockAddComment)

		assert.NoError(t, err)

		_, err = govalidator.ValidateStruct(tempMockAddComment)

		assert.Error(t, err)
		assert.Equal(t, mockAddedComment.ID, tempMockAddComment.ID)
		assert.NotEqual(t, mockAddedComment.Message, tempMockAddComment.Message)
		assert.Equal(t, mockAddedComment.PhotoID, tempMockAddComment.PhotoID)
		mockCommentRepository.AssertExpectations(t)
	})
}

func TestUpdateComment(t *testing.T) {
	now := time.Now()
	mockUpdatedComment := domain.Comment{
		ID:        "comment-123",
		Message:   "A new comment",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		UpdatedAt: &now,
	}

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success update comment correctly", func(t *testing.T) {
		tempMockCommentID := "comment-123"
		tempMockUpdateComment := domain.Comment{
			Message: "A new comment",
		}

		mockCommentRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Comment"), mock.AnythingOfType("string")).Return(mockUpdatedComment, nil).Once()

		comment, err := commentUseCase.Update(context.Background(), tempMockUpdateComment, tempMockCommentID)

		assert.NoError(t, err)

		tempMockUpdatedComment := domain.Comment{
			ID:        "comment-123",
			Message:   "A new comment",
			UserID:    "user-123",
			PhotoID:   "photo-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdateComment)

		assert.NoError(t, err)
		assert.Equal(t, comment, tempMockUpdatedComment)
		assert.Equal(t, mockUpdatedComment.ID, tempMockUpdatedComment.ID)
		assert.Equal(t, mockUpdatedComment.Message, tempMockUpdatedComment.Message)
		assert.Equal(t, mockUpdatedComment.UserID, tempMockUpdatedComment.UserID)
		assert.Equal(t, mockUpdatedComment.PhotoID, tempMockUpdatedComment.PhotoID)
		assert.Equal(t, mockUpdatedComment.UpdatedAt, tempMockUpdatedComment.UpdatedAt)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail update comment with empty message", func(t *testing.T) {
		tempMockCommentID := "Comment-123"
		tempMockUpdateComment := domain.Comment{
			Message: "",
		}

		mockCommentRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Comment"), mock.AnythingOfType("string")).Return(mockUpdatedComment, nil).Once()

		comment, err := commentUseCase.Update(context.Background(), tempMockUpdateComment, tempMockCommentID)

		assert.NoError(t, err)

		tempMockUpdatedComment := domain.Comment{
			ID:        "comment-123",
			Message:   "A new comment",
			UserID:    "user-123",
			PhotoID:   "photo-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdateComment)

		assert.Error(t, err)
		assert.Equal(t, comment, tempMockUpdatedComment)
		assert.Equal(t, mockUpdatedComment.ID, tempMockUpdatedComment.ID)
		assert.Equal(t, mockUpdatedComment.Message, tempMockUpdatedComment.Message)
		assert.Equal(t, mockUpdatedComment.UserID, tempMockUpdatedComment.UserID)
		assert.Equal(t, mockUpdatedComment.PhotoID, tempMockUpdatedComment.PhotoID)
		assert.Equal(t, mockUpdatedComment.UpdatedAt, tempMockUpdatedComment.UpdatedAt)
		mockCommentRepository.AssertExpectations(t)
	})
	t.Run("should fail update comment with not contain property", func(t *testing.T) {
		tempMockCommentID := "comment-123"
		tempMockUpdateComment := domain.Comment{}

		mockCommentRepository.On("Update", mock.Anything, mock.AnythingOfType("domain.Comment"), mock.AnythingOfType("string")).Return(mockUpdatedComment, nil).Once()

		comment, err := commentUseCase.Update(context.Background(), tempMockUpdateComment, tempMockCommentID)

		assert.NoError(t, err)

		tempMockUpdatedComment := domain.Comment{
			ID:        "comment-123",
			Message:   "A new comment",
			UserID:    "user-123",
			PhotoID:   "photo-123",
			UpdatedAt: &now,
		}

		_, err = govalidator.ValidateStruct(tempMockUpdateComment)

		assert.Error(t, err)
		assert.Equal(t, comment, tempMockUpdatedComment)
		assert.Equal(t, mockUpdatedComment.ID, tempMockUpdatedComment.ID)
		assert.Equal(t, mockUpdatedComment.Message, tempMockUpdatedComment.Message)
		assert.Equal(t, mockUpdatedComment.UserID, tempMockUpdatedComment.UserID)
		assert.Equal(t, mockUpdatedComment.PhotoID, tempMockUpdatedComment.PhotoID)
		assert.Equal(t, mockUpdatedComment.UpdatedAt, tempMockUpdatedComment.UpdatedAt)
		mockCommentRepository.AssertExpectations(t)
	})
}

func TestDeleteComment(t *testing.T) {
	now := time.Now()
	mockComment := domain.Comment{
		ID:        "comment-123",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		Message:   "A message",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success delete comment correctly", func(t *testing.T) {
		mockCommentRepository.On("DeleteById", mock.Anything, mock.AnythingOfType("string")).Return(nil).Once()

		err := commentUseCase.DeleteById(context.Background(), mockComment.ID)

		assert.NoError(t, err)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail delete comment with not found Comment", func(t *testing.T) {
		mockCommentRepository.On("DeleteById", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := commentUseCase.DeleteById(context.Background(), "comment-234")

		assert.Error(t, err)
		mockCommentRepository.AssertExpectations(t)
	})
}

func TestFindAllByUser(t *testing.T) {
	now := time.Now()
	mockComment := domain.Comment{
		ID:        "comment-123",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		Message:   "A message",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockComments := make([]domain.Comment, 0)

	mockComments = append(mockComments, mockComment)

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success find all comments by user", func(t *testing.T) {
		mockCommentRepository.On("FindAllByUser", mock.Anything, mock.AnythingOfType("*[]domain.Comment"), mock.AnythingOfType("string")).Return(nil).Once()

		err := commentUseCase.FindAllByUser(context.Background(), &mockComments, mockComment.UserID)

		assert.NoError(t, err)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail find all comment cause empty by user", func(t *testing.T) {
		mockCommentRepository.On("FindAllByUser", mock.Anything, mock.AnythingOfType("*[]domain.Comment"), mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := commentUseCase.FindAllByUser(context.Background(), &mockComments, "user-345")

		assert.Error(t, err)
		mockCommentRepository.AssertExpectations(t)
	})
}

func TestFindAllByPhoto(t *testing.T) {
	now := time.Now()
	mockComment := domain.Comment{
		ID:        "comment-123",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		Message:   "A message",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockComments := make([]domain.Comment, 0)

	mockComments = append(mockComments, mockComment)

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success find all comments by photo", func(t *testing.T) {
		mockCommentRepository.On("FindAllByPhoto", mock.Anything, mock.AnythingOfType("*[]domain.Comment"), mock.AnythingOfType("string")).Return(nil).Once()

		err := commentUseCase.FindAllByPhoto(context.Background(), &mockComments, mockComment.PhotoID)

		assert.NoError(t, err)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail find all comment cause empty by photo", func(t *testing.T) {
		mockCommentRepository.On("FindAllByPhoto", mock.Anything, mock.AnythingOfType("*[]domain.Comment"), mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := commentUseCase.FindAllByPhoto(context.Background(), &mockComments, "photo-345")

		assert.Error(t, err)
		mockCommentRepository.AssertExpectations(t)
	})
}

func TestFindCommentById(t *testing.T) {
	now := time.Now()
	mockComment := domain.Comment{
		ID:        "comment-123",
		UserID:    "user-123",
		PhotoID:   "photo-123",
		Message:   "A message",
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	mockCommentRepository := new(mocks.CommentRepository)
	commentUseCase := commentUseCase.NewCommentUseCase(mockCommentRepository)

	t.Run("should success find a comment", func(t *testing.T) {
		mockCommentRepository.On("FindById", mock.Anything, mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(nil).Once()

		mockPhotoId := "comment-123"

		err := commentUseCase.FindById(context.Background(), &mockComment, mockPhotoId)

		assert.NoError(t, err)
		mockCommentRepository.AssertExpectations(t)
	})

	t.Run("should fail find a comment with id not found", func(t *testing.T) {
		mockCommentRepository.On("FindById", mock.Anything, mock.AnythingOfType("*domain.Comment"), mock.AnythingOfType("string")).Return(errors.New("fail")).Once()

		err := commentUseCase.FindById(context.Background(), &mockComment, "comment-345")

		assert.Error(t, err)
		mockCommentRepository.AssertExpectations(t)
	})
}
