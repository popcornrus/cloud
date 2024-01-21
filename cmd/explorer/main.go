package main

import (
	"cloud/internal/explorer"
	"log/slog"
)

func main() {
	slog.Info("starting app")

	fx := explorer.NewApp()
	fx.Run()
}
