package main

import (
	"log"

	"git.dmitriygnatenko.ru/dima/homethings/internal/fiber"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	_ "github.com/lib/pq"
)

// @title 						Homethings API
// @version 					1.0
// @securityDefinitions.basic 	BasicAuth
// @BasePath 					/api
func main() {
	serviceProvider, err := sp.Init()
	if err != nil {
		log.Fatal(err)
	}

	fiberApp, err := fiber.Init(serviceProvider)
	if err != nil {
		log.Fatal(err)
	}

	if err = fiberApp.Listen(":" + serviceProvider.GetEnvService().GetAppPort()); err != nil {
		log.Fatal(err)
	}
}
