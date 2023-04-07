package usecase

import (
	"context"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
)

type commentUseCase struct {
	commentRepository domain.CommentRepository
}

func NewCommentUseCase(commentRepository domain.CommentRepository) *commentUseCase {
	return &commentUseCase{commentRepository}
}

func (commentUseCase *commentUseCase) Fetch(ctx context.Context, comments *[]domain.Comment, userID string) (err error) {
	if err = commentUseCase.commentRepository.Fetch(ctx, comments, userID); err != nil {
		return err
	}

	return
}

func (commentUseCase *commentUseCase) Store(ctx context.Context, comment *domain.Comment) (err error) {
	if err = commentUseCase.commentRepository.Store(ctx, comment); err != nil {
		return err
	}

	return
}

func (commentUseCase *commentUseCase) GetByID(ctx context.Context, comment *domain.Comment, id string) (err error) {
	if err = commentUseCase.commentRepository.GetByID(ctx, comment, id); err != nil {
		return err
	}

	return
}

func (commentUseCase *commentUseCase) Update(ctx context.Context, comment domain.Comment, id string) (photo domain.Photo, err error) {
	if photo, err = commentUseCase.commentRepository.Update(ctx, comment, id); err != nil {
		return photo, err
	}

	return photo, nil
}

func (commentUseCase *commentUseCase) Delete(ctx context.Context, id string) (err error) {
	if err = commentUseCase.commentRepository.Delete(ctx, id); err != nil {
		return err
	}

	return
}
