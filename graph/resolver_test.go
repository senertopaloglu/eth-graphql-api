package graph

import (
	"context"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-redis/redis/v8"
	"github.com/senertopaloglu/eth-graphql-api/graph/generated"
	"github.com/senertopaloglu/eth-graphql-api/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestBalanceResolver(t *testing.T) {
	// use a local ephemeral Redis (DB=1)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), DB: 1,
	})
	defer rdb.FlushDB(context.Background())

	exec := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{Cache: cache.New(os.Getenv("REDIS_ADDR"), "", 1)}}),
	)
	c := client.New(exec)

	var resp struct{ Balance string }
	query := `query { balance(address: "0xde0B295669a9FD93d5F28D9Ec85E40f4cb697BAe") }`
	c.MustPost(query, &resp)
	assert.NotEmpty(t, resp.Balance)
}

func TestTokenPricesResolver(t *testing.T) {
	exec := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{Cache: cache.New(os.Getenv("REDIS_ADDR"), "", 1)}}),
	)
	c := client.New(exec)

	var resp struct {
		TokenPrices []struct {
			ID          string
			Currency    string
			Price       float64
			LastUpdated string
		}
	}
	query := `query { tokenPrices(ids: ["ethereum"]) { id currency price last_updated } }`
	c.MustPost(query, &resp)
	assert.Equal(t, "ethereum", resp.TokenPrices[0].ID)
	assert.Equal(t, "usd", resp.TokenPrices[0].Currency)
	assert.Greater(t, resp.TokenPrices[0].Price, 0.0)
	assert.NotEmpty(t, resp.TokenPrices[0].LastUpdated)
}
