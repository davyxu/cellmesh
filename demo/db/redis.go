package db

import (
	"github.com/davyxu/cellmesh/demo/db/model"
	"github.com/go-redis/redis"
	"time"
)

type ObjectKeyFetcher interface {

	// 获取Redis中主Key（最外面的一层）
	GetMainKey(env interface{}) string

	// 获取HASH类型值的key, hash容器的key，一般为玩家身上的
	GetHashKey(env interface{}) string
}

type RedisClient = redis.Client

var (
	GamePlayDB *RedisClient
)

func Init() {
	GamePlayDB = redis.NewClient(&redis.Options{
		Addr:         ":16379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	var err error

	err = GamePlayDB.Ping().Err()
	if err != nil {
		panic(err)
	}

	var ri model.RoleInfo
	ri.Name = "davy"
	ri.Level = 60

	data, err := ri.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	GamePlayDB.HSet(ri.GetMainKey("2315"), ri.GetHashKey(nil), data)

	var back model.RoleInfo

	v, err := GamePlayDB.HGet(ri.GetMainKey("2315"), ri.GetHashKey(nil)).Bytes()
	if err != nil {
		panic(err)
	}

	back.UnmarshalMsg(v)
}

func SaveObject(client *RedisClient, fetcher ObjectKeyFetcher, env interface{}) {

	data, err := fetcher.(interface {
		MarshalMsg(b []byte) (o []byte, err error)
	}).MarshalMsg(nil)

	if err != nil {
		panic(err)
	}

	client.HSet(fetcher.GetMainKey(env), fetcher.GetHashKey(nil), data)
}
