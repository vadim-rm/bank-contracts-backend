package main

import "github.com/vadim-rm/bmstu-web-backend/internal/app"

// @title           BMSTU Web Backend
// @version         1.0

// @host      localhost:8000

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

//go:generate swag init --parseDepth 2 --parseInternal -g cmd/bmstu_web_backend/bmstu_web_backend.go -d ../../ -o ../../docs
func main() {
	app.Run()
}
