package repository

import (
	"context"

	"github.com/redis/rueidis"
)

type CasheRepository interface {
	SetEx(ctx context.Context, key string, value string, ttl int64) error
	Get(ctx context.Context, key string) (string, error)
}

type casheRepository struct {
	client rueidis.Client
}

func NewCacheRepository(client rueidis.Client) CasheRepository {
	repo := &casheRepository{
		client: client,
	}

	return repo
}

func (r *casheRepository) SetEx(ctx context.Context, key string, value string, ttl int64) error {
	return r.client.Do(ctx, r.client.B().Setex().Key(key).Seconds(ttl).Value(value).Build()).Error()
}

func (r *casheRepository) Get(ctx context.Context, key string) (string, error) {
	result := r.client.Do(ctx, r.client.B().Get().Key(key).Build())
	return result.ToString()
}
