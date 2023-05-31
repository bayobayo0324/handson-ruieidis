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
	expire := 86400

	val := reflect.ValueOf(target).Elem()

	for i := 0; i < val.NumField(); i++ {
		hSetFieldValue := r.client.B().Hset().Key(key).FieldValue()
		typeField := val.Type().Field(i)
		fieldName := typeField.Name
		fieldValue := val.FieldByName(typeField.Name).Interface()
		hSetFieldValue = hSetFieldValue.FieldValue(fieldName, fieldValue.(string))
		if err := r.client.Do(ctx, hSetFieldValue.Build()).Error(); err != nil {
			return err
		}
	}

	ttl := int64(expire * int(time.Second))
	if err := r.client.Do(ctx, r.client.B().Expire().Key(key).Seconds(ttl).Build()).Error(); err != nil {
		return err
	}
	return nil
}

func (r *casheRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result := r.client.Do(ctx, r.client.B().Hgetall().Key(key).Build())
	return result.AsStrMap()
}
