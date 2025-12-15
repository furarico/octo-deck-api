package main

import (
	"log"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/handler"
	"github.com/furarico/octo-deck-api/internal/repository"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
)

func main() {
	spec, err := openapi3.NewLoader().LoadFromFile("openapi/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}
	spec.Servers = nil

	router := gin.Default()

	router.Use(middleware.OapiRequestValidator(spec))

	cardRepository := repository.NewMockCardRepository()
	cardService := service.NewCardService(cardRepository)
	cardHandler := handler.NewHandler(cardService)

	api.RegisterHandlers(router, cardHandler)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
