package entity

import (
	"github.com/mephistolie/chefbook-backend-common/random"
	"strconv"
)

const deleteVaultCodeLength = 6

func GenerateDeleteCode() string {
	return random.DigitString(deleteVaultCodeLength)
}

func IsDeleteCode(str string) bool {
	if len(str) != deleteVaultCodeLength {
		return false
	}
	_, err := strconv.Atoi(str)
	return err == nil
}
