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
	ErrProductExists = errors.New("product already exists")
)

// AddProduct добавляет новый продукт в каталог
func AddProduct(ctx context.Context, db *pgxpool.Pool, product *models.Product) (uuid.UUID, error) {
	productID := uuid.New()

	query := `
		INSERT INTO products (product_id, product_name, price, rating, image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(ctx, query,
		productID,
		product.Product_Name,
		product.Price,
		product.Rating,
		product.Image,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	if err != nil {
		return uuid.Nil, err
	}

	return productID, nil
}
