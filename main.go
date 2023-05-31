package main

import (
	"context"
	"fmt"
	"handson-ruieidis/infra"
	"handson-ruieidis/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	repo := repository.NewCacheRepository(infra.Redis())
	defer infra.Redis().Close()

	ctx := context.Background()
	// SET key1 val1
	if err := repo.SetEx(ctx, "key1", "val1", 65535); err != nil {
		log.Fatal(err)
	}
	// GET key1
	val, err := repo.Get(ctx, "key1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("value", val) // "val1"
}
