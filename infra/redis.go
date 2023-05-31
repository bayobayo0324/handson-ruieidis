package infra

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/rueidis"
)

var _redis rueidis.Client

func Redis() rueidis.Client {
	if _redis == nil {
		maxConnection, err := strconv.Atoi(os.Getenv("REDIS_MAX_CONNECTION"))
		if err != nil {
			log.Fatal(err)
		}
		client, err := rueidis.NewClient(
			rueidis.ClientOption{
				InitAddress:      []string{fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))},
				BlockingPoolSize: maxConnection,
				DisableCache:     true,
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		_redis = client
	}
	return _redis
}
