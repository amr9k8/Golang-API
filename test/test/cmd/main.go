package main

import (
	"test/pkg"
)

func main() {
	app := pkg.NewApp()
	app.Run("0.0.0.0:8080")
}
