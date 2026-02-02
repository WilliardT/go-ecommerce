package controllers

import (
	"context"
	"ec-platform/database"
	"ec-platform/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (app *Application) AddAdress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем email пользователя из контекста
		email, exists := c.Get("email")

		if !exists {
			log.Println("user email not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Парсим адрес из request body
		var address models.Address

		if err := c.BindJSON(&address); err != nil {
			log.Printf("invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Получаем user_id по email
		var userID string

		err := app.DB.QueryRow(ctx, "SELECT user_id FROM users WHERE email = $1", email).Scan(&userID)

		if err != nil {
			log.Printf("error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
			return
		}

		// Вызываем функцию из database слоя
		addressID, err := database.AddAddress(ctx, app.DB, userID, &address)

		if err != nil {
			log.Printf("error adding address: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add address"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "address added successfully",
			"address_id": addressID,
		})
	}
}

func (app *Application) EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		addressQueryID := c.Query("id")

		if addressQueryID == "" {
			log.Println("address ID is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "address ID is required"})
			return
		}

		// Получаем email пользователя из контекста (установлен middleware)
		email, exists := c.Get("email")

		if !exists {
			log.Println("user email not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Парсим UUID адреса
		addressID, err := uuid.Parse(addressQueryID)

		if err != nil {
			log.Printf("invalid address ID format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID format"})
			return
		}

		// Парсим обновленные данные адреса из request body
		var address models.Address

		if err := c.BindJSON(&address); err != nil {
			log.Printf("invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		err = database.UpdateAddress(ctx, app.DB, userID, addressID, &address)

		if err != nil {
			if err == database.ErrAddressNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})

			} else {
				log.Printf("error updating address: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address updated successfully"})
	}
}

func (app *Application) EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		addressQueryID := c.Query("id")

		if addressQueryID == "" {
			log.Println("address ID is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "address ID is required"})
			return
		}

		// Получаем email пользователя из контекста (установлен middleware)
		email, exists := c.Get("email")

		if !exists {
			log.Println("user email not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Парсим UUID адреса
		addressID, err := uuid.Parse(addressQueryID)

		if err != nil {
			log.Printf("invalid address ID format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID format"})
			return
		}

		// Парсим обновленные данные адреса из request body
		var address models.Address

		if err := c.BindJSON(&address); err != nil {
			log.Printf("invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		err = database.UpdateAddress(ctx, app.DB, userID, addressID, &address)

		if err != nil {
			if err == database.ErrAddressNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})

			} else {
				log.Printf("error updating address: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address updated successfully"})
	}
}

func (app *Application) DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		addressQueryID := c.Query("id")

		if addressQueryID == "" {
			log.Println("address ID is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "address ID is required"})
			return
		}

		// Получаем email пользователя из контекста (установлен middleware)
		email, exists := c.Get("email")

		if !exists {
			log.Println("user email not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Парсим UUID адреса
		addressID, err := uuid.Parse(addressQueryID)

		if err != nil {
			log.Printf("invalid address ID format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID format"})
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
		err = database.DeleteAddress(ctx, app.DB, userID, addressID)

		if err != nil {
			if err == database.ErrAddressNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})

			} else {
				log.Printf("error deleting address: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
			}

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "address deleted successfully"})
	}
}
