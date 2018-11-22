package redis

import (
	"github.com/go-redis/redis"
	"sync"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

var once sync.Once
var client *redis.Client
var lock *sync.Mutex = new(sync.Mutex)

func GetClient() (*redis.Client, error) {

	// 单例模式
	once.Do(func() {
		cfg := config.GetConfig()
		addr := cfg.GetString("redis.addr")
		pass := cfg.GetString("redis.password")

		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass, // password set
			DB:       0,    // use default DB
		})

		util.LogInfo("redis client init.")
	})

	if _, err := client.Ping().Result(); err != nil {
		util.LogError(err)
		return nil, err
	}
	return client, nil
}

// 自增长ID
func GetId() int64 {
	lock.Lock()
	defer lock.Unlock()

	client, cerr := GetClient()
	if cerr != nil {
		util.LogError(cerr)
		return -1
	}

	result, err := client.Incr("alarm:increment").Result()
	if err != nil {
		util.LogError(err)
		return -1
	}
	return result
}

func GetPipelineArr(client *redis.Client, count int) []redis.Pipeliner {
	arr := make([]redis.Pipeliner, 0, count)
	for i := 0; i < count; i++ {
		arr = append(arr, client.Pipeline())
	}
	return arr
}

func ClosePipelineArr(pipes []redis.Pipeliner) {
	for _, pipe := range pipes {
		if pipe != nil {
			pipe.Close()
		}
	}
}

func ClosePipeline(pipe redis.Pipeliner) {
	if pipe != nil {
		pipe.Close()
	}
}
