package encryption

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
	profileApi "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-recipe/api/model"
	recipeApi "github.com/mephistolie/chefbook-backend-recipe/api/proto/implementation/v1"
	"time"
)

func (s *Service) GetRecipeKeyRequests(recipeId uuid.UUID, userId uuid.UUID) ([]entity.RecipeKeyRequest, error) {
	if err := s.checkRequesterIsRecipeOwner(recipeId, userId); err != nil {
		return nil, err
	}

	requests := s.repo.GetRecipeKeyRequests(recipeId, userId)

	var profileIds []string
	for _, request := range requests {
		profileIds = append(profileIds, request.UserId.String())
	}

	if len(requests) > 0 {
		ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)
		res, err := s.grpc.Profile.GetProfilesMinInfo(ctx, &profileApi.GetProfilesMinInfoRequest{
			ProfileIds: profileIds,
		})
		cancelCtx()

		if err == nil {
			for i, request := range requests {
				if profile, ok := res.Infos[request.UserId.String()]; ok {
					requests[i].UserName = profile.VisibleName
					requests[i].UserAvatar = profile.Avatar
				}
			}
		}
	}

	return requests, nil
}

func (s *Service) GetRecipeKey(recipeId, userId uuid.UUID) *[]byte {
	return s.repo.GetRecipeKey(recipeId, userId)
}

func (s *Service) RequestRecipeKeyAccess(recipeId, userId uuid.UUID) error {
	policy, err := s.grpc.Recipe.GetRecipePolicy(context.Background(), &recipeApi.GetRecipePolicyRequest{
		RecipeId: recipeId.String(),
	})
	if err != nil {
		return err
	}
	if policy.OwnerId == userId.String() {
		return nil
	}
	if policy.Visibility != model.VisibilityLink {
		return fail.GrpcAccessDenied
	}

	return s.repo.CreateRecipeKeyAccessRequest(recipeId, userId)
}

func (s *Service) SetRecipeKey(recipeId, userId uuid.UUID, key []byte, requesterId uuid.UUID) error {
	if err := s.checkRequesterIsRecipeOwner(recipeId, requesterId); err != nil {
		return err
	}

	if userId == requesterId {
		return s.repo.SetRecipeAuthorKey(recipeId, userId, key)
	} else {
		return s.repo.GrantRecipeKeyAccessForUser(recipeId, userId, key)
	}
}

func (s *Service) DeleteRecipeKey(recipeId, userId, requesterId uuid.UUID) error {
	if requesterId != userId {
		if err := s.checkRequesterIsRecipeOwner(recipeId, requesterId); err != nil {
			return err
		}
	}

	if userId == requesterId {
		return s.repo.DeleteRecipeAuthorKey(recipeId, userId)
	} else {
		return s.repo.DeclineRecipeKeyAccessForUser(recipeId, userId)
	}
}

func (s *Service) checkRequesterIsRecipeOwner(recipeId, requesterId uuid.UUID) error {
	policy, err := s.grpc.Recipe.GetRecipePolicy(context.Background(), &recipeApi.GetRecipePolicyRequest{
		RecipeId: recipeId.String(),
	})
	if err != nil {
		return err
	}
	if policy.OwnerId != requesterId.String() {
		return fail.GrpcAccessDenied
	}
	return nil
}
