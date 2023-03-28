package routers

import (
	"github.com/Digisata/dts-hactiv8-golang-chap2/controllers"
	"github.com/gin-gonic/gin"
)

func StartServer(controller controllers.Controller) *gin.Engine {
	router := gin.Default()

	router.POST("/books", controller.CreateBook)
	router.GET("/books", controller.GetBook)
	router.GET("/books/:id", controller.GetBookById)
	router.PUT("/books/:id", controller.UpdateBook)
	router.DELETE("/books/:id", controller.DeleteBook)

	return router
}
