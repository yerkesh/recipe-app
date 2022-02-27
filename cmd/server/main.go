package main

import (
	"context"
	"log"
	"net/http"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/service"
)

func main() {
	ctx := context.Background()
	handlerCtx := handler.NewHandlerCtx(ctx, handler.WithRecipeService(service.NewRecipeService()))

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err := http.ListenAndServe(config.Configuration.Server.Port, router.Router(handlerCtx))

	pool.Close()

	log.Fatal(err)

}