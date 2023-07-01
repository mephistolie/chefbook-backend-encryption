package mq

import (
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
	log.Infof("processing message %s with type %s", msg.Id, msg.Type)
	switch msg.Type {
	case recipe.MsgTypeRecipeDeleted:
		return s.handleRecipeDeletedMsg(msg.Id, msg.Body)
	case auth.MsgTypeProfileDeleted:
		return s.handleProfileDeletedMsg(msg.Id, msg.Body)
	default:
		log.Warnf("got unsupported message type %s for message %s", msg.Type, msg.Id)
		return errors.New("not implemented")
	}
}

func (s *Service) handleRecipeDeletedMsg(messageId uuid.UUID, data []byte) error {
	var body recipe.MsgBodyRecipeDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	log.Infof("deleting recipe %s key...", body.RecipeId)
	return s.repo.DeleteRecipeKey(body.RecipeId, messageId)
}

func (s *Service) handleProfileDeletedMsg(messageId uuid.UUID, data []byte) error {
	var body auth.MsgBodyProfileDeleted
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		return err
	}

	log.Infof("deleting user %s...", body.UserId)
	return s.repo.DeleteProfile(userId, messageId)
}
