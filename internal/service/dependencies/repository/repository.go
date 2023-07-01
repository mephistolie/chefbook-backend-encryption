package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

type Encryption interface {
	GetEncryptedVault(userId uuid.UUID) entity.EncryptedVault
	CreateEncryptedVault(vault entity.EncryptedVault) error
	CreateVaultDeletionRequest(userId uuid.UUID) (string, error)
	DeleteEncryptedVault(userId uuid.UUID, deleteCode string) (*model.MessageData, error)

	GetRecipeKeyRequests(recipeId uuid.UUID, ownerId uuid.UUID) []entity.RecipeKeyRequest
	GetRecipeKey(recipeId, userId uuid.UUID) *[]byte
	SetRecipeAuthorKey(recipeId, userId uuid.UUID, key []byte) error
	CreateRecipeKeyAccessRequest(recipeId, userId uuid.UUID) error
	GrantRecipeKeyAccessForUser(recipeId, userId uuid.UUID, key []byte) error
	DeclineRecipeKeyAccessForUser(recipeId, userId uuid.UUID) error
	DeleteRecipeAuthorKey(recipeId, userId uuid.UUID) error

	DeleteRecipeKey(recipeId, messageId uuid.UUID) error
	DeleteProfile(userId, messageId uuid.UUID) error
}
