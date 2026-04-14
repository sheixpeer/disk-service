package main

import (
	"fmt"

	"github.com/sheixpeer/disk-service/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
	// TODO: init logger: slog

	// TODO: init storage/repository: postgres

	// TODO: init router: chi, "chi render"

	// TODO: run server

}
