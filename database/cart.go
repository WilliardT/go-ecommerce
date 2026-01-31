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

// BuyItemFromCart выполняет покупку всех товаров из корзины пользователя
func BuyItemFromCart(ctx context.Context, db *pgxpool.Pool, userID string) (orderID uuid.UUID, totalPrice uint64, err error) {
	// Начинаем транзакцию
	tx, err := db.Begin(ctx)

	if err != nil {
		return uuid.Nil, 0, err
	}

	defer tx.Rollback(ctx)

	// Получаем все товары из корзины с их ценами
	query := `
		SELECT c.product_id, p.price, c.quantity
		FROM cart c
		JOIN products p ON c.product_id = p.product_id
		WHERE c.user_id = $1
	`

	rows, err := tx.Query(ctx, query, userID)

	if err != nil {
		return uuid.Nil, 0, err
	}

	defer rows.Close()

	type OrderItem struct {
		ProductID uuid.UUID
		Price     uint64
		Quantity  int
	}

	var orderItems []OrderItem
	var total uint64

	for rows.Next() {
		var item OrderItem

		err := rows.Scan(&item.ProductID, &item.Price, &item.Quantity)

		if err != nil {
			return uuid.Nil, 0, err
		}

		total += item.Price * uint64(item.Quantity)

		orderItems = append(orderItems, item)
	}

	if len(orderItems) == 0 {
		return uuid.Nil, 0, ErrCantGetItem
	}

	// Создаем заказ
	orderID = uuid.New()

	orderQuery := `
		INSERT INTO orders (order_id, user_id, total_price, ordered_at, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(ctx, orderQuery, orderID, userID, total, time.Now().UTC(), "pending")

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	// Добавляем товары в order_items
	for _, item := range orderItems {
		_, err = tx.Exec(ctx,
			"INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			uuid.New(), orderID, item.ProductID, item.Quantity, item.Price)

		if err != nil {
			return uuid.Nil, 0, ErrCantBuyCartItem
		}
	}

	// Очищаем корзину
	_, err = tx.Exec(ctx, "DELETE FROM cart WHERE user_id = $1", userID)

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	// Коммитим транзакцию
	err = tx.Commit(ctx)

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	return orderID, total, nil
}

// InstantBuyer выполняет мгновенную покупку одного товара
func InstantBuyer(ctx context.Context, db *pgxpool.Pool, userID string, productID uuid.UUID) (orderID uuid.UUID, totalPrice uint64, err error) {
	// Начинаем транзакцию
	tx, err := db.Begin(ctx)

	if err != nil {
		return uuid.Nil, 0, err
	}

	defer tx.Rollback(ctx)

	// Получаем информацию о продукте
	var price uint64

	err = tx.QueryRow(ctx, "SELECT price FROM products WHERE product_id = $1", productID).Scan(&price)

	if err != nil {
		return uuid.Nil, 0, ErrRecordNotFound
	}

	// Создаем заказ
	orderID = uuid.New()

	orderQuery := `
		INSERT INTO orders (order_id, user_id, total_price, ordered_at, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, orderQuery, orderID, userID, price, time.Now().UTC(), "pending")

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	// Добавляем товар в order_items
	_, err = tx.Exec(ctx,
		"INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
		uuid.New(), orderID, productID, 1, price)

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	// Коммитим транзакцию
	err = tx.Commit(ctx)

	if err != nil {
		return uuid.Nil, 0, ErrCantBuyCartItem
	}

	return orderID, price, nil
}
