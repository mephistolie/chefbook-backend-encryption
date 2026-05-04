package encryption

import (
	"context"

	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-auth/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

func (s *Service) HasEncryptedVault(ctx context.Context, userId uuid.UUID) bool {
	return s.repo.HasEncryptedVault(ctx, userId)
}

func (s *Service) GetEncryptedVault(ctx context.Context, userId uuid.UUID) entity.EncryptedVault {
	return s.repo.GetEncryptedVault(ctx, userId)
}

func (s *Service) CreateEncryptedVault(ctx context.Context, key entity.EncryptedVault) error {
	return s.repo.CreateEncryptedVault(ctx, key)
}

func (s *Service) RequestEncryptedVaultDeletion(ctx context.Context, userId uuid.UUID) error {
	deleteCode, err := s.repo.CreateVaultDeletionRequest(ctx, userId)
	if err == nil {
		go func() {
			ctx := context.WithoutCancel(ctx)
			info, err := s.grpc.Auth.GetAuthInfo(ctx, &api.GetAuthInfoRequest{Id: userId.String()})
			if err == nil {
				s.mail.SendEncryptedVaultDeletionMail(info.Email, deleteCode)
			}
		}()
	}
	return err
}

func (s *Service) DeleteEncryptedVault(ctx context.Context, userId uuid.UUID, deleteCode string) error {
	msg, err := s.repo.ConfirmEncryptedVaultDeletion(ctx, userId, deleteCode)
	if err == nil && msg != nil {
		go s.mqPublisher.PublishMessage(msg)
	}
	return err
}
