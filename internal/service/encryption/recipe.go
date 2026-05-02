package encryption

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
	encryptionFail "github.com/mephistolie/chefbook-backend-encryption/internal/entity/fail"
	profileApi "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	recipeModel "github.com/mephistolie/chefbook-backend-recipe/api/model"
	recipeApi "github.com/mephistolie/chefbook-backend-recipe/api/proto/implementation/v1"
)

func (s *Service) GetRecipeKeyRequests(ctx context.Context, recipeId uuid.UUID, userId uuid.UUID) ([]entity.RecipeKeyRequest, error) {
	if err := s.checkRecipePolicy(ctx, recipeId, userId); err != nil {
		return nil, err
	}

	requests := s.repo.GetRecipeKeyRequests(ctx, recipeId)

	var profileIds []string
	for _, request := range requests {
		profileIds = append(profileIds, request.UserId.String())
	}

	if len(requests) > 0 {
		callCtx, cancelCtx := context.WithTimeout(ctx, 3*time.Second)
		res, err := s.grpc.Profile.GetProfilesMinInfo(callCtx, &profileApi.GetProfilesMinInfoRequest{
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

func (s *Service) GetRecipeKey(ctx context.Context, recipeId, userId uuid.UUID) *[]byte {
	policy, err := s.grpc.Recipe.GetRecipePolicy(ctx, &recipeApi.GetRecipePolicyRequest{
		RecipeId: recipeId.String(),
	})
	if err != nil || !policy.IsEncrypted || policy.OwnerId != userId.String() && policy.Visibility == recipeModel.VisibilityPrivate {
		return nil
	}

	return s.repo.GetRecipeKey(ctx, recipeId, userId)
}

func (s *Service) RequestRecipeKeyAccess(ctx context.Context, recipeId, userId uuid.UUID) error {
	if userVaultKey := s.repo.GetEncryptedVault(ctx, userId); userVaultKey.PrivateKey == nil {
		return encryptionFail.GrpcNoVault
	}

	policy, err := s.grpc.Recipe.GetRecipePolicy(ctx, &recipeApi.GetRecipePolicyRequest{
		RecipeId: recipeId.String(),
	})
	if err != nil {
		return err
	}
	if policy.OwnerId == userId.String() {
		return nil
	}
	if !policy.IsEncrypted || policy.Visibility != recipeModel.VisibilityLink {
		return fail.GrpcAccessDenied
	}

	return s.repo.CreateRecipeKeyAccessRequest(ctx, recipeId, userId)
}

func (s *Service) SetRecipeKey(ctx context.Context, recipeId, userId uuid.UUID, key []byte, requesterId uuid.UUID) error {
	if err := s.checkRecipePolicy(ctx, recipeId, requesterId); err != nil {
		return err
	}

	if userVaultKey := s.repo.GetEncryptedVault(ctx, userId); userVaultKey.PrivateKey == nil {
		return encryptionFail.GrpcNoVault
	}

	if userId == requesterId {
		return s.repo.SetRecipeAuthorKey(ctx, recipeId, userId, key)
	} else {
		if ownerKey := s.repo.GetRecipeKey(ctx, recipeId, requesterId); ownerKey == nil {
			return encryptionFail.GrpcNoOwnerRecipeKey
		}

		return s.repo.GrantRecipeKeyAccessForUser(ctx, recipeId, userId, key)
	}
}

func (s *Service) DeleteRecipeUserKey(ctx context.Context, recipeId, userId, requesterId uuid.UUID) error {
	if err := s.checkRecipePolicy(ctx, recipeId, requesterId); err != nil {
		return err
	}
	if userId == requesterId {
		return encryptionFail.GrpcOwnedRecipeKeyDeletion
	}

	return s.repo.DeclineRecipeKeyAccessForUser(ctx, recipeId, userId)
}

func (s *Service) checkRecipePolicy(ctx context.Context, recipeId, requesterId uuid.UUID) error {
	policy, err := s.grpc.Recipe.GetRecipePolicy(ctx, &recipeApi.GetRecipePolicyRequest{
		RecipeId: recipeId.String(),
	})
	if err != nil {
		return err
	}
	if policy.OwnerId != requesterId.String() || !policy.IsEncrypted || policy.Visibility == recipeModel.VisibilityPublic {
		return fail.GrpcAccessDenied
	}
	return nil
}
