package grpc

import (
	api "github.com/mephistolie/chefbook-backend-encryption/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-encryption/internal/transport/dependencies/service"
)

type EncryptionServer struct {
	api.UnsafeEncryptionServiceServer
	service           service.Encryption
	checkSubscription bool
}

func NewServer(service service.Encryption, checkSubscription bool) *EncryptionServer {
	return &EncryptionServer{
		service:           service,
		checkSubscription: checkSubscription,
	}
}
