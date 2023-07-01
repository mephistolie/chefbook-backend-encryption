package app

import (
	auth "github.com/mephistolie/chefbook-backend-auth/api/mq"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/config"
	mqConsumer "github.com/mephistolie/chefbook-backend-common/mq/consumer"
	mqApi "github.com/mephistolie/chefbook-backend-common/mq/dependencies"
	mqPublisher "github.com/mephistolie/chefbook-backend-common/mq/publisher"
	api "github.com/mephistolie/chefbook-backend-encryption/api/mq"
	mqRecipeApi "github.com/mephistolie/chefbook-backend-encryption/api/mq"
	"github.com/mephistolie/chefbook-backend-encryption/internal/repository/postgres"
	recipe "github.com/mephistolie/chefbook-backend-recipe/api/mq"
	amqp "github.com/wagslane/go-rabbitmq"
)

const queueProfiles = "encryption.profiles"
const queueRecipes = "encryption.recipes"

var supportedMsgTypes = []string{
	recipe.MsgTypeRecipeDeleted,
	auth.MsgTypeProfileDeleted,
}

func NewMqPublisher(
	cfg config.Amqp,
	repository *postgres.Repository,
) (*mqPublisher.Publisher, error) {
	var publisher *mqPublisher.Publisher = nil
	var err error

	if len(*cfg.Host) > 0 {
		publisher, err = mqPublisher.New(mqRecipeApi.AppId, cfg, repository)
		if err != nil {
			return nil, err
		}
		if err = publisher.Start(
			amqp.WithPublisherOptionsExchangeName(api.ExchangeEncryption),
			amqp.WithPublisherOptionsExchangeKind("fanout"),
			amqp.WithPublisherOptionsExchangeDurable,
			amqp.WithPublisherOptionsExchangeDeclare,
		); err != nil {
			return nil, err
		}

		log.Info("MQ Publisher initialized")
	}

	return publisher, nil
}

func NewMqSubscriber(
	cfg config.Amqp,
	service mqApi.Inbox,
) (*mqConsumer.Consumer, error) {
	var consumer *mqConsumer.Consumer = nil
	var err error

	if len(*cfg.Host) > 0 {

		consumer, err = mqConsumer.New(cfg, service, supportedMsgTypes)
		if err != nil {
			return nil, err
		}
		if err = consumer.Start(
			mqConsumer.Params{
				QueueName: queueProfiles,
				Options: []func(*amqp.ConsumerOptions){
					amqp.WithConsumerOptionsQueueQuorum,
					amqp.WithConsumerOptionsQueueDurable,
					amqp.WithConsumerOptionsExchangeName(auth.ExchangeProfiles),
					amqp.WithConsumerOptionsExchangeKind("fanout"),
					amqp.WithConsumerOptionsExchangeDurable,
					amqp.WithConsumerOptionsExchangeDeclare,
					amqp.WithConsumerOptionsRoutingKey(""),
				},
			},
			mqConsumer.Params{
				QueueName: queueRecipes,
				Options: []func(*amqp.ConsumerOptions){
					amqp.WithConsumerOptionsQueueQuorum,
					amqp.WithConsumerOptionsQueueDurable,
					amqp.WithConsumerOptionsExchangeName(recipe.ExchangeRecipes),
					amqp.WithConsumerOptionsExchangeKind("fanout"),
					amqp.WithConsumerOptionsExchangeDurable,
					amqp.WithConsumerOptionsExchangeDeclare,
					amqp.WithConsumerOptionsRoutingKey(""),
				},
			},
		); err != nil {
			return nil, err
		}

		log.Info("MQ Consumer initialized")
	}

	return consumer, nil
}
