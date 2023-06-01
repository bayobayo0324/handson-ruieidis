package main

import (
	"context"
	"fmt"
	"handson-ruieidis/infra"
	"handson-ruieidis/repository"
	"log"

	"github.com/joho/godotenv"
)

type account struct {
	Name string
	Age  string
}

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
	fmt.Println("Get key1:", val) // "val1"

	hSetKey := "account_1"
	sampleAccount := &account{
		Name: "kojikoji",
		Age:  "41",
	}
	if err := repo.HSet(ctx, hSetKey, sampleAccount); err != nil {
		log.Fatal(err)
	}

	hGetAll, err := repo.HGetAll(ctx, hSetKey)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range hGetAll {
		fmt.Printf("key:%s\nvalue:%s\n", key, value)
	}

	if err := repo.FlushDB(ctx); err != nil {
		log.Fatal(err)
	}
}
