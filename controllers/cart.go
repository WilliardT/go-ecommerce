package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")

		if productQueryID == "" {
			log.Println("product ID is empty")

			c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})

			return
		}

		// Получаем email пользователя из контекста (установлен middleware)
		email, exists := c.Get("email")

		if !exists {
			log.Println("user email not found in context")

			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})

			return
		}

		// Парсим UUID продукта
		productID, err := uuid.Parse(productQueryID)
		if err != nil {
			log.Printf("invalid product ID format: %v", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})

			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		// Получаем user_id по email
		var userID string

		err = app.DB.QueryRow(ctx, "SELECT user_id FROM users WHERE email = $1", email).Scan(&userID)

		if err != nil {
			log.Printf("error finding user: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})

			return
		}

		// Проверяем существование продукта
		var productExists bool

		err = app.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM products WHERE product_id = $1)", productID).Scan(&productExists)

		if err != nil {
			log.Printf("error checking product existence: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify product"})

			return
		}

		if !productExists {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})

			return
		}

		// Проверяем, есть ли уже этот продукт в корзине
		var cartItemExists bool

		err = app.DB.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM cart WHERE user_id = $1 AND product_id = $2)",
			userID, productID).Scan(&cartItemExists)

		if err != nil {
			log.Printf("error checking cart: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check cart"})

			return
		}

		if cartItemExists {
			// Увеличиваем количество, если товар уже в корзине
			_, err = app.DB.Exec(ctx,
				"UPDATE cart SET quantity = quantity + 1, updated_at = $1 WHERE user_id = $2 AND product_id = $3",
				time.Now().UTC(), userID, productID)

		} else {
			// Добавляем новый товар в корзину
			_, err = app.DB.Exec(ctx,
				"INSERT INTO cart (id, user_id, product_id, quantity, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
				uuid.New(), userID, productID, 1, time.Now().UTC(), time.Now().UTC())
		}

		if err != nil {
			log.Printf("error adding product to cart: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product to cart"})

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product added to cart", "product_id": productID})
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return nil
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return nil
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return nil
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return nil
}
