package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"recipe-app/internal/config"
	"recipe-app/pkg/handler"
	"recipe-app/pkg/service"
	"recipe-app/router"
)

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *config.Config) {
	f, err := os.Open("resources/configs/config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func main() {
	var cfg config.Config
	readFile(&cfg)
	//fmt.Printf("%+v", cfg)

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, cfg.Database.URI)
	if err != nil {
		log.Printf("Unable to connect to database err: %v", err)

		return
	}

	handlerCtx := handler.NewHandlerCtx(ctx, handler.WithRecipeService(service.NewRecipeService()))

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	err = http.ListenAndServe(cfg.Server.Port, router.Router(handlerCtx))

	pool.Close()

	log.Fatal(err)
}