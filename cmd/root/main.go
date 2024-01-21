package main

import "cloud/internal/root"

func main() {
	fx := root.NewApp()
	fx.Run()
}
