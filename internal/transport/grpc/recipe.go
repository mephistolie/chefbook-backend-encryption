package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-encryption/api/proto/implementation/v1"
)

func (s *EncryptionServer) GetRecipeKeyRequests(_ context.Context, req *api.GetRecipeKeyRequestsRequest) (*api.GetRecipeKeyRequestsResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	recipeId, err := uuid.Parse(req.RecipeId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	requests, err := s.service.GetRecipeKeyRequests(recipeId, userId)
	if err != nil {
		return nil, err
	}

	response := make([]*api.RecipeKeyRequest, len(requests))
	for i, request := range requests {
		var key []byte
		if request.PublicKey != nil {
			key = *request.PublicKey
		}

		response[i] = &api.RecipeKeyRequest{
			UserId:     request.UserId.String(),
			UserName:   request.UserName,
			UserAvatar: request.UserAvatar,
			Status:     request.Status,
			PublicKey:  key,
		}
	}

	return &api.GetRecipeKeyRequestsResponse{Requests: response}, nil
}

func (s *EncryptionServer) GetRecipeKey(_ context.Context, req *api.GetRecipeKeyRequest) (*api.GetRecipeKeyResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	recipeId, err := uuid.Parse(req.RecipeId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	key := s.service.GetRecipeKey(recipeId, userId)
	var response []byte
	if key != nil {
		response = *key
	}

	return &api.GetRecipeKeyResponse{EncryptedKey: response}, nil
}

func (s *EncryptionServer) RequestRecipeKeyAccess(_ context.Context, req *api.RequestRecipeKeyAccessRequest) (*api.RequestRecipeKeyAccessResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	recipeId, err := uuid.Parse(req.RecipeId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	if err = s.service.RequestRecipeKeyAccess(recipeId, userId); err != nil {
		return nil, err
	}

	return &api.RequestRecipeKeyAccessResponse{Message: "recipe key access requested"}, nil
}

func (s *EncryptionServer) SetRecipeKey(_ context.Context, req *api.SetRecipeKeyRequest) (*api.SetRecipeKeyResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	recipeId, err := uuid.Parse(req.RecipeId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	requesterId, err := uuid.Parse(req.RequesterId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	if err = s.service.SetRecipeKey(recipeId, userId, req.EncryptedKey, requesterId); err != nil {
		return nil, err
	}

	return &api.SetRecipeKeyResponse{Message: "recipe key set"}, nil
}

func (s *EncryptionServer) DeleteRecipeKey(_ context.Context, req *api.DeleteRecipeKeyRequest) (*api.DeleteRecipeKeyResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	recipeId, err := uuid.Parse(req.RecipeId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	requesterId, err := uuid.Parse(req.RequesterId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	if err = s.service.DeleteRecipeKey(recipeId, userId, requesterId); err != nil {
		return nil, err
	}

	return &api.DeleteRecipeKeyResponse{Message: "recipe key deleted"}, nil
}
