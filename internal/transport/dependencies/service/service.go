package service

import (
	"context"

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
	HasEncryptedVault(ctx context.Context, userId uuid.UUID) bool
	GetEncryptedVault(ctx context.Context, userId uuid.UUID) entity.EncryptedVault
	CreateEncryptedVault(ctx context.Context, key entity.EncryptedVault) error
	RequestEncryptedVaultDeletion(ctx context.Context, userId uuid.UUID) error
	DeleteEncryptedVault(ctx context.Context, userId uuid.UUID, deleteCode string) error

	GetRecipeKeyRequests(ctx context.Context, recipeId uuid.UUID, userId uuid.UUID) ([]entity.RecipeKeyRequest, error)
	GetRecipeKey(ctx context.Context, recipeId, userId uuid.UUID) *[]byte
	RequestRecipeKeyAccess(ctx context.Context, recipeId, userId uuid.UUID) error
	SetRecipeKey(ctx context.Context, recipeId, userId uuid.UUID, key []byte, requesterId uuid.UUID) error
	DeleteRecipeUserKey(ctx context.Context, recipeId, userId, requesterId uuid.UUID) error
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
