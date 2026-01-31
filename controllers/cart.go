package controllers

import (
	"context"
	"ec-platform/database"
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

		// Вызываем функцию из database слоя
		err = database.AddProductToCart(ctx, app.DB, userID, productID)

		if err != nil {
			if err == database.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})

			} else {
				log.Printf("error adding product to cart: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product to cart"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "product added to cart", "product_id": productID})
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
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

		// Вызываем функцию из database слоя
		err = database.RemoveCartItem(ctx, app.DB, userID, productID)

		if err != nil {
			if err == database.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "item not found in cart"})

			} else {
				log.Printf("error removing item from cart: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item from cart"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
	}
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
