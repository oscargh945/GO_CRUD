package main

import (
	"github.com/joho/godotenv"
	"github.com/oscargh945/go-crud-graphql/domain/repositories"
	"github.com/oscargh945/go-crud-graphql/domain/usecase"
	"github.com/oscargh945/go-crud-graphql/graph"
	"github.com/oscargh945/go-crud-graphql/infrastructure"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mongoConnection := infrastructure.Connect()

	repository := &repositories.UserRepository{
		Client: mongoConnection,
	}

	userUseCase := usecase.UserUseCase{
		Repository: repositories.UserRepository{
			Client: mongoConnection,
		},
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		UserUseCase: userUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", repository.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
