package entity

import "github.com/google/uuid"

type EncryptedVault struct {
	UserId     uuid.UUID
	PublicKey  *[]byte
	PrivateKey *[]byte
}
