package repository

import (
	"context"
	"reflect"
	"time"

	"github.com/redis/rueidis"
)

type CasheRepository interface {
	SetEx(ctx context.Context, key string, value string, ttl int64) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, target interface{}) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	FlushDB(ctx context.Context) error
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

func (r *casheRepository) HSet(ctx context.Context, key string, target interface{}) error {
	hSetFieldValue := r.client.B().Hset().Key(key).FieldValue()

	val := reflect.ValueOf(target).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fieldValue := val.Field(i).String()
		hSetFieldValue = hSetFieldValue.FieldValue(typeField.Name, fieldValue)
	}

	for _, resp := range r.client.DoMulti(
		ctx,
		hSetFieldValue.Build(),
		r.client.B().Expire().Key(key).Seconds(86400*int64(time.Second)).Build(),
	) {
		if err := resp.Error(); err != nil {
			return err
		}
	}
	return nil
}

func (r *casheRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result := r.client.Do(ctx, r.client.B().Hgetall().Key(key).Build())
	return result.AsStrMap()
}

func (r *casheRepository) FlushDB(ctx context.Context) error {
	return r.client.Do(ctx, r.client.B().Flushdb().Build()).Error()
}
