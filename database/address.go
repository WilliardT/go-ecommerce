package database

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAddressNotFound = errors.New("address not found")
	ErrUnauthorized    = errors.New("address does not belong to user")
)

// удаляет адрес пользователя по address_id
func DeleteAddress(ctx context.Context, db *pgxpool.Pool, userID string, addressID uuid.UUID) error {
	// Проверяем, что адрес принадлежит пользователю, и удаляем его
	result, err := db.Exec(ctx,
		"DELETE FROM addresses WHERE address_id = $1 AND user_id = $2",
		addressID, userID)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrAddressNotFound
	}

	return nil
}
