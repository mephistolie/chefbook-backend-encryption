package entity

import "github.com/google/uuid"

const (
	RecipeKeyRequestStatusOwned    = "owned"
	RecipeKeyRequestStatusPending  = "pending"
	RecipeKeyRequestStatusApproved = "approved"
	RecipeKeyRequestStatusDeclined = "declined"
)

type RecipeKeyRequest struct {
	UserId     uuid.UUID
	UserName   *string
	UserAvatar *string
	PublicKey  *[]byte
	Status     string
}
