package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"ec-platform/database"
	"ec-platform/models"
	generate "ec-platform/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Инициализируем валидатор один раз для всего приложения
var validate = validator.New()

// Application будет хранить зависимости, такие как подключение к БД
type Application struct {
	DB *pgxpool.Pool
}

// хеширует пароль с использованием bcrypt
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Panic(err)
	}

	return string(hashedPassword)
}

// проверяет соответствие пароля хешу
func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "Login or password is incorrect"
		valid = false
	}

	return valid, msg
}

// обработчик для регистрации нового пользователя.
func (app *Application) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// Проводим валидацию на основе тегов `validate` в модели
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + validationErr.Error()})
			return
		}

		// Проверяем, существует ли пользователь с таким email или телефоном
		var count int

		query := "SELECT COUNT(*) FROM users WHERE email = $1 OR phone = $2"

		err := app.DB.QueryRow(ctx, query, user.Email, user.Phone).Scan(&count)

		if err != nil {
			log.Printf("Error checking for existing user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email or phone already exists"})
			return
		}

		// Хешируем пароль перед сохранением
		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

		// Генерируем ID и устанавливаем временные метки
		user.ID = uuid.New()
		user.Created_At = time.Now().UTC()
		user.Updated_At = time.Now().UTC()
		user.User_ID = user.ID.String()

		// Генерируем JWT токены
		token, refreshToken, err := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name)
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication tokens"})
			return
		}

		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.PoductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		// Вставляем нового пользователя в базу данных
		insertQuery := `
            INSERT INTO users (id, first_name, last_name, password, email, phone, user_id, token, refresh_token, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        `
		_, err = app.DB.Exec(ctx, insertQuery,
			user.ID, user.First_Name, user.Last_Name, user.Password, user.Email, user.Phone, user.User_ID, user.Token, user.Refresh_Token, user.Created_At, user.Updated_At,
		)

		if err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.User_ID})
	}
}

func (app *Application) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// Ищем пользователя в базе данных по email
		var foundUser models.User

		query := "SELECT id, first_name, last_name, password, email, phone, user_id, created_at, updated_at FROM users WHERE email = $1"

		err := app.DB.QueryRow(ctx, query, user.Email).Scan(
			&foundUser.ID,
			&foundUser.First_Name,
			&foundUser.Last_Name,
			&foundUser.Password,
			&foundUser.Email,
			&foundUser.Phone,
			&foundUser.User_ID,
			&foundUser.Created_At,
			&foundUser.Updated_At,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login or password incorrect"})
			return
		}

		// Проверяем пароль
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		// Генерируем новые токены
		token, refreshToken, err := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name)

		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication tokens"})
			return
		}

		// Обновляем токены в базе данных
		err = generate.UpdateAllTokens(app.DB, token, refreshToken, foundUser.User_ID)

		if err != nil {
			log.Printf("Error updating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		foundUser.Token = &token
		foundUser.Refresh_Token = &refreshToken

		c.JSON(http.StatusOK, foundUser)
	}
}

func (app *Application) ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			log.Printf("invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Проверяем обязательные поля
		if product.Product_Name == nil || *product.Product_Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_name is required"})
			return
		}

		if product.Price == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "price is required"})
			return
		}

		// Добавляем продукт в базу данных
		productID, err := database.AddProduct(ctx, app.DB, &product)

		if err != nil {
			log.Printf("error adding product: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "product added successfully",
			"product_id": productID,
		})
	}
}

func (app *Application) SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Получаем все Product из базы данных
		query := "SELECT product_id, product_name, price, rating, image FROM products ORDER BY product_name"

		rows, err := app.DB.Query(ctx, query)

		if err != nil {
			log.Printf("error fetching products: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}

		defer rows.Close()

		var productList []models.Product

		for rows.Next() {
			var product models.Product

			err := rows.Scan(&product.Product_ID, &product.Product_Name, &product.Price, &product.Rating, &product.Image)

			if err != nil {
				log.Printf("error scanning product: %v", err)
				continue
			}

			productList = append(productList, product)
		}

		c.JSON(http.StatusOK, gin.H{"products": productList})
	}
}

func (app *Application) SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("search query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "search query is required"})
			c.Abort()
			return
		}

		// Используем ILIKE для case-insensitive поиска в PostgreSQL
		query := "SELECT product_id, product_name, price, rating, image FROM products WHERE product_name ILIKE '%' || $1 || '%' ORDER BY product_name"

		rows, err := app.DB.Query(ctx, query, queryParam)

		if err != nil {
			log.Printf("error searching products: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search products"})
			return
		}

		defer rows.Close()

		var searchProducts []models.Product

		for rows.Next() {
			var product models.Product

			err := rows.Scan(&product.Product_ID, &product.Product_Name, &product.Price, &product.Rating, &product.Image)

			if err != nil {
				log.Printf("error scanning product: %v", err)
				continue
			}

			searchProducts = append(searchProducts, product)
		}

		c.JSON(http.StatusOK, gin.H{"products": searchProducts})
	}
}
