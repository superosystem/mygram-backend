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

func (commentUseCase *commentUseCase) Save(ctx context.Context, comment *domain.Comment) (err error) {
	if err = commentUseCase.commentRepository.Save(ctx, comment); err != nil {
		return err
	}

	return
}

func (commentUseCase *commentUseCase) Update(ctx context.Context, c domain.Comment, id string) (comment domain.Comment, err error) {
	if comment, err = commentUseCase.commentRepository.Update(ctx, c, id); err != nil {
		return comment, err
	}

	return comment, nil
}

func (commentUseCase *commentUseCase) DeleteById(ctx context.Context, id string) (err error) {
	if err = commentUseCase.commentRepository.DeleteById(ctx, id); err != nil {
		return err
	}

	return
}

/*
func (commentUseCase *commentUseCase) GetAll(ctx context.Context, comments *[]domain.Comment, userID string) (err error) {
	if err = commentUseCase.commentRepository.Fetch(ctx, comments, userID); err != nil {
		return err
	}

	return
}

*/

func (commentUseCase *commentUseCase) FindById(ctx context.Context, comment *domain.Comment, id string) (err error) {
	if err = commentUseCase.commentRepository.FindById(ctx, comment, id); err != nil {
		return err
	}

	return
}

func (commentUseCase *commentUseCase) FindByPhoto(ctx context.Context, comment *[]domain.Comment, id string) (err error) {
	if err = commentUseCase.commentRepository.FindByPhoto(ctx, comment, id); err != nil {
		return err
	}

	return
}
