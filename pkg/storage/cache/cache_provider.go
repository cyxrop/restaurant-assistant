package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Provider struct {
	ctx    context.Context
	Client *redis.Client
}

func NewProvider(ctx context.Context, client *redis.Client) *Provider {
	return &Provider{
		ctx:    ctx,
		Client: client,
	}
}
