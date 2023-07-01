package mq

import "github.com/google/uuid"

const ExchangeEncryption = "encryption"
const AppId = "encryption-service"

const MsgTypeVaultDeleted = "vault.deleted"

type MsgBodyVaultDeleted struct {
	UserId uuid.UUID `json:"userId"`
}
