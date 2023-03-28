package main

import (
	"github.com/Digisata/dts-hactiv8-golang-chap2/routers"
)

func main() {
	PORT := ":3000"

	routers.StartServer().Run(PORT)
}
