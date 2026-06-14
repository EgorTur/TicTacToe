package main

import (
	"tic-tac-toe/internal/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(di.App).Run()
}
