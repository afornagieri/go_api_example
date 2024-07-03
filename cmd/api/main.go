package main

import (
	"fmt"
	"net/http"

	"github.com/afornagieri/go_api_template/internal/adapter/router"
	"github.com/afornagieri/go_api_template/internal/infra/di"
)

func main() {
	container := di.NewContainer()

	r := router.NewRouter(container.ItemController)

	port := ":8080"
	fmt.Printf("Server initialized. Running on port %s\n", port)

	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}
