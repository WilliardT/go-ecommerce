package database

import (
	"context"
	"ec-platform/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrAddressNotFound = errors.New("address not found")
	ErrUnauthorized    = errors.New("address does not belong to user")
)

// добавляет новый адрес для пользователя
func AddAddress(ctx context.Context, db *pgxpool.Pool, userID string, address *models.Address) (uuid.UUID, error) {
	addressID := uuid.New()

	query := `
		INSERT INTO addresses (address_id, user_id, house, street, city, pincode, state, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := db.Exec(ctx, query,
		addressID,
		userID,
		address.House,
		address.Street,
		address.City,
		address.Pincode,
		address.State,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	if err != nil {
		return uuid.Nil, err
	}

	return addressID, nil
}

// UpdateAddress обновляет адрес пользователя
func UpdateAddress(ctx context.Context, db *pgxpool.Pool, userID string, addressID uuid.UUID, address *models.Address) error {
	query := `
		UPDATE addresses
		SET house = $1, street = $2, city = $3, pincode = $4, state = $5, updated_at = $6
		WHERE address_id = $7 AND user_id = $8
	`

	result, err := db.Exec(ctx, query,
		address.House,
		address.Street,
		address.City,
		address.Pincode,
		address.State,
		time.Now().UTC(),
		addressID,
		userID,
	)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrAddressNotFound
	}

	return nil
}

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
