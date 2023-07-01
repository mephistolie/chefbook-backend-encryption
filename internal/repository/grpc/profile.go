package grpc

import (
	api "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Profile struct {
	api.ProfileServiceClient
	Conn *grpc.ClientConn
}

func NewProfile(addr string) (*Profile, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, opts)
	if err != nil {
		return nil, err
	}
	return &Profile{
		ProfileServiceClient: api.NewProfileServiceClient(conn),
		Conn:                 conn,
	}, nil
}
