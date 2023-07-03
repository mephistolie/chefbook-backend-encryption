package service

import (
	"github.com/google/uuid"
	mq "github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	mqPublisher "github.com/mephistolie/chefbook-backend-common/mq/publisher"
	"github.com/mephistolie/chefbook-backend-encryption/internal/config"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
	"github.com/mephistolie/chefbook-backend-encryption/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/encryption"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/mail"
	mqInbox "github.com/mephistolie/chefbook-backend-encryption/internal/service/mq"
)

type Service struct {
	Encryption Encryption
	MQ         mq.Inbox
}

type Encryption interface {
	GetEncryptedVaultKey(userId uuid.UUID) *[]byte
	CreateEncryptedVault(key entity.EncryptedVault) error
	RequestEncryptedVaultDeletion(userId uuid.UUID) error
	DeleteEncryptedVault(userId uuid.UUID, deleteCode string) error

	GetRecipeKeyRequests(recipeId uuid.UUID, userId uuid.UUID) ([]entity.RecipeKeyRequest, error)
	GetRecipeKey(recipeId, userId uuid.UUID) *[]byte
	RequestRecipeKeyAccess(recipeId, userId uuid.UUID) error
	SetRecipeKey(recipeId, userId uuid.UUID, key []byte, requesterId uuid.UUID) error
	DeleteRecipeUserKey(recipeId, userId, requesterId uuid.UUID) error
}

func New(
	repo repository.Encryption,
	grpc *grpc.Repository,
	mqPublisher *mqPublisher.Publisher,
	cfg *config.Config,
) (*Service, error) {

	mailService, err := mail.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		Encryption: encryption.NewService(repo, grpc, mqPublisher, mailService),
		MQ:         mqInbox.NewService(repo),
	}, nil
}
