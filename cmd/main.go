package main

import (
	"github.com/Girag/calc-api-service/internal/application"
)

func main() {
	app := application.NewApp()
	app.RunServer()
}
