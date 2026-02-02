package main

import (
	"ec-platform/controllers"
	"ec-platform/database"
	"ec-platform/middleware"
	"ec-platform/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// Инициализируем подключение к базе данных
	db := database.DBSet()
	defer db.Close()

	// Создаем экземпляр приложения
	app := &controllers.Application{
		DB: db,
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Публичные роуты (без аутентификации)
	routes.UserRoutes(router, app)

	// Защищенные роуты (с аутентификацией)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", app.GetItemFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
