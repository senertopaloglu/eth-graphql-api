package graph

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/senertopaloglu/eth-graphql-api/graph/generated"
	"github.com/senertopaloglu/eth-graphql-api/graph/model"
	"github.com/senertopaloglu/eth-graphql-api/internal/cache"
	"github.com/senertopaloglu/eth-graphql-api/internal/client"
)

type Resolver struct {
	Cache *cache.Cache
}

func (r *queryResolver) Balance(ctx context.Context, address string) (string, error) {
	key := "balance:" + address
	if val, err := r.Cache.Get(key); err == nil {
		return val, nil
	}
	bal, err := client.GetBalance(address)
	if err != nil {
		return "", err
	}
	// cache for 1 minute
	_ = r.Cache.Set(key, bal, time.Minute)
	return bal, nil
}

func (r *queryResolver) Transactions(ctx context.Context, address string) ([]*model.Transaction, error) {
	key := "txns:" + address
	if raw, err := r.Cache.Get(key); err == nil {
		var cached []*model.Transaction
		_ = json.Unmarshal([]byte(raw), &cached)
		return cached, nil
	}
	txs, err := client.GetTransactions(address)
	if err != nil {
		return nil, err
	}
	out := make([]*model.Transaction, len(txs))
	for i, t := range txs {
		out[i] = &model.Transaction{
			Hash:      t.Hash,
			From:      t.From,
			To:        t.To,
			Value:     t.Value,
			Timestamp: t.TimeStamp, // string
		}
	}
	data, _ := json.Marshal(out)
	_ = r.Cache.Set(key, string(data), 5*time.Minute)
	return out, nil
}

func (r *queryResolver) TokenPrices(ctx context.Context, ids []string) ([]*model.TokenPrice, error) {
	key := "prices:" + strings.Join(ids, ",")
	if raw, err := r.Cache.Get(key); err == nil {
		var cached []*model.TokenPrice
		_ = json.Unmarshal([]byte(raw), &cached)
		return cached, nil
	}
	priceMap, err := client.GetTokenPrices(ids, "usd")
	if err != nil {
		return nil, err
	}
	var out []*model.TokenPrice
	for _, id := range ids {
		p := priceMap[id]
		out = append(out, &model.TokenPrice{
			ID:          id,
			Currency:    "usd",
			Price:       p.Price,
			LastUpdated: p.LastUpdated,
		})
	}
	data, _ := json.Marshal(out)
	_ = r.Cache.Set(key, string(data), 1*time.Minute)
	return out, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }
