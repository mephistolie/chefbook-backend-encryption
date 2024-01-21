package encryption

import (
	"context"
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-auth/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

func (s *Service) HasEncryptedVault(userId uuid.UUID) bool {
	return s.repo.HasEncryptedVault(userId)
}

func (s *Service) GetEncryptedVaultKey(userId uuid.UUID) *[]byte {
	return s.repo.GetEncryptedVault(userId).PrivateKey
}

func (s *Service) CreateEncryptedVault(key entity.EncryptedVault) error {
	return s.repo.CreateEncryptedVault(key)
}

func (s *Service) RequestEncryptedVaultDeletion(userId uuid.UUID) error {
	deleteCode, err := s.repo.CreateVaultDeletionRequest(userId)
	if err == nil {
		go func() {
			info, err := s.grpc.Auth.GetAuthInfo(context.Background(), &api.GetAuthInfoRequest{Id: userId.String()})
			if err == nil {
				s.mail.SendEncryptedVaultDeletionMail(info.Email, deleteCode)
			}
		}()
	}
	return err
}

func (s *Service) DeleteEncryptedVault(userId uuid.UUID, deleteCode string) error {
	msg, err := s.repo.ConfirmEncryptedVaultDeletion(userId, deleteCode)
	if err == nil && msg != nil {
		go s.mqPublisher.PublishMessage(msg)
	}
	return err
}
