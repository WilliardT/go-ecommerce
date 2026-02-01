package routes

import (
	"ec-platform/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine, app *controllers.Application) {
	incomingRoutes.POST("/users/signup", app.SignUp())
	incomingRoutes.POST("/users/login", app.Login())
	incomingRoutes.POST("/admin/addproduct", app.ProductViewerAdmin())
	incomingRoutes.GET("/users/productview", app.SearchProduct())
	incomingRoutes.GET("/users/search", app.SearchProductByQuery())
}
