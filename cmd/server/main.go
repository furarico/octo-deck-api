package main

import (
	"log"

	api "github.com/furarico/octo-deck-api/generated"
	"github.com/furarico/octo-deck-api/internal/database"
	"github.com/furarico/octo-deck-api/internal/handler"
	authmiddleware "github.com/furarico/octo-deck-api/internal/middleware"
	"github.com/furarico/octo-deck-api/internal/repository"
	"github.com/furarico/octo-deck-api/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	oapimiddleware "github.com/oapi-codegen/gin-middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	db, err := database.ConnectWithConnectorIAMAuthN()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	spec, err := openapi3.NewLoader().LoadFromFile("openapi/openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}
	spec.Servers = nil

	router := gin.Default()

	router.Use(oapimiddleware.OapiRequestValidatorWithOptions(spec, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}))
	router.Use(authmiddleware.AuthMiddleware())

	// identiconGen := identicon.NewGenerator()

	cardRepository := repository.NewCardRepository(db)
	communityRepository := repository.NewCommunityRepository(db)
	//cardRepository := repository.NewMockCardRepository()
	// TODO: 後ほどGitHub API Clientを注入
	cardService := service.NewCardService(cardRepository)
	communityService := service.NewCommunityService(communityRepository)
	h := handler.NewHandler(cardService, communityService)

	// StrictServerInterface を使用してハンドラーを登録
	strictHandler := api.NewStrictHandler(h, nil)
	api.RegisterHandlers(router, strictHandler)

	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
