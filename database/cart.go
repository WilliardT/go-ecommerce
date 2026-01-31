package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound     = errors.New("can't find the product")
	ErrCantDecodeProducts = errors.New("can't find (decode) products")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from the cart")
	ErrCantGetItem        = errors.New("can't get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
)

// добавляет продукт в корзину пользователя или увеличивает количество
func AddProductToCart(ctx context.Context, db *pgxpool.Pool, userID string, productID uuid.UUID) error {
	// проверяем существование продукта
	var productExists bool

	err := db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE product_id = $1)", productID).Scan(&productExists)

	if err != nil {
		return err
	}

	if !productExists {
		return ErrRecordNotFound
	}

	// проверяем, есть ли уже этот продукт в корзине
	var cartItemExists bool

	err = db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM cart WHERE user_id = $1 AND product_id = $2)",
		userID, productID).Scan(&cartItemExists)

	if err != nil {
		return err
	}

	if cartItemExists {
		// увеличиваем количество, если товар уже в корзине
		_, err = db.Exec(ctx,
			"UPDATE cart SET quantity = quantity + 1, updated_at = $1 WHERE user_id = $2 AND product_id = $3",
			time.Now().UTC(), userID, productID)

	} else {
		// добавляем новый товар в корзину
		_, err = db.Exec(ctx,
			"INSERT INTO cart (id, user_id, product_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			uuid.New(), userID, productID, 1, time.Now().UTC(), time.Now().UTC())
	}

	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

// удаляет продукт из корзины пользователя
func RemoveCartItem(ctx context.Context, db *pgxpool.Pool, userID string, productID uuid.UUID) error {
	result, err := db.Exec(ctx,
		"DELETE FROM cart WHERE user_id = $1 AND product_id = $2",
		userID, productID)

	if err != nil {
		return ErrCantRemoveItemCart
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
