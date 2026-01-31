package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"ec-platform/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DBSet инициализирует соединение с PostgreSQL
func DBSet() *pgxpool.Pool {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Формат: postgres://username:password@localhost:5432/database_name
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	client, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return client
}

// FindUserByID - заглушка для поиска пользователя по ID.
// В реальном приложении здесь будет SQL-запрос к таблице "users".
func FindUserByID(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (*models.User, error) {
	// Это заглушка. Реальная функция будет выполнять SQL-запрос, например:
	// var user models.User
	// err := db.QueryRow(ctx, "SELECT ... FROM users WHERE id=$1", id).Scan(&user.ID, ...)
	// if err != nil {
	//     return nil, err
	// }
	// return &user, nil
	fmt.Printf("Заглушка: Поиск пользователя с ID %s\n", id)
	return nil, nil
}

// FindProductByID - заглушка для поиска продукта по ID.
// В реальном приложении здесь будет SQL-запрос к таблице "products".
func FindProductByID(ctx context.Context, db *pgxpool.Pool, id uuid.UUID) (*models.Product, error) {
	// Это заглушка.
	fmt.Printf("Заглушка: Поиск продукта с ID %s\n", id)
	return nil, nil
}
