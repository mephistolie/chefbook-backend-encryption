package grpc

import (
	api "github.com/mephistolie/chefbook-backend-recipe/api/proto/implementation/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Recipe struct {
	api.RecipeServiceClient
	Conn *grpc.ClientConn
}

func NewRecipe(addr string) (*Recipe, error) {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(addr, opts)
	if err != nil {
		return nil, err
	}
	return &Recipe{
		RecipeServiceClient: api.NewRecipeServiceClient(conn),
		Conn:                conn,
	}, nil
}
