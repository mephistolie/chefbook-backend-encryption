package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-encryption/internal/entity"
)

func (r *Repository) GetRecipeKeyRequests(recipeId uuid.UUID) []entity.RecipeKeyRequest {
	var requests []entity.RecipeKeyRequest

	query := fmt.Sprintf(`
		SELECT %[1]v.user_id, %[1]v.status, %[2]v.public_key
		FROM
			%[1]v
		LEFT JOIN
			%[2]v ON %[1]v.user_id=%[2]v.user_id
		WHERE
			%[1]v.recipe_id=$1 AND %[1]v.status<>'%[3]v'
	`, recipeKeysTable, vaultKeysTable, entity.RecipeKeyRequestStatusOwned)

	rows, err := r.db.Query(query, recipeId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to get recipe %s key requests: %s", recipeId, err)
		return []entity.RecipeKeyRequest{}
	}

	for rows.Next() {
		request := entity.RecipeKeyRequest{}

		if err := rows.Scan(&request.UserId, &request.Status, &request.PublicKey); err != nil {
			continue
		}
		if request.Status == entity.RecipeKeyRequestStatusApproved {
			request.PublicKey = nil
		}

		requests = append(requests, request)
	}

	return requests
}

func (r *Repository) GetRecipeKey(recipeId, userId uuid.UUID) (*[]byte, *[]byte) {
	var key *[]byte
	var iv *[]byte

	query := fmt.Sprintf(`
		SELECT %[1]v.key, %[2]v.iv
		FROM %[1]v
		LEFT JOIN
			%[2]v ON %[1]v.reicpe_id=%[2]v.reicpe_id
		WHERE recipe_id=$1 AND user_id=$2
	`, recipeKeysTable, recipeIvsTable)

	row := r.db.QueryRow(query, recipeId, userId)
	if err := row.Scan(&key, &iv); err != nil {
		log.Debugf("unable to get recipe %s key for user %s: %s", recipeId, userId, err)
	}

	return key, iv
}

func (r *Repository) CreateRecipeKeyAccessRequest(recipeId, userId uuid.UUID) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (recipe_id, user_id)
		VALUES ($1, $2)
	`, recipeKeysTable)

	if _, err := r.db.Exec(query, recipeId, userId); err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to create recipe %s key access request for user %s: %s", recipeId, userId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) SetRecipeAuthorKey(recipeId, userId uuid.UUID, key []byte, iv []byte) error {
	tx, err := r.startTransaction()
	if err != nil {
		return err
	}

	ivQuery := fmt.Sprintf(`
		INSERT INTO %s (recipe_id, iv)
		VALUES ($1, $2)
	`, recipeIvsTable)

	if _, err := tx.Exec(ivQuery, recipeId, key, iv); err != nil {
		_ = tx.Rollback()
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to set recipe %s IV for user %s: %s", recipeId, userId, err)
		return fail.GrpcUnknown
	}

	recipeQuery := fmt.Sprintf(`
		INSERT INTO %s (recipe_id, user_id, key, status)
		VALUES ($1, $2, $3, '%s')
	`, recipeKeysTable, entity.RecipeKeyRequestStatusOwned)

	if _, err := tx.Exec(recipeQuery, recipeId, userId, key); err != nil {
		_ = tx.Rollback()
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to set recipe %s key for user %s: %s", recipeId, userId, err)
		return fail.GrpcUnknown
	}

	return commitTransaction(tx)
}

func (r *Repository) GrantRecipeKeyAccessForUser(recipeId, userId uuid.UUID, key []byte) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET key=$3, status='approved'
		WHERE recipe_id=$1 AND user_id=$2
	`, recipeKeysTable)

	if _, err := r.db.Exec(query, recipeId, userId, key); err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to grant recipe %s key access for user %s: %s", recipeId, userId, err)
		return fail.GrpcNotFound
	}

	return nil
}

func (r *Repository) DeclineRecipeKeyAccessForUser(recipeId, userId uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET key=null, status='declined'
		WHERE recipe_id=$1 AND user_id=$2
	`, recipeKeysTable)

	if _, err := r.db.Exec(query, recipeId, userId); err != nil {
		if isUniqueViolationError(err) {
			return nil
		}
		log.Warnf("unable to decline recipe %s key access for user %s: %s", recipeId, userId, err)
		return fail.GrpcNotFound
	}

	return nil
}
