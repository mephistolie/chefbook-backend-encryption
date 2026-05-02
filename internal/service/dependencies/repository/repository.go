package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

type Encryption interface {
	HasEncryptedVault(ctx context.Context, userId uuid.UUID) bool
	GetEncryptedVault(ctx context.Context, userId uuid.UUID) entity.EncryptedVault
	CreateEncryptedVault(ctx context.Context, vault entity.EncryptedVault) error
	CreateVaultDeletionRequest(ctx context.Context, userId uuid.UUID) (string, error)
	ConfirmEncryptedVaultDeletion(ctx context.Context, userId uuid.UUID, deleteCode string) (*model.MessageData, error)

	GetRecipeKeyRequests(ctx context.Context, recipeId uuid.UUID) []entity.RecipeKeyRequest
	GetRecipeKey(ctx context.Context, recipeId, userId uuid.UUID) *[]byte
	SetRecipeAuthorKey(ctx context.Context, recipeId, userId uuid.UUID, key []byte) error
	CreateRecipeKeyAccessRequest(ctx context.Context, recipeId, userId uuid.UUID) error
	GrantRecipeKeyAccessForUser(ctx context.Context, recipeId, userId uuid.UUID, key []byte) error
	DeclineRecipeKeyAccessForUser(ctx context.Context, recipeId, userId uuid.UUID) error

	DeleteRecipeKeys(ctx context.Context, recipeId, messageId uuid.UUID) error
	DeleteProfile(ctx context.Context, userId, messageId uuid.UUID) error
}
