package main

import (
	"fmt"

	"github.com/Digisata/dts-hactiv8-golang-chap2/controllers"
	"github.com/Digisata/dts-hactiv8-golang-chap2/database"
	"github.com/Digisata/dts-hactiv8-golang-chap2/routers"
)

func main() {
	database.StartDB()
	db := database.GetDB()

	bookController := controllers.NewBookController(db)

	routers.StartServer(bookController).Run(fmt.Sprintf(":%s", "3000"))
}
