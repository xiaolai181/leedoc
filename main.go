package main

import (
	_ "leedoc/conf"

	"leedoc/routers"
)

func main() {
	r := routers.Init_router()
	r.Run(":8080")
}
 