package grpc

import (
	"context"
	"crypto/x509"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-encryption/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
	encryptionFail "github.com/mephistolie/chefbook-backend-encryption/internal/entity/fail"
)

const (
	vaultPrivateKeyMinLength = 4000
	vaultPrivateKeyMaxLength = 5000
)

func (s *EncryptionServer) HasEncryptedVault(_ context.Context, req *api.HasEncryptedVaultRequest) (*api.HasEncryptedVaultResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	hasEncryptedVault := s.service.HasEncryptedVault(userId)

	return &api.HasEncryptedVaultResponse{HasEncryptedVault: hasEncryptedVault}, nil
}

func (s *EncryptionServer) GetEncryptedVaultKey(_ context.Context, req *api.GetEncryptedVaultKeyRequest) (*api.GetEncryptedVaultKeyResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	key := s.service.GetEncryptedVaultKey(userId)
	var response []byte
	if key != nil {
		response = *key
	}

	return &api.GetEncryptedVaultKeyResponse{EncryptedPrivateKey: response}, nil
}

func (s *EncryptionServer) CreateEncryptedVault(_ context.Context, req *api.CreateEncryptedVaultRequest) (*api.CreateEncryptedVaultResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	log.Debugf("supplied keys length is: public - %d; private - %d", len(req.PublicKey), len(req.EncryptedPrivateKey))

	privateKeyLength := len(req.EncryptedPrivateKey)
	if privateKeyLength < vaultPrivateKeyMinLength || privateKeyLength > vaultPrivateKeyMaxLength {
		return nil, encryptionFail.GrpcPrivateKeyLengthOutOfRange
	}

	_, err = x509.ParsePKCS1PublicKey(req.PublicKey)
	if err != nil {
		_, err = jwt.ParseRSAPublicKeyFromPEM(req.PublicKey)
		if err != nil {
			return nil, encryptionFail.GrpcInvalidPublicKey
		}
	}

	err = s.service.CreateEncryptedVault(entity.EncryptedVault{
		UserId:     userId,
		PublicKey:  &req.PublicKey,
		PrivateKey: &req.EncryptedPrivateKey,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateEncryptedVaultResponse{Message: "encrypted vault created"}, nil
}

func (s *EncryptionServer) RequestEncryptedVaultDeletion(_ context.Context, req *api.RequestEncryptedVaultDeletionRequest) (*api.RequestEncryptedVaultDeletionResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.RequestEncryptedVaultDeletion(userId)
	if err != nil {
		return nil, err
	}

	return &api.RequestEncryptedVaultDeletionResponse{Message: "encrypted vault deletion requested"}, nil
}

func (s *EncryptionServer) DeleteEncryptedVault(_ context.Context, req *api.DeleteEncryptedVaultRequest) (*api.DeleteEncryptedVaultResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil || !entity.IsDeleteCode(req.DeleteCode) {
		return nil, encryptionFail.GrpcInvalidCode
	}

	err = s.service.DeleteEncryptedVault(userId, req.DeleteCode)
	if err != nil {
		return nil, err
	}

	return &api.DeleteEncryptedVaultResponse{Message: "encrypted vault deleted"}, nil
}
