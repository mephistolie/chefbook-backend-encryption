package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

func (r *Repository) CreateVaultDeletionRequest(userId uuid.UUID) (string, error) {
	var deleteCode string

	getExistingDeleteCodeQuery := fmt.Sprintf(`
		SELECT delete_code
		FROM %s
		WHERE user_id=$1
	`, vaultDeletionsTable)

	if err := r.db.Get(&deleteCode, getExistingDeleteCodeQuery, userId); err == nil {
		log.Infof("found existing vault delete code for user %s", userId)
		return deleteCode, nil
	}

	deleteCode = entity.GenerateDeleteCode()

	createDeleteCodeQuery := fmt.Sprintf(`
		INSERT INTO %s (user_id, delete_code)
		VALUES ($1, $2)
	`, vaultDeletionsTable)

	if _, err := r.db.Exec(createDeleteCodeQuery, userId, deleteCode); err != nil {
		log.Errorf("error while creating vault delete code for user %s: %s", userId, err)
		return "", fail.GrpcUnknown
	}

	return deleteCode, nil
}
