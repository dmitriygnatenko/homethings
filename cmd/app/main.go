package main

import (
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"git.dmitriygnatenko.ru/dima/homethings/internal/fiber"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
)

// @title 						Homethings API
// @version 					1.0
// @securitydefinitions.apikey 	APIKey
// @in 							header
// @name 						Authorization
func main() {
	serviceProvider, err := sp.Init()
	if err != nil {
		log.Fatal(err)
	}

	fiberApp, err := fiber.Init(serviceProvider)
	if err != nil {
		log.Fatal(err)
	}

	port := strconv.FormatUint(uint64(serviceProvider.ConfigService().AppPort()), 10)
	if err = fiberApp.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func init() {
	time.Local = time.UTC
}
