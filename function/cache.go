package function

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	m "github.com/Faizal-Asep/crud-app/model"
	"github.com/go-redis/redis/v8"
)

func SetCache(rds *redis.Client, key string, data interface{}) (err error) {

	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	ttl := time.Duration(5) * time.Minute
	op1 := rds.Set(context.Background(), key, b, ttl)
	return op1.Err()
}

func GetCache(rds *redis.Client, key string) (exist bool, result string, err error) {
	result, err = rds.Get(context.Background(), key).Result()
	if err != redis.Nil {
		exist = true
	}
	return
}

func CacheNewsKey(filter m.Newsfilter) string {
	value := reflect.ValueOf(filter)
	name := value.Type()
	var key string = "news"
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).Interface() == "" {
			continue
		}
		if key == "" {
			key += fmt.Sprintf("/%s:%s", name.Field(i).Name, value.Field(i).Interface())
		} else {
			key += fmt.Sprintf("/%s:%s", name.Field(i).Name, value.Field(i).Interface())
		}
	}
	return key

}
