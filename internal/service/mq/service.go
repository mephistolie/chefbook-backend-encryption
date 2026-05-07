package mq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	"github.com/mephistolie/chefbook-backend-encryption/internal/service/dependencies/repository"
	recipe "github.com/mephistolie/chefbook-backend-recipe/api/mq"
)

type Service struct {
	repo repository.Encryption
}

func NewService(
	repo repository.Encryption,
) *Service {
	return &Service{repo: repo}
}

func (s *Service) HandleMessage(msg model.MessageData) error {
	ctx := context.Background()
	log.Log(ctx, log.Event{
		Event:     "mq.message.processing",
		Message:   "processing message",
		Component: log.ComponentAMQP,
		MessageID: msg.Id.String(),
		Payload: map[string]any{
			"message_type": msg.Type,
		},
	})
	switch msg.Type {
	case recipe.MsgTypeRecipeDeleted:
		return s.handleRecipeDeletedMsg(ctx, msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(ctx, msg.Id, msg.Body)
	default:
		log.LogWarn(ctx, log.Event{
			Event:     "mq.message.unsupported_type",
			Message:   "got unsupported message type",
			Component: log.ComponentAMQP,
			MessageID: msg.Id.String(),
			Payload: map[string]any{
				"message_type": msg.Type,
			},
		})
		return errors.New("not implemented")
	}
}

func (s *Service) handleRecipeDeletedMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body recipe.MsgBodyRecipeDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	log.Log(ctx, log.Event{
		Event:     "recipe.deleted.message.processing",
		Message:   "processing recipe deleted message",
		Component: log.ComponentAMQP,
		MessageID: messageId.String(),
		Payload: map[string]any{
			"recipe_id": body.RecipeId,
		},
	})
	return s.repo.DeleteRecipeKeys(ctx, body.RecipeId, messageId)
}

func (s *Service) handleProfileDeletedMsg(ctx context.Context, messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	log.Log(ctx, log.Event{
		Event:     "profile.deleted.message.processing",
		Message:   "processing profile deleted message",
		Component: log.ComponentAMQP,
		MessageID: messageId.String(),
		UserID:    body.UserId,
	})
	return s.repo.DeleteProfile(ctx, userId, messageId)
}
