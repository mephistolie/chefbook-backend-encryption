package encryption

import (
	mqPublisher "github.com/mephistolie/chefbook-backend-common/mq/publisher"
	"github.com/mephistolie/chefbook-backend-encryption/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/mail"
)

type Service struct {
	repo        repository.Encryption
	grpc        *grpc.Repository
	mqPublisher *mqPublisher.Publisher
	mail        *mail.Service
}

func NewService(
	repo repository.Encryption,
	grpc *grpc.Repository,
	mqPublisher *mqPublisher.Publisher,
	mail *mail.Service,
) *Service {
	return &Service{
		repo:        repo,
		grpc:        grpc,
		mqPublisher: mqPublisher,
		mail:        mail,
	}
}
