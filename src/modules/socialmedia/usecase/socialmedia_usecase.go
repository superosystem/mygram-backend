package usecase

import (
	"context"

	"github.com/gusrylmubarok/mygram-backend/src/domain"
)

type socialMediaUseCase struct {
	socialMediaRepository domain.SocialMediaRepository
}

func NewSocialMediaUseCase(socialMediaRepository domain.SocialMediaRepository) *socialMediaUseCase {
	return &socialMediaUseCase{socialMediaRepository}
}

func (socialMediaUseCase *socialMediaUseCase) Fetch(ctx context.Context, socialMedias *[]domain.SocialMedia, userID string) (err error) {
	if err = socialMediaUseCase.socialMediaRepository.Fetch(ctx, socialMedias, userID); err != nil {
		return err
	}

	return
}

func (socialMediaUseCase *socialMediaUseCase) Store(ctx context.Context, socialMedia *domain.SocialMedia) (err error) {
	if err = socialMediaUseCase.socialMediaRepository.Store(ctx, socialMedia); err != nil {
		return err
	}

	return
}

func (socialMediaUseCase *socialMediaUseCase) GetByID(ctx context.Context, socialMedia *domain.SocialMedia, id string) (err error) {
	if err = socialMediaUseCase.socialMediaRepository.GetByID(ctx, socialMedia, id); err != nil {
		return err
	}

	return
}

func (socialMediaUseCase *socialMediaUseCase) Update(ctx context.Context, socialMedia domain.SocialMedia, id string) (socmed domain.SocialMedia, err error) {
	if socmed, err = socialMediaUseCase.socialMediaRepository.Update(ctx, socialMedia, id); err != nil {
		return socmed, err
	}

	return socmed, nil
}

func (socialMediaUseCase *socialMediaUseCase) Delete(ctx context.Context, id string) (err error) {
	if err = socialMediaUseCase.socialMediaRepository.Delete(ctx, id); err != nil {
		return err
	}

	return
}
