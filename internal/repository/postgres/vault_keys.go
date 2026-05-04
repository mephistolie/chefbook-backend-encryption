package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/mq/model"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-encryption/api/mq"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
	encryptionFail "github.com/mephistolie/chefbook-backend-encryption/internal/entity/fail"
)

func (r *Repository) HasEncryptedVault(ctx context.Context, userId uuid.UUID) bool {
	var rowsCount int

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE user_id=$1 AND public_key IS NOT NULL
	`, vaultKeysTable)

	row := r.db.QueryRowContext(ctx, query, userId)
	if err := row.Scan(&rowsCount); err != nil {
		log.Debugf("unable to get user %s encrypted vault exist state: %s", userId, err)
	}

	return rowsCount > 0
}

func (r *Repository) GetEncryptedVault(ctx context.Context, userId uuid.UUID) entity.EncryptedVault {
	vault := entity.EncryptedVault{UserId: userId}

	query := fmt.Sprintf(`
		SELECT public_key, private_key, salt
		FROM %s
		WHERE user_id=$1
	`, vaultKeysTable)

	row := r.db.QueryRowContext(ctx, query, userId)
	if err := row.Scan(&vault.PublicKey, &vault.PrivateKey, &vault.Salt); err != nil {
		log.Debugf("unable to get user %s encrypted vault: %s", userId, err)
	}

	return vault
}

func (r *Repository) CreateEncryptedVault(ctx context.Context, vault entity.EncryptedVault) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, public_key, private_key, salt)
		VALUES ($1, $2, $3, $4)
	`, vaultKeysTable)

	if _, err := r.db.ExecContext(ctx, query, vault.UserId, *vault.PublicKey, *vault.PrivateKey, *vault.Salt); err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to set profile %s vault: %s", vault.UserId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) ConfirmEncryptedVaultDeletion(ctx context.Context, userId uuid.UUID, deleteCode string) (*model.MessageData, error) {
	tx, err := r.startTransaction(ctx)
	if err != nil {
		return nil, err
	}

	consumeCodeQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id=$1 AND delete_code=$2
	`, vaultDeletionsTable)

	result, err := tx.ExecContext(ctx, consumeCodeQuery, userId, deleteCode)
	if err != nil {
		log.Warnf("unable to consume profile %s vault delete code: %s", userId, err)
		return nil, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return nil, errorWithTransactionRollback(tx, encryptionFail.GrpcInvalidCode)
	}

	if err = r.deleteVaultWithOwnedRecipeKeys(ctx, userId, tx); err != nil {
		return nil, errorWithTransactionRollback(tx, err)
	}

	msg, err := r.addOutboxVaultDeletedMsg(ctx, userId, tx)
	if err != nil {
		return nil, errorWithTransactionRollback(tx, err)
	}

	return msg, commitTransaction(tx)
}

func (r *Repository) deleteVaultWithOwnedRecipeKeys(ctx context.Context, userId uuid.UUID, tx *sql.Tx) error {
	getOwnedRecipesQuery := fmt.Sprintf(`
		SELECT recipe_id
		FROM %s
		WHERE user_id=$1 AND status='%s'
	`, recipeKeysTable, entity.RecipeKeyRequestStatusOwned)

	rows, err := tx.QueryContext(ctx, getOwnedRecipesQuery, userId)
	if err != nil {
		log.Errorf("unable to get owned recipes for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}
	defer rows.Close()

	var ownedRecipeIds []uuid.UUID
	for rows.Next() {
		var recipeId uuid.UUID
		if err = rows.Scan(&recipeId); err != nil {
			log.Errorf("unable to parse owned recipe ID for user %s: %s", userId, err)
			continue
		}
		ownedRecipeIds = append(ownedRecipeIds, recipeId)
	}
	if err = rows.Err(); err != nil {
		log.Errorf("unable to iterate owned recipes for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	deleteVaultQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id=$1
	`, vaultKeysTable)

	_, err = tx.ExecContext(ctx, deleteVaultQuery, userId)
	if err != nil {
		log.Warnf("unable to delete profile %s encrypted vault: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	deleteOwnedRecipesQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE recipe_id=ANY($1)
	`, recipeKeysTable)

	if _, err = tx.ExecContext(ctx, deleteOwnedRecipesQuery, ownedRecipeIds); err != nil {
		log.Errorf("unable to delete owned recipe keys for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return nil
}

func (r *Repository) addOutboxVaultDeletedMsg(ctx context.Context, userId uuid.UUID, tx *sql.Tx) (*model.MessageData, error) {
	msgBody := api.MsgBodyVaultDeleted{UserId: userId}
	msgBodyBson, err := json.Marshal(msgBody)
	if err != nil {
		log.Error("unable to marshal vault deleted message body: ", err)
		return nil, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}
	msgInfo := model.MessageData{
		Id:       uuid.New(),
		Exchange: api.ExchangeEncryption,
		Type:     api.MsgTypeVaultDeleted,
		Body:     msgBodyBson,
	}

	return &msgInfo, r.createOutboxMsg(ctx, &msgInfo, tx)
}
