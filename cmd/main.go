package main

import (
	"Service/internal/app"

	"go.uber.org/fx"
)

func main() {
	fx.New(app.NewApp()).Run()
}
