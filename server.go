package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"

	"github.com/senertopaloglu/eth-graphql-api/graph"
	"github.com/senertopaloglu/eth-graphql-api/graph/generated"
	"github.com/senertopaloglu/eth-graphql-api/internal/cache"
)

func main() {
	godotenv.Load()
	addr := ":8080"

	c := cache.New(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"), 0)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Cache: c}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost%s/ for GraphQL playground", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
